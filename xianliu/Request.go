package xianliu


type List struct{
	Head *Request
	Tail *Request
}

type Request struct{
	Next *Request
}

func (l *List)addToTail(req *Request){
	l.Tail.Next = req
	l.Tail = req
}

func (l *List)removeHead()*Request{
	req := l.Head
	if req != nil{
		l.Head = l.Head.Next
	}
	return req
}
func (req *Request)handle() {}
