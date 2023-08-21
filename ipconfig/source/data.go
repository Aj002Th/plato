package source

import (
	"context"
	"log"
	"plato/common/config"
	"plato/common/discovery"
)

func Init() {
	eventChan = make(chan *Event)
	ctx := context.Background()

	// 服务发现
	go DataHandle(&ctx)

	// mock 服务注册, 提供一些 gateway 上报的模拟数据
	if config.IsDebug() {
		ctx := context.Background()
		testServiceRegister(&ctx, "7896", "node1")
		testServiceRegister(&ctx, "7897", "node2")
		testServiceRegister(&ctx, "7898", "node3")
	}
}

// DataHandle 初始化并启动服务发现
func DataHandle(ctx *context.Context) {
	dis := discovery.NewServiceDiscovery(ctx)

	// set & del 的工作是将 etcd 中的信息做转换后放入 eventChan
	// 等待 domain 层做处理
	set := func(k, v string) {
		endpointInfo, err := discovery.UnMarshal([]byte(v))
		if err != nil {
			log.Printf("DataHandler.setFunc.err :%s\n", err.Error())
			return
		}
		event := NewEvent(endpointInfo)
		if event == nil {
			log.Printf("DataHandler.setFunc.err :%s\n", "event is nil")
			return
		}
		event.Type = AddNodeEvent
		eventChan <- event
	}
	del := func(k, v string) {
		endpointInfo, err := discovery.UnMarshal([]byte(v))
		if err != nil {
			log.Printf("DataHandler.delFunc.err :%s\n", err.Error())
			return
		}
		event := NewEvent(endpointInfo)
		if event == nil {
			log.Printf("DataHandler.delFunc.err :%s\n", "event is nil")
			return
		}
		event.Type = DelNodeEvent
		eventChan <- event
	}

	err := dis.Watch(config.GetServicePathForIPConf(), set, del)
	if err != nil {
		panic(err)
	}
}
