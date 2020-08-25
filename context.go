package rattle

import (
	"encoding/json"

	"github.com/jjeffcaii/rattle/internal"
	"github.com/rsocket/rsocket-go"
	"github.com/rsocket/rsocket-go/payload"
	"github.com/rsocket/rsocket-go/rx/mono"
)

type Context struct {
	*internal.PathVariables
	sink    mono.Sink
	reqData []byte
	client  rsocket.RSocket
}

func (c *Context) Request(namespace string) *Requester {
	// TODO:
	return nil
}

func (c *Context) Data() []byte {
	return c.reqData
}

func (c *Context) Bind(v interface{}) error {
	return json.Unmarshal(c.reqData, v)
}

func (c *Context) JSON(data interface{}) (err error) {
	b, err := json.Marshal(data)
	c.sink.Success(payload.New(b, nil))
	return
}

func newContext(v *internal.PathVariables, sink mono.Sink) *Context {
	return &Context{
		PathVariables: v,
		sink:          sink,
	}
}
