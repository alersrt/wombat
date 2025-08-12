package pkg

type Request[T any] struct {
	value *T
}

func NewReq[T any](value *T) *Request[T] {
	return &Request[T]{value}
}

func (req *Request[T]) Value() *T {
	return req.value
}

type Response[T any] struct {
	value *T
}

func NewRes[T any](value *T) *Response[T] {
	return &Response[T]{value}
}

func (res *Response[T]) Value() *T {
	return res.value
}

type Router[REQ, RES any] struct {
	requests  chan *Request[REQ]
	responses chan *Response[RES]
}

func NewRouter[REQ, RES any]() *Router[REQ, RES] {
	return &Router[REQ, RES]{
		requests:  make(chan *Request[REQ]),
		responses: make(chan *Response[RES]),
	}
}

func (r *Router[REQ, RES]) SendReq(req *Request[REQ]) {
	r.requests <- req
}

func (r *Router[REQ, RES]) SendRes(res *Response[RES]) {
	r.responses <- res
}

func (r *Router[REQ, RES]) ReqChan() chan *Request[REQ] {
	return r.requests
}

func (r *Router[REQ, RES]) ResChan() chan *Response[RES] {
	return r.responses
}
