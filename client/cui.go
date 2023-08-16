package client

import (
	"errors"
	"fmt"
	"github.com/jroimartin/gocui"
	"io/ioutil"
	"log"
	"plato/client/sdk"
)

func init() {
}

var (
	chat *sdk.Chat
	buf  string // 记录 output 的信息打印到日志

	pos int // 控制 pasteUp 和 pasteDown
)

const (
	serverAddr = "127.0.0.1:8080"
	nick       = "heng"
	userID     = "123456"
	sessionID  = "1"

	echoUserID = "654321"
	echoNick   = "simple"
)

// VOT view to terminal
type VOT struct {
	Name, Msg, Sep string
}

func (vot VOT) Show(g *gocui.Gui) error {
	v, err := g.View("output")
	if err != nil {
		return err
	}

	// 展示信息
	_, _ = fmt.Fprintf(v, "%v:%v%v\n", vot.Name, vot.Sep, vot.Msg)
	return nil
}

// 打印新消息
func viewPrint(g *gocui.Gui, name, msg string, newline bool) {
	var out VOT
	out.Name = name
	out.Msg = msg
	if newline {
		out.Sep = "\n"
	} else {
		out.Sep = " "
	}
	g.Update(out.Show)
}

// 接收消息并打印
func doRecv(g *gocui.Gui) {
	recvChan := chat.Recv()
	for msg := range recvChan {
		switch msg.Type {
		case sdk.MsgTypeText:
			viewPrint(g, msg.Name, msg.Content, false)
		}
	}
}

// 发送消息并打印
// cv 是 input 的 view
func doSay(g *gocui.Gui, cv *gocui.View) {
	v, err := g.View("output")
	if err != nil || cv == nil {
		return
	}

	// 获取 input 中的输入
	//p := make([]byte, 300)
	//n, _ := cv.Read(p)
	//if n > 0 {
	//	p = p[:n]
	//} else {
	//	return
	//}
	p := cv.Buffer()

	// 发 echo 消息
	// 构造以别人的视角发同样内容的消息
	msg := &sdk.Message{
		Type:       sdk.MsgTypeText,
		Name:       echoNick,
		FromUserID: echoUserID,
		ToUserID:   userID,
		Content:    string(p),
		SessionID:  sessionID,
	}
	// 先展示自己发的消息
	viewPrint(g, "me", string(p), false)
	// 发出消息
	chat.Send(msg)

	v.Autoscroll = true
}

// gui 样式
func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if err := viewHead(g, 1, 1, maxX-1, 3); err != nil {
		return err
	}
	if err := viewOutput(g, 1, 4, maxX-1, maxY-4); err != nil {
		return err
	}
	if err := viewInput(g, 1, maxY-3, maxX-1, maxY-1); err != nil {
		return err
	}
	return nil
}
func viewHead(g *gocui.Gui, x0, y0, x1, y1 int) error {
	if v, err := g.SetView("head", x0, y0, x1, y1); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Wrap = false
		v.Overwrite = true
		setHeadText(g, "start to chat")
	}
	return nil
}
func setHeadText(g *gocui.Gui, msg string) {
	v, err := g.View("head")
	if err != nil {
		return
	}
	v.Clear()
	_, _ = fmt.Fprintln(v, msg)
}
func viewInput(g *gocui.Gui, x0, y0, x1, y1 int) error {
	if v, err := g.SetView("input", x0, y0, x1, y1); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Wrap = false
		v.Overwrite = true

		// 可输入, 光标聚焦
		v.Editable = true
		if _, err := g.SetCurrentView("input"); err != nil {
			return err
		}
	}
	return nil
}
func viewOutput(g *gocui.Gui, x0, y0, x1, y1 int) error {
	v, err := g.SetView("output", x0, y0, x1, y1)
	if err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Wrap = true
		v.Overwrite = false
		v.Autoscroll = true
		v.SelBgColor = gocui.ColorRed
		v.Title = "Messages"
	}
	return nil
}

// gui 键位绑定
func keyBind(g *gocui.Gui) error {
	if err := g.SetKeybinding("input", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("input", gocui.KeyEnter, gocui.ModNone, viewUpdate); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("input", gocui.KeyPgup, gocui.ModNone, viewUpScroll); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("input", gocui.KeyPgdn, gocui.ModNone, viewDownScroll); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("input", gocui.KeyArrowDown, gocui.ModNone, pasteDown); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("input", gocui.KeyArrowUp, gocui.ModNone, pasteUP); err != nil {
		log.Panicln(err)
	}
	return nil
}
func quit(g *gocui.Gui, cv *gocui.View) error {
	// 关闭 chat
	chat.Close()
	// 获取 output 内容
	v, _ := g.View("output")
	buf = v.Buffer()
	// 关闭 gui
	return gocui.ErrQuit
}
func viewUpdate(g *gocui.Gui, cv *gocui.View) error {
	doSay(g, cv)
	l := len(cv.Buffer())
	cv.MoveCursor(-l, 0, true)
	cv.Clear()
	return nil
}
func viewUpScroll(*gocui.Gui, *gocui.View) error {
	return nil
}
func viewDownScroll(*gocui.Gui, *gocui.View) error {
	return nil
}
func pasteDown(*gocui.Gui, *gocui.View) error {
	return nil
}
func pasteUP(*gocui.Gui, *gocui.View) error {
	return nil
}

func RunMain() {
	// 初始化核心对象 chat
	chat = sdk.NewChat(serverAddr, nick, userID, sessionID)

	// 创建 gui, 配置样式和键位绑定
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Cursor = true
	//g.Mouse = true
	g.ASCII = true // 兼容 windows terminal

	// 设置编排函数
	g.SetManagerFunc(layout)
	// 注册回调事件
	if err := keyBind(g); err != nil {
		log.Panicln(err)
	}

	// 启动消费协程
	// quit 时 channel 关闭, goroutine 退出
	go doRecv(g)

	// 启动 gui 主循环
	if err := g.MainLoop(); err != nil && !errors.Is(err, gocui.ErrQuit) {
		log.Panicln(err)
	}

	// 记录聊天信息
	_ = ioutil.WriteFile("chat.log", []byte(buf), 0644)
}
