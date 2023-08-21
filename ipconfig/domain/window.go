package domain

const (
	WindowSize = 5
)

type stateWindow struct {
	stateQueue []*State    // 窗口内的状态信息
	statChan   chan *State // 新状态信息
	sumStat    *State      // 汇总后的总值
	idx        int64
}

func newWindow() *stateWindow {
	return &stateWindow{
		stateQueue: make([]*State, 0),
		statChan:   make(chan *State),
		sumStat:    &State{},
		idx:        0,
	}
}

func (sw *stateWindow) getState() *State {
	// 将窗口总值平均后返回
	result := sw.sumStat.Clone()
	result.Avg(WindowSize)
	return result
}

func (sw *stateWindow) appendState(state *State) {
	// 减去淘汰的
	sw.sumStat.Sub(sw.stateQueue[sw.idx%WindowSize])
	// 加上新增的
	sw.sumStat.Add(state)
	// 放入窗口中
	sw.stateQueue[sw.idx%WindowSize] = state
	sw.idx++
}
