package pkg

type Method string

const (
	HEAD   Method = "HEAD"
	GET           = "GET"
	POST          = "POST"
	PUT           = "PUT"
	PATCH         = "PATCH"
	DELETE        = "DELETE"
)

type Routing struct {
	Domain string
	Method Method
	Path   string
}

