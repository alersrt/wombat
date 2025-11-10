package internal

import (
	"context"
	"wombat/internal/cel"

	"github.com/alersrt/wombat/pkg"
)

type BroadcastServer interface {
	Subscribe() <-chan []byte
	CancelSubscription(<-chan []byte)
	Serve(ctx context.Context, input <-chan []byte) error
	Close() error
}

type broadcastServer struct {
	listeners      []chan []byte
	addListener    chan chan []byte
	removeListener chan (<-chan []byte)
	processor      pkg.Processor
	filter         *cel.Cel
	transform      *cel.Cel
}

func (s *broadcastServer) Subscribe() <-chan []byte {
	newListener := make(chan []byte)
	s.addListener <- newListener
	return newListener
}

func (s *broadcastServer) CancelSubscription(channel <-chan []byte) {
	s.removeListener <- channel
}

func (s *broadcastServer) Close() error {
	for _, listener := range s.listeners {
		if listener != nil {
			close(listener)
		}
	}
	return s.processor.Close()
}

func NewBroadcastServer(processor pkg.Processor) BroadcastServer {
	service := &broadcastServer{
		listeners:      make([]chan []byte, 0),
		addListener:    make(chan chan []byte),
		removeListener: make(chan (<-chan []byte)),
		processor:      processor,
	}
	return service
}

func (s *broadcastServer) Serve(ctx context.Context, input <-chan []byte) error {
	defer func() {
		_ = s.Close()
	}()
	ctx, cancel := context.WithCancelCause(ctx)
	defer cancel(nil)

	in, err := s.processor.Process(ctx, input)
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case newListener := <-s.addListener:
			s.listeners = append(s.listeners, newListener)
		case listenerToRemove := <-s.removeListener:
			for i, ch := range s.listeners {
				if ch == listenerToRemove {
					s.listeners[i] = s.listeners[len(s.listeners)-1]
					s.listeners = s.listeners[:len(s.listeners)-1]
					close(ch)
					break
				}
			}
		case val, ok := <-in:
			if !ok {
				return nil
			}
			for _, listener := range s.listeners {
				if listener != nil {
					select {
					case listener <- val:
					case <-ctx.Done():
						return nil
					}
				}
			}
		}
	}
}
