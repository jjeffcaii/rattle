package rattle

import (
	"context"
	"encoding/json"

	"github.com/rsocket/rsocket-go"
	"github.com/rsocket/rsocket-go/payload"
)

type Response struct {
	p payload.Payload
}

type Requester struct {
	ns string
	c  rsocket.RSocket
}

func (r Response) Bind(v interface{}) error {
	return json.Unmarshal(r.p.Data(), v)
}

func (r *Requester) createRoutingMetadata(path string) []byte {
	// TODO:
	return nil
}

func (r *Requester) GET(ctx context.Context, path string, data interface{}) (*Response, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	res, err := r.c.RequestResponse(payload.New(b, r.createRoutingMetadata(path))).Block(ctx)
	if err != nil {
		return nil, err
	}
	return &Response{
		p: res,
	}, nil
}
