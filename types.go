package rattle

const (
	HEAD   Method = "HEAD"
	GET           = "GET"
	POST          = "POST"
	PUT           = "PUT"
	PATCH         = "PATCH"
	DELETE        = "DELETE"
)

type Handler func(c *Context) error

type Method string

type Routing struct {
	Domain string
	Method Method
	Path   string
}
