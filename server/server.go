package server

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/rsocket/rsocket-go"
	"github.com/rsocket/rsocket-go/extension"
	"github.com/rsocket/rsocket-go/payload"
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

type Server struct {
	port    int
	node    *Node
	factory *ServiceFactory
}

func (r *Server) Serve(ctx context.Context) (err error) {
	err = r.node.Start()
	if err != nil {
		return
	}

	responder := &Responder{
		repo: r.factory,
	}
	t := rsocket.TCPServer().SetAddr(fmt.Sprintf(":%d", r.port)).Build()
	err = rsocket.Receive().
		Acceptor(func(setup payload.SetupPayload, s rsocket.CloseableRSocket) (rsocket.RSocket, error) {
			if setup.MetadataMimeType() != extension.MessageCompositeMetadata.String() {
				return nil, errBadSetup
			}
			return responder, nil
		}).
		Transport(t).
		Serve(ctx)
	return
}

type Matcher interface {
	Match(path string) bool
}

func NewRattle(c *Config) (r *Server, err error) {
	if c == nil {
		err = errors.New("no config found")
		return
	}
	node := NewNode(c.ID, c.Discovery.Port, c.Discovery.Seeds)
	r = &Server{
		port:    c.Port,
		node:    node,
		factory: NewServiceFactory(),
	}
	return
}
