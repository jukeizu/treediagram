package command

import (
	"regexp"
	"strings"
)

type Time interface {
	Command
}

type time struct{}

func NewTimeCommand() Time {
	return &time{}
}

func (c *time) Handle(request Request) Response {
	response := buildResponse(request)

	response.Content = "Time!"

	return response
}

func (c *time) IsCommand(request Request) bool {
	words := strings.Fields(request.Content)

	if len(words) < 1 {
		return false
	}

	lastWord := words[(len(words))-1]

	re := regexp.MustCompile("^(?i)(time)(?-i)[[:punct:]]*$")

	found := re.FindString(lastWord)

	return len(found) > 0
}
