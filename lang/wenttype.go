package lang

import (
	"bytes"
	"fmt"
)

// WType is an interface where all other `went` language data structures
// should implemented, null is not within these types
type WType interface {
	IsZeroValue() WBool
	String() string
}

// WNull is the null/none type in went, it is a value for no values
type WNull struct{}

// IsZeroValue always returns true for null type
func (w WNull) IsZeroValue() WBool { return true }
func (w WNull) String() string     { return "null" }

// WNum is a number type in went, it combines both integers as well as floats
type WNum float64

// IsZeroValue returns the zero value of a went number value
func (w WNum) IsZeroValue() WBool { return w == 0 }
func (w WNum) String() string     { return fmt.Sprintf("%v", float64(w)) }

// IsInt checks if WNum is an integer, if not its a float
func (w WNum) IsInt() bool { return float64(w) == float64(int64(w)) }

// WString is a string
type WString string

// IsZeroValue returns the zero value of a went string value
func (w WString) IsZeroValue() WBool { return w == "" }
func (w WString) String() string     { return fmt.Sprintf("'%v'", string(w)) }

// WBool is a boolean
type WBool bool

// IsZeroValue returns the zero value of a went boolean value
func (w WBool) IsZeroValue() WBool { return !w }
func (w WBool) String() string     { return fmt.Sprintf("%v", bool(w)) }

// WList is a list
type WList []interface{}

// IsZeroValue returns the zero value of a went list
func (w WList) IsZeroValue() WBool { return len(w) == 0 }
func (w WList) String() string     { return fmt.Sprintf("%v", []interface{}(w)) }

var (
	tab        = "\t"
	twoSpaces  = "  "
	fourSpaces = "    "
)

// Wmap is a naive implementation of a went "map" data structure
// a data structure that maps strings to other values in wentlang
type Wmap map[string]interface{}

func (w Wmap) toString(tabLevel int) string {
	var buffer bytes.Buffer
	buffer.WriteString("{\n")
	for k, v := range w {
		// adds a new tab in addition to the number of tabLevels while inside the body
		for i := 0; i < tabLevel+1; i++ {
			buffer.WriteString(twoSpaces)
		}
		switch vTyped := v.(type) {
		case Wmap:
			buffer.WriteString(fmt.Sprintf("%s: %v,\n", k, vTyped.toString(tabLevel+1)))
		default:
			buffer.WriteString(fmt.Sprintf("%s: %v,\n", k, vTyped))
		}
	}
	for i := 0; i < tabLevel; i++ {
		buffer.WriteString(twoSpaces)
	}
	buffer.WriteString("}")
	return buffer.String()
}

// IsZeroValue returns the zero value of a went map
func (w Wmap) IsZeroValue() WBool { return len(w) == 0 }
func (w Wmap) String() string     { return w.toString(0) }
