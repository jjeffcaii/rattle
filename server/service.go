package server

import (
	"context"

	"github.com/jjeffcaii/rattle"
	"github.com/pkg/errors"
	"github.com/rsocket/rsocket-go"
	"github.com/rsocket/rsocket-go/balancer"
	"github.com/rsocket/rsocket-go/payload"
	"github.com/rsocket/rsocket-go/rx/mono"
)

var errServerNotAvailable = errors.New("server is not available")

type ServiceFactory struct {
	store *balancer.Group
}

func (s *ServiceFactory) Register(socket rsocket.Client, routing *rattle.Routing, others ...*rattle.Routing) (err error) {
	s.registerOne(socket, routing)
	for _, it := range others {
		s.registerOne(socket, it)
	}
	return
}

func (s *ServiceFactory) Request(req payload.Payload) mono.Mono {
	m, ok := req.Metadata()
	if !ok {
		return mono.Error(errBadRequest)
	}
	routing, err := rattle.ParseRouting(m)
	if err != nil {
		return mono.Error(err)
	}

	sid := s.toServiceID(routing)
	next, ok := s.store.Get(sid).Next(context.Background())
	if !ok {
		return mono.Error(errServerNotAvailable)
	}
	return next.RequestResponse(req)
}

func (s *ServiceFactory) toServiceID(routing *rattle.Routing) string {
	return routing.Domain
}

func (s *ServiceFactory) registerOne(socket rsocket.Client, routing *rattle.Routing) {
	s.store.Get(s.toServiceID(routing)).Put(socket)
}

func NewServiceFactory() *ServiceFactory {
	return &ServiceFactory{
		store: balancer.NewGroup(balancer.NewRoundRobinBalancer),
	}
}
