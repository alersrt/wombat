package pkg

type Router[REQ, RES any] struct {
	requests  chan *REQ
	responses chan *RES
}

func NewRouter[REQ, RES any]() *Router[REQ, RES] {
	return &Router[REQ, RES]{
		requests:  make(chan *REQ),
		responses: make(chan *RES),
	}
}

func (r *Router[REQ, RES]) SendReq(req *REQ) {
	r.requests <- req
}

func (r *Router[REQ, RES]) SendRes(res *RES) {
	r.responses <- res
}

func (r *Router[REQ, RES]) ReqChan() chan *REQ {
	return r.requests
}

func (r *Router[REQ, RES]) ResChan() chan *RES {
	return r.responses
}
