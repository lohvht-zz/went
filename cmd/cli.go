package cmd

import (
	"flag"
	"io/ioutil"
	"log"
	"path/filepath"
)

// NOTE: write-up on how to decouple CLI and Running commands
// https://npf.io/2016/10/reusable-commands/

// Run starts the command line process, returning an error code when the process is
// finished
func Run() int {
	// if len(os.Args) < 2 {
	// 	fmt.Println("Entering Interpreter Mode")
	// }

	filePtr := flag.String("f", "", "Script file to read and parse (Required)")
	flag.Parse()

	if *filePtr == "" {
		flag.PrintDefaults()
		return 1
	}
	// Read the entire script into file, this is how they handle it for golang's html/template: https://golang.org/src/html/template/template.go (LINE 420)
	// NOTE: If this proves to be an issue later on, use a buffer a la: https://stackoverflow.com/questions/13514184/how-can-i-read-a-whole-file-into-a-string-variable-in-golang
	// Not likely though, since our scripts are meant to be literally all text (i.e. no finicky business with images)
	// Worst case scenario would be to restrict the file extension?
	b, err := ioutil.ReadFile(*filePtr)
	if err != nil {
		log.Fatalf("Encountered error with opening/reading the file input: %s.\n", *filePtr)
		return 1
	}
	s := string(b) // string value of input
	name := filepath.Base(*filePtr)
	parseInput(name, s)
	return 0
}

// func interpreterMode() {
// REVIEW:
// }

// parseInput takes in the string input and runs the language
func parseInput(name, input string) {
	// p, err := lang.Parse(name, input)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// i := lang.NewInterpreter(p.Root)
	// i.Interpret()
}
