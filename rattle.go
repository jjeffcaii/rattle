package rattle

import (
	"context"
	"sync"

	"github.com/pkg/errors"
	"github.com/rsocket/rsocket-go"
	"go.uber.org/atomic"
)

type Config struct {
	Namespace string
	Token     string
	Bootstrap string
}

type Rattle struct {
	mu      sync.RWMutex
	c       *Config
	router  *Router
	started *atomic.Bool
	client  rsocket.CloseableRSocket
	done    chan struct{}
}

func (r *Rattle) Close() (err error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	select {
	case <-r.done:
		// already closed
		break
	default:
		close(r.done)
		if r.client == nil {
			break
		}
		err = r.client.Close()
	}
	return
}

func (r *Rattle) GET(path string, h Handler) error {
	return r.router.Route(GET, path, h)
}

func (r *Rattle) Request(domain string) *Requester {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return &Requester{
		ns: domain,
		c:  r.client,
	}
}

func (r *Rattle) Start(ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.client != nil {
		return errors.New("rattle: already started")
	}

	client, err := rsocket.Connect().
		Acceptor(func(socket rsocket.RSocket) rsocket.RSocket {
			return NewRouteResponder(r.router)
		}).
		Transport(rsocket.TCPClient().SetAddr(r.c.Bootstrap).Build()).
		Start(ctx)
	if err != nil {
		return err
	}
	r.client = client
	return nil
}

func NewRattle(c *Config) (*Rattle, error) {
	return &Rattle{
		c:      c,
		router: NewRouter(),
		done:   make(chan struct{}),
	}, nil
}
