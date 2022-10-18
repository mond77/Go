package xianliu

import "time"

type TokenBucketLimiter struct {
	cap    int
	//每秒rate个
	rate   int
	//令牌数量
	amount int
}

func NewTokenBucketLimiter(cap, rate int) *TokenBucketLimiter {
	t := &TokenBucketLimiter{cap, rate, 0}
	go t.addToken()
	return t
}

func (t *TokenBucketLimiter) addToken() {
	ticker := time.Tick(time.Second/time.Duration(t.rate))
	for range ticker{
		if t.amount<t.cap{
			t.amount++
		}
	}
}

//只有leakybucket算法不立即处理请求，其他限流算法实现tryAcquire时可以选择把req.handle()耦和
func (t *TokenBucketLimiter) tryAcquire(req *Request)bool{
	if t.amount<=0{
		return false
	}
	t.amount--
	req.handle()
	return true
}
