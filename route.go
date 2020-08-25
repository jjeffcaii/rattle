package rattle

import (
	"context"
	"sync"

	"github.com/jjeffcaii/rattle/internal"
	"github.com/rsocket/rsocket-go/payload"
	"github.com/rsocket/rsocket-go/rx"
	"github.com/rsocket/rsocket-go/rx/flux"
	"github.com/rsocket/rsocket-go/rx/mono"
)

type Router struct {
	mu sync.RWMutex
	m  map[Method]*internal.PathTrie
}

type RouteResponder struct {
	router *Router
}

func (r *Router) Route(method Method, pattern string, h Handler) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	t, ok := r.m[method]
	if !ok {
		t = internal.NewPathTrie()
		r.m[method] = t
	}
	return t.AddPath(pattern, h)
}

func (r *Router) Get(method Method, path string) (*internal.PathVariables, Handler, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.m[method]
	if !ok {
		return nil, nil, false
	}
	params, v, ok := t.Find(path)
	if !ok {
		return nil, nil, false
	}
	return params, v.(Handler), true
}

func (r RouteResponder) FireAndForget(message payload.Payload) {

}

func (r RouteResponder) MetadataPush(message payload.Payload) {
	panic("implement me")
}

func (r RouteResponder) RequestResponse(request payload.Payload) mono.Mono {
	metadata, ok := request.Metadata()
	if !ok {
		return mono.Error(errInvalidRoutingFormat)
	}
	routing, err := ParseRouting(metadata)
	if err != nil {
		return mono.Error(err)
	}
	params, h, ok := r.router.Get(routing.Method, routing.Path)
	if !ok {
		return mono.Error(errNoRouting)
	}
	return mono.Create(func(ctx context.Context, sink mono.Sink) {
		c := newContext(params, sink)
		if err := h(c); err != nil {
			sink.Error(err)
		}
	})
}

func (r RouteResponder) RequestStream(message payload.Payload) flux.Flux {
	panic("implement me")
}

func (r RouteResponder) RequestChannel(messages rx.Publisher) flux.Flux {
	panic("implement me")
}

func NewRouter() *Router {
	return &Router{
		m: make(map[Method]*internal.PathTrie),
	}
}

func NewRouteResponder(router *Router) *RouteResponder {
	return &RouteResponder{
		router: router,
	}
}
