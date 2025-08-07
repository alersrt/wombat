package internal

type Router struct {
	requests  chan *Request
	responses chan *Response
}

func NewRouter() *Router {
	return &Router{
		requests:  make(chan *Request),
		responses: make(chan *Response),
	}
}

func (r *Router) SendReq(req *Request) {
	r.requests <- req
}

func (r *Router) SendRes(res *Response) {
	r.responses <- res
}

func (r *Router) GetReq() chan *Request {
	return r.requests
}

func (r *Router) GetRes() chan *Response {
	return r.responses
}
