package parser

import "strings"

// ParseCommand parses command into command object
func ParseCommand(text string) Command {
	var parts = strings.Split(text, " ")
	if len(parts) != 2 {
		return Command{"", ""}
	}
	return Command{parts[0], parts[1]}
}

type Command struct {
	Action  string
	Payload string
}
