package cmd

import (
	"fmt"

	prompt "github.com/c-bata/go-prompt"
	"github.com/lohvht/went/lang/runtime"
)

var promptState struct {
	LivePrefix          string
	LivePrefixIsEnabled bool
	brackets            bracketStack
}

var wentprefix = "went> "
var multiprefix = "..... "

var query = ""

var brackets = map[string]string{
	"(": ")",
	"{": "}",
	"[": "]",
}

// bracketStack is a string slice used to collect brackets
type bracketStack []string

func (s *bracketStack) empty() bool { return len(*s) == 0 }

// push a string to the top of the bracketStack
func (s *bracketStack) push(r string) { *s = append(*s, r) }

// pop removes a string from the top of the bracketStack, you should always check if
// the bracketStack is empty prior to popping
func (s *bracketStack) pop() (r string) {
	r, *s = (*s)[len(*s)-1], (*s)[:len(*s)-1]
	return
}

// peek looks at the top of the bracketStack you should always check if the bracketStack is
// empty prior to peeking
func (s *bracketStack) peek() string { return (*s)[len(*s)-1] }

type bracketLineStatus int

const (
	normal   bracketLineStatus = iota // no multiline needed due to brackets
	open                              // bracket stack still has open brackets that needs to be closed in subsequent lines, but no error
	errbrack                          // mismatched brackets, terminate do not
)

// collectBrackets traverses the input string and keeps track of the brackets seen
// returns statuses after scanning the input
func (s *bracketStack) collectBrackets(in string) bracketLineStatus {
	for _, r := range in {
		switch rStr := string(r); rStr {
		case "(", "[", "{": // If its an opening (left) bracket
			s.push(rStr)
		case ")", "]", "}":
			if s.empty() {
				return errbrack
			}
			// compare the expected closing bracket, v with rStr
			if v, ok := brackets[s.pop()]; ok && v != rStr {
				// prematurely terminate and return false
				return errbrack
			}
		}
	}
	if s.empty() {
		return normal
	}
	return open
}

func interpretExecutor(interpreter *runtime.Interpreter) func(string) {
	return func(in string) {
		// fmt.Println("l1: ", in)
		status := promptState.brackets.collectBrackets(in)
		query += in + "\n"
		switch status {
		case open:
			promptState.LivePrefix = multiprefix
			promptState.LivePrefixIsEnabled = true
		case errbrack:
			promptState.brackets = nil // empty bracket stack
			fallthrough
		case normal:
			// fmt.Printf("\nLin1e:\n\"%s\"\n\n", query)
			runOnce(query, in, interpreter)
			// clear the query
			query = ""
		}
	}
}

func changeLivePrefix() (string, bool) {
	return promptState.LivePrefix, promptState.LivePrefixIsEnabled
}

func completer(in prompt.Document) []prompt.Suggest {
	s := []prompt.Suggest{
		// {Text: "users", Description: "Store the username and age"},
	}
	return prompt.FilterHasPrefix(s, in.GetWordBeforeCursor(), true)
}

func runOnce(query, in string, interpreter *runtime.Interpreter) {
	promptState.LivePrefixIsEnabled = false
	promptState.LivePrefix = in
	err := run("", query, interpreter)
	if err != nil {
		fmt.Println(err.Error())
	}
}
