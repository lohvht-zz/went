package lang

import (
	"bytes"
	"fmt"
)

// WType is an interface where all other `went` language data structures
// should implemented, null is not within these types
type WType interface {
	IsZeroValue() WBool // returns true if the value is zero value
	Equals(WType) WBool // returns true if the object compared to it is equals
	String() string
}

// WNull is the null/none type in went, it is a value for no values
type WNull struct{}

// IsZeroValue always returns true for null type
func (w WNull) IsZeroValue() WBool { return true }

// Equals checks if the type compared to is equal
func (w WNull) Equals(w2 WType) WBool {
	if _, ok := w2.(WNull); ok {
		return true
	}
	return false
}
func (w WNull) String() string { return "null" }

// WNum is a number type in went, it combines both integers as well as floats
type WNum float64

// IsZeroValue returns the zero value of a went number value
func (w WNum) IsZeroValue() WBool { return w == 0 }

// Equals checks if the type compared to is equal
func (w WNum) Equals(w2 WType) WBool {
	if v, ok := w2.(WNum); ok {
		return w == v
	}
	return false
}

func (w WNum) String() string { return fmt.Sprintf("%v", float64(w)) }

// IsInt checks if WNum is an integer, if not its a float
func (w WNum) IsInt() bool { return float64(w) == float64(int64(w)) }

// WString is a string
type WString string

// IsZeroValue returns the zero value of a went string value
func (w WString) IsZeroValue() WBool { return w == "" }

// Equals checks if the type compared to is equal
func (w WString) Equals(w2 WType) WBool {
	if v, ok := w2.(WString); ok {
		return w == v
	}
	return false
}

func (w WString) String() string { return fmt.Sprintf("'%v'", string(w)) }

// WBool is a boolean
type WBool bool

// IsZeroValue returns the zero value of a went boolean value
func (w WBool) IsZeroValue() WBool { return !w }

// Equals checks if the type compared to is equal
func (w WBool) Equals(w2 WType) WBool {
	if v, ok := w2.(WBool); ok {
		return w == v
	}
	return false
}

func (w WBool) String() string { return fmt.Sprintf("%v", bool(w)) }

// WList is a list
type WList []WType

// IsZeroValue returns the zero value of a went list
func (w WList) IsZeroValue() WBool { return len(w) == 0 }

// Equals checks if the type compared to is equal
func (w WList) Equals(w2 WType) WBool {
	v, ok := w2.(WList)
	if !ok {
		return false
	} else if len(w) != len(v) {
		return false
	}
	for i := 0; i < len(w); i++ {
		if !w[i].Equals(v[i]) {
			return false
		}
	}
	return true
}

func (w WList) String() string { return fmt.Sprintf("%v", []WType(w)) }

var (
	tab        = "\t"
	twoSpaces  = "  "
	fourSpaces = "    "
)

// Wmap is a naive implementation of a went "map" data structure
// a data structure that maps strings to other values in wentlang
type Wmap map[string]WType

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

// Equals checks if the type compared to is equal
func (w Wmap) Equals(w2 WType) WBool {
	map2, ok := w2.(Wmap)
	if !ok {
		return false
	} else if len(w) != len(map2) {
		return false
	}
	for k1, v1 := range w {
		if !v1.Equals(map2[k1]) {
			return false
		}
	}
	return true
}

func (w Wmap) String() string { return w.toString(0) }
