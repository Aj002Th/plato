package sdk

const (
	MsgTypeText = "text"
)

type Message struct {
	Type       string
	Name       string
	FromUserID string
	ToUserID   string
	Content    string
	SessionID  string
}

type Chat struct {
	Nick      string
	UserID    string
	SessionID string
	conn      *connect
}

func NewChat(serverAddr, nick, userID, sessionID string) *Chat {
	return &Chat{
		Nick:      nick,
		UserID:    userID,
		SessionID: sessionID,
		conn:      newConnect(serverAddr),
	}
}

// Send 发消息
func (c *Chat) Send(msg *Message) {
	c.conn.send(msg)
}

// Recv 接收消息
func (c *Chat) Recv() <-chan *Message {
	return c.conn.recv()
}

// Close 关闭连接
func (c *Chat) Close() {
	c.conn.close()
}
