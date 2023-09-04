package sdk

import (
	"fmt"
	"github.com/cloudwego/hertz/pkg/common/json"
	"net"
	"plato/common/logger"
	"plato/common/tcp"
)

type connect struct {
	conn               *net.TCPConn
	sendChan, recvChan chan *Message
}

func newConnect(ip net.IP, port int) *connect {
	logger.Debug("new connection")
	clientConn := &connect{
		sendChan: make(chan *Message),
		recvChan: make(chan *Message),
	}
	addr := &net.TCPAddr{
		IP:   ip,
		Port: port,
	}
	// 建立连接
	tcpConn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		fmt.Printf("DialTCP.err=%+v\n", err)
		return nil
	}
	clientConn.conn = tcpConn
	// 监听连接, 接收消息
	logger.Debug("start listen: tcp conn")
	go func() {
		for {
			logger.Debug("in the for: tcp conn")
			data, err := tcp.ReadData(tcpConn)
			logger.Debug("get data: tcp conn")
			if err != nil {
				fmt.Printf("ReadData err:%+v\n", err)
			}
			msg := &Message{}
			_ = json.Unmarshal(data, msg)
			clientConn.recvChan <- msg
		}
	}()
	return clientConn
}

func (c *connect) send(msg *Message) {
	// 还是先不用上 send chan
	// 直接发回 recv 对应的 tcp 连接
	//bytes, _ := json.Marshal(msg)
	//dataPgk := tcp.DataPkg{
	//	Data: bytes,
	//	Len:  uint32(len(bytes)),
	//}
	//xx := dataPgk.Marshal()
	//_, _ = c.conn.Write(xx)
	//n, err := c.conn.Write(xx)
	//fmt.Debug("conn.write, ", n, "; ", err)

	dataBuf, _ := json.Marshal(msg)
	logger.Debug(string(dataBuf))
	bytes := tcp.DataPkg{
		Len:  uint32(len(dataBuf)),
		Data: dataBuf,
	}
	err := tcp.SendData(c.conn, bytes.Marshal())
	if err != nil {
		logger.Debug("tcp.SendData err in chat.Send: ", err.Error())
	}
}

func (c *connect) recv() <-chan *Message {
	return c.recvChan
}

func (c *connect) close() {
	_ = c.conn.Close()
}
