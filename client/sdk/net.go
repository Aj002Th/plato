package sdk

type connect struct {
	serverAddr         string
	sendChan, recvChan chan *Message
}

func newConnect(serverAddr string) *connect {
	return &connect{
		serverAddr: serverAddr,

		// todo: no buffer?
		sendChan: make(chan *Message),
		recvChan: make(chan *Message),
	}
}

func (c *connect) send(msg *Message) {
	// 先本地调通, 直接发给接收方
	c.recvChan <- msg
}

func (c *connect) recv() <-chan *Message {
	return c.recvChan
}

func (c *connect) close() {

}
