package gateway

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"plato/common/config"
	"plato/common/logger"
	"plato/common/tcp"
)

func RunMain(path string) {
	config.Init(path) // 初始化配置
	// 建立服务端 tcp 连接
	ln, err := net.ListenTCP(
		"tcp",
		&net.TCPAddr{Port: config.GetGatewayServerPort()})
	if err != nil {
		log.Fatalf("Start Gateway Server err:%s", err.Error())
	}
	initWorkPool()
	initEpollPool(ln, runProc)
	select {}
}

// 注入 epoll 的回调处理函数
func runProc(c *connection, epoll *epoller) {
	logger.Debug("in runProc")

	// 1. 获取连接中的数据包
	dataBuf, err := tcp.ReadData(c.conn)
	if err != nil {
		// 如果读取conn时发现连接关闭，则直接端口连接
		if errors.Is(err, io.EOF) {
			_ = epoll.remove(c)
		}
		return
	}
	// 2. 交给协程池处理 - 调用 state server rpc, 业务rpc 等
	err = wPool.Submit(func() {
		// do something
		logger.Debug("write action have summit to wPool")

		// 实现 echo
		bytes := tcp.DataPkg{
			Len:  uint32(len(dataBuf)),
			Data: dataBuf,
		}
		err := tcp.SendData(c.conn, bytes.Marshal())
		if err != nil {
			logger.Debug("tcp.SendData err: ", err.Error())
		}
	})
	if err != nil {
		fmt.Errorf("runProc err: %+v\n", err.Error())
	}
}
