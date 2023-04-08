package goshrt

const (
	ErrNotFound = Error("item not found")
	ErrInvalid  = Error("invalid data")
	ErrMultiple = Error("multiple data with same domain and slug")
)

type Error string

func (e Error) Error() string { return string(e) }
