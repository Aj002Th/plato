package gateway

import "net"

type connection struct {
	fd   int
	conn *net.TCPConn
}

func (c connection) Close() {
	err := c.conn.Close()
	panic(err)
}

// RemoteAddr 获取连接地址
func (c connection) RemoteAddr() string {
	return c.conn.RemoteAddr().String()
}
