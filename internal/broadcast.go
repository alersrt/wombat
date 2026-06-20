package internal

import (
	"context"

	"github.com/alersrt/wombat/pkg"
)

type BroadcastServer interface {
	Subscribe() <-chan pkg.Message
	CancelSubscription(<-chan pkg.Message)
	Serve(ctx context.Context, input <-chan pkg.Message) error
	Close() error
}

type broadcastServer struct {
	listeners      []chan pkg.Message
	addListener    chan chan pkg.Message
	removeListener chan (<-chan pkg.Message)
	processor      pkg.Component
}

func (s *broadcastServer) Subscribe() <-chan pkg.Message {
	newListener := make(chan pkg.Message)
	s.addListener <- newListener
	return newListener
}

func (s *broadcastServer) CancelSubscription(channel <-chan pkg.Message) {
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

func NewBroadcastServer(processor pkg.Component) BroadcastServer {
	service := &broadcastServer{
		listeners:      make([]chan pkg.Message, 0),
		addListener:    make(chan chan pkg.Message),
		removeListener: make(chan (<-chan pkg.Message)),
		processor:      processor,
	}
	return service
}

func (s *broadcastServer) Serve(ctx context.Context, input <-chan pkg.Message) error {
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
