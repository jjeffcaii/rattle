package rattle

import (
	"context"
	"fmt"

	"github.com/jjeffcaii/rattle/internal"
	"github.com/pkg/errors"
	"github.com/rsocket/rsocket-go"
	"github.com/rsocket/rsocket-go/extension"
	"github.com/rsocket/rsocket-go/payload"
	"github.com/rsocket/rsocket-go/rx"
	"github.com/rsocket/rsocket-go/rx/flux"
	"github.com/rsocket/rsocket-go/rx/mono"
)

var (
	errBadSetup   = errors.New("bad setup")
	errBadRequest = errors.New("bad request")
	errNoProvider = errors.New("no such service")
)

type Config struct {
	ID        string
	Port      int
	Discovery struct {
		Seeds []string
		Port  int
	}
}

type Rattle struct {
	port    int
	node    *internal.Node
	factory *ServiceFactory
}

type Responder struct {
	repo *ServiceFactory
}

func (r *Responder) FireAndForget(message payload.Payload) {
	panic("implement me")
}

func (r *Responder) MetadataPush(message payload.Payload) {
	panic("implement me")
}

func (r *Responder) RequestResponse(message payload.Payload) mono.Mono {
	return r.repo.Request(message)
}

func (r *Responder) RequestStream(message payload.Payload) flux.Flux {
	panic("implement me")
}

func (r *Responder) RequestChannel(messages rx.Publisher) flux.Flux {
	panic("implement me")
}

func (r *Rattle) Serve(ctx context.Context) (err error) {
	err = r.node.Start()
	if err != nil {
		return
	}

	responder := &Responder{
		repo: r.factory,
	}
	err = rsocket.Receive().
		Acceptor(func(setup payload.SetupPayload, sendingSocket rsocket.CloseableRSocket) (rsocket.RSocket, error) {
			if setup.MetadataMimeType() != extension.MessageCompositeMetadata.String() {
				return nil, errBadSetup
			}
			return responder, nil
		}).
		Transport(fmt.Sprintf("tcp://0.0.0.0:%d", r.port)).
		Serve(ctx)
	return
}

type Matcher interface {
	Match(path string) bool
}

func NewRattle(c *Config) (r *Rattle, err error) {
	if c == nil {
		err = errors.New("no config found")
		return
	}
	node := internal.NewNode(c.ID, c.Discovery.Port, c.Discovery.Seeds)
	r = &Rattle{
		port:    c.Port,
		node:    node,
		factory: NewServiceFactory(),
	}
	return
}
