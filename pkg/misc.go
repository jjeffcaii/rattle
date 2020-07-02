package pkg

import (
	"github.com/pkg/errors"
	"github.com/rsocket/rsocket-go/extension"
)

var (
	errNoRouting            = errors.New("no routing found")
	errInvalidRoutingFormat = errors.New("invalid routing format")
)

func ParseRouting(metadata []byte) (routing *Routing, err error) {
	scanner := extension.NewCompositeMetadataBytes(metadata).Scanner()
	var (
		m    string
		b    []byte
		tags []string
	)
	for scanner.Scan() {
		m, b, err = scanner.Metadata()
		if err != nil {
			return
		}
		switch m {
		case extension.MessageRouting.String():
			tags, err = extension.ParseRoutingTags(b)
			if err != nil {
				return
			}
			if len(tags) != 3 {
				err = errInvalidRoutingFormat
				return
			}
			routing = &Routing{
				Domain: tags[0],
				Method: Method(tags[1]),
				Path:   tags[1],
			}
			return
		}
	}
	err = errNoRouting
	return
}
