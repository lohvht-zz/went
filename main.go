package main

import (
	"os"

	"github.com/lohvht/nondescript/cmd"
)

// REVIEW: Temporary dump for a preliminary implementation of the "map":
// Naive implementation of a went "map" data structure, a data structure that maps strings
// to other values in wentlang
// type Wmap map[string]interface{}
// tab = "\t"
// twoSpaces = "  "
// fourSpaces = "    "
// func (w Wmap) toString(tabLevel int) string {
// 	var buffer bytes.Buffer
// 	buffer.WriteString("{\n")
// 	for k, v := range w {
// 		for i := 0; i < tabLevel+1; i++ {
// 			buffer.WriteString("  ")
// 		}
// 		switch vTyped := v.(type) {
// 		case Wmap:
// 			buffer.WriteString(fmt.Sprintf("%s: %v,\n", k, vTyped.toString(tabLevel+1)))
// 		default:
// 			buffer.WriteString(fmt.Sprintf("%s: %v,\n", k, vTyped))
// 		}
//
// 	}
// 	for i := 0; i < tabLevel; i++ {
// 		buffer.WriteString("  ")
// 	}
// 	buffer.WriteString("}")
// 	return buffer.String()
// }
//
// func (w Wmap) String() string {
// 	return w.toString(0)
// }

func main() {
	os.Exit(cmd.Run())
}
