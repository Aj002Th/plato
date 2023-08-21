package domain

import (
	"plato/ipconfig/source"
	"sort"
	"sync"
)

type Dispatcher struct {
	candidateTable map[string]*EndPoint
	sync.RWMutex
}

var dispatcher *Dispatcher

func Init() {
	dispatcher = &Dispatcher{
		candidateTable: make(map[string]*EndPoint),
	}

	// 监听 eventChan 来获得 etcd 中所有 ipconfig endpoint 的状态变化
	// 异步更新状态信息
	go func() {
		for event := range source.EventChan() {
			switch event.Type {
			case source.AddNodeEvent:
				dispatcher.addNode(event)
			case source.DelNodeEvent:
				dispatcher.delNode(event)
			}
		}
	}()
}

// Dispatch ip调度
func Dispatch(ctx *IpConfigContext) []*EndPoint {
	// 1. 获取 endpoints
	eds := dispatcher.getCandidateEndpoint(ctx)
	// 2. 逐一计算分数
	for _, ed := range eds {
		ed.CalculateScore(ctx)
	}
	// 3. 依据分数进行排序
	sort.Slice(eds, func(i, j int) bool {
		// 优先看动态分
		if eds[i].ActiveScore > eds[j].ActiveScore {
			return true
		}
		// 动态分相同比较静态分
		if eds[i].ActiveScore == eds[j].ActiveScore {
			if eds[i].StaticScore > eds[j].StaticScore {
				return true
			}
		}
		return false
	})

	return eds
}

func (dp *Dispatcher) getCandidateEndpoint(ctx *IpConfigContext) []*EndPoint {
	dp.RLock()
	defer dp.RUnlock()

	eds := make([]*EndPoint, 0, len(dp.candidateTable))
	for _, ed := range dp.candidateTable {
		eds = append(eds, ed)
	}
	return eds
}

func (dp *Dispatcher) addNode(event *source.Event) {
	dp.Lock()
	defer dp.Unlock()

	ed, ok := dp.candidateTable[event.Key()]
	if !ok {
		ed = NewEndPoint(event.IP, event.Port)
		dp.candidateTable[event.Key()] = ed
	}
	ed.UpdateStat(&State{
		ConnectNum:   event.ConnectNum,
		MessageBytes: event.MessageBytes,
	})
}

func (dp *Dispatcher) delNode(event *source.Event) {
	dp.Lock()
	defer dp.Unlock()
	delete(dp.candidateTable, event.Key())
}
