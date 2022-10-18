package xianliu

import (
	"fmt"
	"sync/atomic"
	"time"
)

type CounterLimiter struct {
	windowSize time.Duration
	limit      int64
	count      int64
}

func NewCounterLimiter(w time.Duration, limit int64) *CounterLimiter {
	c := &CounterLimiter{w, limit, 0}
	//注意要开一个goroutine执行定时器
	go c.run()
	return c
}

func (c *CounterLimiter) run() {
	t := time.Tick(c.windowSize)
	for range t {
		atomic.StoreInt64(&(c.count), 0)
	}
}

func (c *CounterLimiter) tryAcquire() bool {
	newcount := atomic.AddInt64(&(c.count), 1)
	return newcount <= c.limit
}

func T1() {
	c := NewCounterLimiter(1*time.Second, 20)
	count := 0
	for i := 0; i < 50; i++ {
		if c.tryAcquire() {
			count++
		}
	}
	fmt.Printf("第一拨50次请求中通过: %v,限流：%v \n" ,count, (50 - count))
	//过一秒再请求
	time.Sleep(time.Second)
	//模拟50次请求，看多少能通过
	count = 0
	for i := 0; i < 50; i++ {
		if c.tryAcquire() {
			count++
		}
	}
	fmt.Printf("第一拨50次请求中通过: %v,限流：%v \n" ,count, (50 - count))
}
