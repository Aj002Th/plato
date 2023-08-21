package domain

import (
	"sync/atomic"
	"unsafe"
)

type EndPoint struct {
	IP          string       `json:"ip"`
	Port        string       `json:"port"`
	ActiveScore float64      `json:"-"`
	StaticScore float64      `json:"-"`
	State       *State       `json:"-"` // 代表当前的统计结果(均值)
	windows     *stateWindow `json:"-"` // 用于统计&计算State
}

func NewEndPoint(ip, port string) *EndPoint {
	// 填充数据
	ed := &EndPoint{
		IP:   ip,
		Port: port,
	}
	ed.windows = newWindow()
	ed.State = ed.windows.sumStat

	// 启动更新worker, 监听并处理 window 中的 state update
	go func() {
		for state := range ed.windows.statChan {
			ed.windows.appendState(state)
			newState := ed.windows.getState()

			// 原子操作更新 State, 省把锁, 其他地方读取 ed.State 时也更方便
			atomic.SwapPointer((*unsafe.Pointer)((unsafe.Pointer)(ed.State)), unsafe.Pointer(newState))
		}
	}()

	return ed
}

func (ed *EndPoint) UpdateStat(s *State) {
	ed.windows.statChan <- s
}

func (ed *EndPoint) CalculateScore(ctx *IpConfigContext) {
	if ed.State == nil {
		ed.ActiveScore = ed.State.CalculateActiveScore()
		ed.StaticScore = ed.State.CalculateStaticScore()
	}
}
