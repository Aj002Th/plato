package gateway

import (
	"log"
	"net"
	"plato/common/config"
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
func runProc(c *connection, ep *epoller) {}
