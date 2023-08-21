package domain

import "math"

// State gateway上报的一些指标
// 用于 ip 分配的决策
// 考虑到部署服务器的能力可能有所不同, 剩余值比当前值更能体现负载的情况
// 目前调度的评判指标和评分方式都比较简陋, 但是留出接口, 可以方便地修改
type State struct {
	// 静态分指标
	ConnectNum float64 // 业务上，im gateway 总体持有的长连接数量 的剩余值

	// 动态分指标
	MessageBytes float64 // 业务上，im gateway 每秒收发消息的总字节数(GB) 的剩余值
}

func (s *State) CalculateActiveScore() float64 {
	decimal := func(value float64) float64 {
		// 保留两位小数
		return math.Trunc(value*1e2+0.5) * 1e-2
	}
	getGB := func(m float64) float64 {
		// B -> GB
		return decimal(m / (1 << 30))
	}

	// 每秒收发消息的总字节数, 单位 GB, 保留两位小数
	return getGB(s.MessageBytes)
}

func (s *State) CalculateStaticScore() float64 {
	// 等同于 剩余连接数
	return s.ConnectNum
}

func (s *State) Clone() *State {
	newStat := &State{
		MessageBytes: s.MessageBytes,
		ConnectNum:   s.ConnectNum,
	}
	return newStat
}

// 一些对于 State 的数据操作

func (s *State) Add(st *State) {
	if st == nil {
		return
	}
	s.ConnectNum += st.ConnectNum
	s.MessageBytes += st.MessageBytes
}

func (s *State) Sub(st *State) {
	if st == nil {
		return
	}
	s.ConnectNum -= st.ConnectNum
	s.MessageBytes -= st.MessageBytes
}

func (s *State) Avg(num float64) {
	s.ConnectNum /= num
	s.MessageBytes /= num
}
