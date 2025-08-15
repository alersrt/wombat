package internal

type Router struct {
	requests  chan *Request
	responses chan *Response
}

func NewRouter() *Router {
	return &Router{
		requests:  make(chan *Request, 1),
		responses: make(chan *Response, 1),
	}
}

func (r *Router) SendReq(req *Request) {
	r.requests <- req
}

func (r *Router) SendRes(res *Response) {
	r.responses <- res
}

func (r *Router) ReqChan() chan *Request {
	return r.requests
}

func (r *Router) ResChan() chan *Response {
	return r.responses
}
