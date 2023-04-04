package goshrt

const (
	ErrNotFound = Error("item not found")
	ErrInvalid  = Error("invalid data")
)

type Error string

func (e Error) Error() string { return string(e) }
