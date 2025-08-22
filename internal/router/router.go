package router

import (
	"wombat/internal/domain"
)

type Router struct {
	requests  chan *domain.Request
	responses chan *domain.Response
}

func NewRouter() *Router {
	return &Router{
		requests:  make(chan *domain.Request),
		responses: make(chan *domain.Response),
	}
}

func (r *Router) SendReq(req *domain.Request) {
	r.requests <- req
}

func (r *Router) SendRes(res *domain.Response) {
	r.responses <- res
}

func (r *Router) ReqChan() chan *domain.Request {
	return r.requests
}

func (r *Router) ResChan() chan *domain.Response {
	return r.responses
}
