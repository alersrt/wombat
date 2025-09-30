package broadcast

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/alersrt/wombat/pkg"
)

type BroadcastServer interface {
	Subscribe() <-chan []byte
	CancelSubscription(<-chan []byte)
	Serve(context.Context)
	// Close() error
}

type broadcastServer struct {
	producer       pkg.Producer
	listeners      []chan []byte
	addListener    chan chan []byte
	removeListener chan (<-chan []byte)
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
	_ = s.producer.Close()
	return nil
}

func NewBroadcastServer(producer pkg.Producer) BroadcastServer {
	service := &broadcastServer{
		producer:       producer,
		listeners:      make([]chan []byte, 0),
		addListener:    make(chan chan []byte),
		removeListener: make(chan (<-chan []byte)),
	}
	return service
}

func (s *broadcastServer) Serve(ctx context.Context) {
	defer func() {
		_ = s.Close()
	}()

	ctx, cancel := context.WithCancelCause(ctx)

	go func() {
		err := s.producer.Run(ctx)
		if err != nil {
			slog.Error(fmt.Sprintf("broadcast: serve: %+v", err))
			cancel(err)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return
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
		case val, ok := <-s.producer.Publish():
			if !ok {
				return
			}
			for _, listener := range s.listeners {
				if listener != nil {
					select {
					case listener <- val:
					case <-ctx.Done():
						return
					}
				}
			}
		}
	}
}
