package goshrt

type Authorizer interface {
	Create(id string) (string, error)
	Validate(t string) (bool, error)
}
