package utils

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func main() {
	filePtr := flag.String("file", "", "Script file to read and parse (Required)")
	flag.Parse()

	if *filePtr == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	// Read the entire script into file, this is how they handle it for golang's html/template: https://golang.org/src/html/template/template.go (LINE 420)
	// (If this proves to be an issue later on, use a buffer a la: https://stackoverflow.com/questions/13514184/how-can-i-read-a-whole-file-into-a-string-variable-in-golang)
	// Not likely though, since our scripts are meant to be literally all text (i.e. no finicky business with images)
	// Worst case scenario would be to restrict the file extension?
	b, err := ioutil.ReadFile(*filePtr)
	if err != nil {
		log.Fatalf("Encountered error with opening/reading the file input: %s.\n", *filePtr)
		os.Exit(1)
	}
	s := string(b) // string value of input
	name := filepath.Base(*filePtr)
	fmt.Println(name, s)
}
