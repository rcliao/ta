package parser

import (
	"errors"
	"strings"
)

type Command struct {
	Action  string
	Payload string
}

// ParseCommand parses command into command object
func ParseCommand(text string) (Command, error) {
	var parts = strings.Split(text, " ")
	if len(parts) != 2 {
		return Command{"", ""}, errors.New("Cannot recognize command. Please follow format of \"Action Payload\"")
	}
	return Command{parts[0], parts[1]}, nil
}
