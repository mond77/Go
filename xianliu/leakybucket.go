package xianliu

import (
	"time"
)

type LeakyBucketLimiter struct {
	capacity int
	//每秒rate个
	rate int
	left int
	//undefined，最好自己实现request的List
	ReqList *List
}

func NewLeakyBucketLimiter(capacity, rate int) *LeakyBucketLimiter {
	l := &LeakyBucketLimiter{capacity, rate, capacity, &List{}}
	go l.handleReq()
	return l
}

func (l *LeakyBucketLimiter) handleReq() {
	t := time.Tick(time.Second / time.Duration(l.rate))
	for range t {
		if req := l.ReqList.removeHead(); req != nil {
			go req.handle()
		}
	}
}

func (l *LeakyBucketLimiter) TryAcquire(req *Request) bool {
	if l.left <= 0 {
		return false
	}
	l.left--
	l.ReqList.addToTail(req)
	return true
}


