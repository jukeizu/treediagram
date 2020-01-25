package builtin

import (
	"github.com/jukeizu/contract"
)

type HandlerRegistration struct {
	Name    string
	Handler func(contract.Request) (*contract.Response, error)
}

type Handler interface {
	Registrations() []HandlerRegistration
}
