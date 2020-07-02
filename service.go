package rattle

import (
	"github.com/jjeffcaii/rattle/pkg"
	"github.com/rsocket/rsocket-go"
	"github.com/rsocket/rsocket-go/balancer"
	"github.com/rsocket/rsocket-go/payload"
	"github.com/rsocket/rsocket-go/rx/mono"
)

type ServiceFactory struct {
	store *balancer.Group
}

func (s *ServiceFactory) Register(socket rsocket.Client, routing *pkg.Routing, others ...*pkg.Routing) (err error) {
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
	routing, err := pkg.ParseRouting(m)
	if err != nil {
		return mono.Error(err)
	}

	sid := s.toServiceID(routing)
	next := s.store.Get(sid).Next()
	return next.RequestResponse(req)
}

func (s *ServiceFactory) toServiceID(routing *pkg.Routing) string {
	return routing.Domain
}

func (s *ServiceFactory) registerOne(socket rsocket.Client, routing *pkg.Routing) {
	s.store.Get(s.toServiceID(routing)).Put(socket)
}

func NewServiceFactory() *ServiceFactory {
	return &ServiceFactory{
		store: balancer.NewGroup(balancer.NewRoundRobinBalancer),
	}
}
