package xianliu

import (
	"time"
)

type SlidewindowLimiter struct {
	//窗口大小
	windowSize time.Duration
	//整个限流数量
	limit int
	//分割小窗口个数
	splitNum int
	//每个小窗口的计时器
	counters []int
	//指示当前小窗口，（index+1）%splitNum 是最早的窗口，及时丢弃
	index int
	//窗口开始时间
	starttime time.Time
}

func NewSlidewindowLimiter(windowSize time.Duration, limit, splitNum int) *SlidewindowLimiter {
	s := &SlidewindowLimiter{
		windowSize: windowSize,
		limit:      limit,
		splitNum:   splitNum,
		counters:   make([]int, splitNum),
		index:      0,
		starttime:  time.Now(),
	}
	return s
}

func (s *SlidewindowLimiter) slide(windowsNum int) {
	if windowsNum == 0 {
		return
	}
	slideNum := min(windowsNum, s.splitNum)
	for i := 0; i < slideNum; i++ {
		s.index = (s.index + 1) % s.splitNum
		//当窗口占满后，清空最早的一个窗口，循环数组
		s.counters[s.index] = 0
	}
	//更新窗口时间
	s.starttime.Add(time.Duration(windowsNum * (int(s.windowSize) / s.splitNum)))

}

func (s *SlidewindowLimiter) tryAcquire() bool {
	now := time.Now()
	windowNum := max(int(now.Sub(s.starttime.Add(s.windowSize))), 0)/(int(s.windowSize)/s.splitNum) + 1
	s.slide(windowNum)
	count := 0
	for _, v := range s.counters {
		count += v
	}
	if count >= s.limit {
		return false
	} else {
		s.counters[s.index]++
		return true
	}
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func max(x, y int) int {
	if x < y {
		return y
	}
	return x
}
