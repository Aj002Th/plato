package tcp

import (
	"net"
	"plato/common/logger"
)

// SendData 发送数据到指定连接中
// 注意拆包问题
func SendData(conn *net.TCPConn, data []byte) error {
	totalLen := len(data)
	writeLen := 0
	for {
		logger.Debug("doing write: writeLen=", writeLen)
		n, err := conn.Write(data[writeLen:])
		if err != nil {
			return err
		}
		writeLen += n
		if writeLen >= totalLen {
			break
		}
	}
	logger.Debug("finish write: n=", totalLen)
	return nil
}
