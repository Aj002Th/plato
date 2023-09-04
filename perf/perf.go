package perf

import (
	"net"
	"plato/common/sdk"
)

var (
	TcpConnNum int32 // 建立的连接数目
)

func RunMain() {
	// 疯狂建立连接
	for i := 0; i < int(TcpConnNum); i++ {
		sdk.NewChat(
			net.ParseIP("127.0.0.1"),
			8900,
			"logic",
			"1223", "123")
	}
}
