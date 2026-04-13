package todo

import (
	"fmt"
	"strconv"
	"strings"
)

type AddCmd struct {
	Title       string
	Description string
	Priority    string
}

type ListCmd struct {
	Filter string // "done", "pending", or "" for all
}

type DeleteCmd struct {
	Id int
}

type UpdateCmd struct {
	Id          int
	Title       string
	Description string
	Priority    string
}

type DoneCmd struct {
	Id int
}

type Command struct {
	Add    *AddCmd
	List   *ListCmd
	Delete *DeleteCmd
	Update *UpdateCmd
	Done   *DoneCmd
	Help   bool
	Exit   bool
}

// Tokenize splits an input line into tokens, keeping quoted strings together.
func Tokenize(line string) []string {
	var tokens []string
	var current strings.Builder
	inQuote := false

	for _, r := range line {
		switch {
		case r == '"':
			inQuote = !inQuote
		case r == ' ' && !inQuote:
			if current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
			}
		default:
			current.WriteRune(r)
		}
	}
	if current.Len() > 0 {
		tokens = append(tokens, current.String())
	}
	return tokens
}

func GetFlag(tokens []string, flag string) string {
	for i, t := range tokens {
		if t == flag && i+1 < len(tokens) {
			return tokens[i+1]
		}
	}
	return ""
}

// ParseCommand parses a raw input line into a structured Command.
// Returns an error if the input is invalid.
func ParseCommand(line string) (Command, error) {
	tokens := Tokenize(line)
	if len(tokens) == 0 {
		return Command{}, fmt.Errorf("empty input")
	}

	switch tokens[0] {
	case "add":
		if len(tokens) < 2 {
			return Command{}, fmt.Errorf("usage: add <title> [--desc <description>] [--priority <low|medium|high>]")
		}
		return Command{Add: &AddCmd{
			Title:       tokens[1],
			Description: GetFlag(tokens, "--desc"),
			Priority:    GetFlag(tokens, "--priority"),
		}}, nil

	case "list":
		filter := ""
		if len(tokens) >= 2 {
			filter = tokens[1]
		}
		return Command{List: &ListCmd{Filter: filter}}, nil

	case "delete":
		if len(tokens) < 2 {
			return Command{}, fmt.Errorf("usage: delete <id>")
		}
		id, err := strconv.Atoi(tokens[1])
		if err != nil {
			return Command{}, fmt.Errorf("invalid task id: %s", tokens[1])
		}
		return Command{Delete: &DeleteCmd{Id: id}}, nil

	case "update":
		if len(tokens) < 2 {
			return Command{}, fmt.Errorf("usage: update <id> [--title <title>] [--desc <description>] [--priority <low|medium|high>]")
		}
		id, err := strconv.Atoi(tokens[1])
		if err != nil {
			return Command{}, fmt.Errorf("invalid task id: %s", tokens[1])
		}
		return Command{Update: &UpdateCmd{
			Id:          id,
			Title:       GetFlag(tokens, "--title"),
			Description: GetFlag(tokens, "--desc"),
			Priority:    GetFlag(tokens, "--priority"),
		}}, nil

	case "done":
		if len(tokens) < 2 {
			return Command{}, fmt.Errorf("usage: done <id>")
		}
		id, err := strconv.Atoi(tokens[1])
		if err != nil {
			return Command{}, fmt.Errorf("invalid task id: %s", tokens[1])
		}
		return Command{Done: &DoneCmd{Id: id}}, nil

	case "help":
		return Command{Help: true}, nil

	case "exit", "quit":
		return Command{Exit: true}, nil

	default:
		return Command{}, fmt.Errorf("unknown command: %s", tokens[0])
	}
}
