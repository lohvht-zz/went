package cmd

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/lohvht/went/lang/token"
)

// NOTE: write-up on how to decouple CLI and Running commands
// https://npf.io/2016/10/reusable-commands/

var usageReminder = "Usage: ./went [script]"

// Run starts the command line process, returning an error code when the process is
// finished
func Run() int {
	if len(os.Args) > 2 {
		log.Fatalln(usageReminder)
	} else if len(os.Args) == 2 {
		filename := os.Args[1]
		if filename == "" {
			log.Fatalln(usageReminder)
		}
		b, err := ioutil.ReadFile(filename)
		if err != nil {
			log.Fatalf("Encountered error with opening/reading the file input: %s.\n", filename)
			return 1
		}
		s := string(b) // string value of input
		name := filepath.Base(filename)
		runFile(name, s)
	} else {
		runPrompt()
	}
	return 0
}

// runPrompt starts a went prompt session
func runPrompt() {
	// REVIEW: Make a mode that runs line-by-line interpretation in a manner similar
	// to Python IDLE or javascript consoles for browsers
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		fmt.Print("> ")
		run(token.Tokenise("Interpreter Mode> ", s.Text()))
		hasError = false
	}
}

// runFile takes in the string input and runs the language
func runFile(name, input string) {
	// p, errp := lang.Parse(name, input)
	// if errp != nil {
	// 	log.Fatal(errp)
	// }
	// _, erri := lang.Interpret(p.Root)
	// if erri != nil {
	// 	log.Fatal(erri)
	// }
	lexer := token.Tokenise(name, input)
	run(lexer)

	if hasError {
		os.Exit(65)
	}
}

func run(lexer *token.Lexer) {
	for {
		t := lexer.Next()
		if t.Type == token.EOF {
			fmt.Println(t)
			break
		}
		fmt.Println(t)
	}
}

//////////////////////////////////////////
// EXTRA
var hasError = false

func error(line int, msg string) { report(line, "", msg) }

func report(line int, where, msg string) {
	log.Fatalf("[line %d] Error %s: %s", line, where, msg)
	hasError = true
}
