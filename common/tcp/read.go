package tcp

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

// ReadData 读取一个数据包
// 需要注意 tcp 的粘包拆包问题
func ReadData(conn *net.TCPConn) ([]byte, error) {
	// 先读取前 4 字节获取数据包长度
	var dataLen uint32
	headBuf := make([]byte, 4)
	if err := readFixedData(conn, headBuf); err != nil {
		return nil, err
	}
	buffer := bytes.NewBuffer(headBuf)
	if err := binary.Read(buffer, binary.BigEndian, &dataLen); err != nil {
		return nil, fmt.Errorf("read headlen error:%s", err.Error())
	}
	if dataLen <= 0 {
		return nil, fmt.Errorf("wrong headlen :%d", dataLen)
	}
	// 再获取数据包内容
	dataBuf := make([]byte, dataLen)
	if err := readFixedData(conn, dataBuf); err != nil {
		return nil, fmt.Errorf("read headlen error:%s", err.Error())
	}
	return dataBuf, nil
}

// 读取固定长度的数据
func readFixedData(conn *net.TCPConn, buf []byte) error {
	// 超时时间 120 s
	_ = conn.SetReadDeadline(time.Now().Add(time.Duration(120) * time.Second))

	totalSize := len(buf)
	size := 0
	for {
		n, err := conn.Read(buf[size:])
		if err != nil {
			return err
		}
		size += n
		if size >= totalSize {
			break
		}
	}
	return nil
}
