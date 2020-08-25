package server

import (
	"github.com/rsocket/rsocket-go/payload"
	"github.com/rsocket/rsocket-go/rx"
	"github.com/rsocket/rsocket-go/rx/flux"
	"github.com/rsocket/rsocket-go/rx/mono"
)

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
