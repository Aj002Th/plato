package source

import (
	"fmt"
	"plato/common/discovery"
)

// Event 对应 etcd 中的信息更新所对应的事件
type Event struct {
	Type         EventType
	IP           string
	Port         string
	ConnectNum   float64
	MessageBytes float64
}

type EventType string

const (
	AddNodeEvent EventType = "addNode"
	DelNodeEvent EventType = "delNode"
)

var eventChan chan *Event

func EventChan() <-chan *Event {
	return eventChan
}

// NewEvent 将 etcd 中的信息做转换
func NewEvent(endpointInfo *discovery.EndpointInfo) *Event {
	if endpointInfo == nil || endpointInfo.MetaData == nil {
		return nil
	}
	var connNum, msgBytes float64
	if data, ok := endpointInfo.MetaData["connect_num"]; ok {
		connNum = data.(float64) // 如果出错，此处应该panic 暴露错误
	}
	if data, ok := endpointInfo.MetaData["message_bytes"]; ok {
		msgBytes = data.(float64) // 如果出错，此处应该panic 暴露错误
	}
	return &Event{
		Type:         AddNodeEvent,
		IP:           endpointInfo.IP,
		Port:         endpointInfo.Port,
		ConnectNum:   connNum,
		MessageBytes: msgBytes,
	}
}

// Key endpoint 存储在 etcd 中的 key
func (e *Event) Key() string {
	return fmt.Sprintf("%v:%v", e.IP, e.Port)
}
