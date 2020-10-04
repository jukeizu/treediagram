package intent

import "github.com/jukeizu/contract"

type Query struct {
	ServerId string
	Type     string
}

func (q Query) Match(intent *contract.Intent) bool {
	if intent == nil {
		return false
	}

	serverIDMatch := (intent.ServerID == "" || intent.ServerID == q.ServerId)
	typeMatch := intent.Type == q.Type

	return serverIDMatch && typeMatch
}
