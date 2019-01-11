package cmd

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/lohvht/went/lang/ast"
	"github.com/lohvht/went/lang/parser"
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
		err = runFile(name, s)
		if err != nil {
			log.Printf(err.Error())
			return 65
		}
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
	var err error
	fmt.Print("> ")
	for s.Scan() {
		err = run("", s.Text())
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Print("> ")
	}
}

// runFile takes in the string input and runs the language
func runFile(name, input string) error { return run(name, input) }

func run(name, input string) error {
	p := parser.New(name, input)
	expr, errs := p.Run()
	if errs != nil {
		return errs
	}
	printer := &ast.AstPrinter{}
	fmt.Println(printer.Print(expr))
	return nil
}
