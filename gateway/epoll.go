package gateway

import (
	"fmt"
	"golang.org/x/sys/unix"
	"log"
	"net"
	"plato/common/config"
	"plato/common/logger"
	"reflect"
	"runtime"
	"sync"
	"sync/atomic"
	"syscall"
)

// 全局对象
var ep *ePool    // epoll 池
var tcpNum int32 // 当前服务允许接入的最大tcp连接数

// epoll 池
type ePool struct {
	eChan chan *connection
	table sync.Map      // 存储 fd 到 连接 的映射
	eSize int           // 池中 epoll 个数
	done  chan struct{} // 关闭信号

	ln      *net.TCPListener                 // 服务连接
	runProc func(c *connection, ep *epoller) // 回调函数
}

func initEpollPool(ln *net.TCPListener, runProc func(c *connection, ep *epoller)) {
	setLimit()
	ep = newEPool(ln, runProc)
	ep.createAcceptProcess()
	ep.startEPoolProcess()
}

func newEPool(ln *net.TCPListener, runProc func(c *connection, ep *epoller)) *ePool {
	return &ePool{
		eChan:   make(chan *connection, config.GetGatewayEpollerChanNum()),
		table:   sync.Map{},
		eSize:   config.GetGatewayEpollerNum(),
		done:    make(chan struct{}),
		ln:      ln,
		runProc: runProc,
	}
}

// 创建专门处理 accept 事件的协程
func (e *ePool) createAcceptProcess() {
	// 协程数与当前cpu的核数对应，能够发挥最大功效
	for i := 0; i < runtime.NumCPU(); i++ {
		go func() {
			for {
				conn, err := e.ln.AcceptTCP()
				if err != nil {
					fmt.Printf("accept err: %v\n", e)
					continue
				}

				// 到达连接上限, 限流熔断
				if !checkTcp() {
					_ = conn.Close()
					continue
				}

				setTcpConifg(conn)

				connect := &connection{
					fd:   socketFD(conn),
					conn: conn,
				}
				e.addTask(connect)
			}
		}()
	}
}

func (e *ePool) addTask(c *connection) {
	e.eChan <- c
}

// 轮询器池的处理器
func (e *ePool) startEPoolProcess() {
	for i := 0; i < e.eSize; i++ {
		go func() {
			// 创建一个 epoll
			epoll, err := newEpoller()
			if err != nil {
				panic(err)
			}

			// 专门 把socket注册到epoll 的协程
			go func() {
				for {
					select {
					case <-e.done:
						return
					case conn := <-e.eChan:
						if err := epoll.add(conn); err != nil {
							fmt.Printf("failed to add connection %v\n", err)
							conn.Close()
							continue
						}
						fmt.Printf(
							"EpollerPool new connection[%v] tcpSize:%d\n",
							conn.RemoteAddr(),
							getTcpNum())
					}
				}
			}()

			// 专门轮询epoll的协程
			for {
				select {
				case <-e.done:
					return
				default:
					connections, err := epoll.wait(200) // 200 ms 轮询一次
					if err != nil && err != syscall.EINTR {
						fmt.Printf("failed to epoll wait %v\n", err)
						return
					}
					for _, conn := range connections {
						//logger.Debug("range conn brfore runProc")
						if conn == nil {
							continue
						}
						logger.Debug("range conn brfore1 runProc")
						e.runProc(conn, epoll)
					}
				}
			}
		}()
	}
}

// 针对 epoll 对象的 轮询器
type epoller struct {
	fd int // epoll 对象的 fd
}

func newEpoller() (*epoller, error) {
	// 通过系统调用创建 epoll
	fd, err := unix.EpollCreate1(0)
	if err != nil {
		_ = unix.Close(fd)
		return nil, err
	}
	return &epoller{fd: fd}, nil
}

// 水平触发模式
func (e *epoller) add(conn *connection) error {
	addTcpNum()
	fd := conn.fd

	// 将连接注册到 epoll 中监听
	err := unix.EpollCtl(
		e.fd,
		syscall.EPOLL_CTL_ADD,
		fd,
		&unix.EpollEvent{
			Events: unix.EPOLLIN | unix.EPOLLHUP,
			Fd:     int32(fd)},
	)
	if err != nil {
		return err
	}

	// 存储 connect fd -> connect 结构的映射
	// epoll 监听到 fd 就绪后可以从中取出 connect 结构进行处理
	ep.table.Store(fd, conn)
	return nil
}

func (e *epoller) remove(conn *connection) error {
	subTcpNum()
	fd := conn.fd
	err := unix.EpollCtl(
		e.fd,
		syscall.EPOLL_CTL_DEL,
		fd,
		nil,
	)
	if err != nil {
		return err
	}
	ep.table.Delete(fd)
	return nil
}

func (e *epoller) wait(msec int) ([]*connection, error) {
	events := make([]unix.EpollEvent, config.GetGatewayEpollWaitQueueSize())
	n, err := unix.EpollWait(e.fd, events, msec)
	if err != nil {
		return nil, err
	}

	// 遍历就绪的 fd, 转换回 connection 结构
	var connections []*connection
	for i := 0; i < n; i++ {
		conn, ok := ep.table.Load(int(events[i].Fd))
		if ok {
			connections = append(connections, conn.(*connection))
		}
	}

	if n > 0 {
		logger.Debug("epoll: some connection is ready")
	}

	return connections, nil
}

//
// 工具函数
//

// 设置进程打开文件数目限制
func setLimit() {
	var rLimit syscall.Rlimit
	var err error

	err = syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		log.Fatalf("setLimit fails: %v", err)
	}
	rLimit.Cur = rLimit.Max
	err = syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		log.Fatalf("setLimit fails: %v", err)
	}

	log.Printf("set cur limit: %d", rLimit.Cur)
}

// 获取、修改当前连接数
func getTcpNum() int32 {
	return atomic.LoadInt32(&tcpNum)
}
func addTcpNum() {
	atomic.AddInt32(&tcpNum, 1)
}
func subTcpNum() {
	atomic.AddInt32(&tcpNum, -1)
}

// 检查连接数是否已大于最大连接限制
func checkTcp() bool {
	num := getTcpNum()
	maxTcpNum := config.GetGatewayMaxTcpNum()
	return num <= maxTcpNum
}

// tcp 设置为长连接
func setTcpConifg(c *net.TCPConn) {
	_ = c.SetKeepAlive(true)
}

// 获取 tcp 连接的 fd
// 利用反射获取私有属性
func socketFD(conn *net.TCPConn) int {
	tcpConn := reflect.Indirect(reflect.ValueOf(*conn)).FieldByName("conn")
	fdVal := tcpConn.FieldByName("fd")
	pfdVal := reflect.Indirect(fdVal).FieldByName("pfd")
	return int(pfdVal.FieldByName("Sysfd").Int())
}
