package lang

/**
 * NOTE: This file is deprecated and should be cleaned up (i.e. Deleted) as soon
 * as the main language is written
 * It's being kept around as a reference while the lox implementation is being
 * written
 */

import (
	"bytes"
	"fmt"
)

// WType is an interface where all other `went` language data structures
// should implemented, null is not within these types
type WType interface {
	IsZeroValue() WBool                    // returns true if the value is zero value
	Equals(w2 WType) WBool                 // returns true if the object compared to it is equals
	Sm(w2 WType, orEq bool) (WBool, error) // returns true if smaller than other type, returns error if type is not supported
	Gr(w2 WType, orEq bool) (WBool, error) // returns true if greather than other type, returns error if type is not supported
	String() string
}

func opError(w1, w2 WType, compString string) error {
	return fmt.Errorf("'%s' not supported between types '%T' and '%T'", compString, w1, w2)
}

var (
	sm   = "<"
	smE  = "<="
	gr   = ">"
	grE  = ">="
	eql  = "=="
	nEql = "!="
)

// WNull is the null/none type in went, it is a value for no values
type WNull struct{}

// IsZeroValue always returns true for null type
func (w WNull) IsZeroValue() WBool { return true }

// Equals checks if the type compared to is equal
func (w WNull) Equals(w2 WType) WBool {
	_, ok := w2.(WNull)
	return WBool(ok)
}

// Sm will always return an error for null as such a relation is not supported yet
func (w WNull) Sm(w2 WType, orEq bool) (WBool, error) {
	switch v := w2.(type) {
	default:
		var operator string
		if orEq {
			operator = smE
		} else {
			operator = sm
		}
		err := opError(w, v, operator)
		return false, err
	}
}

// Gr (see Sm)
// a >= b <==> !(a < b)
// a > b <==> !(a <= b)
func (w WNull) Gr(w2 WType, orEq bool) (WBool, error) {
	smRes, err := w.Sm(w2, !orEq)
	if err != nil {
		var operator string
		if orEq {
			operator = grE
		} else {
			operator = gr
		}
		return false, opError(w, w2, operator)
	}
	// Should be impossible to reach here as WNull will always have err != nil
	return !smRes, nil
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

// Sm returns true if w is smaller than w2, false else, returns an error if the
// 2 are of different types
func (w WNum) Sm(w2 WType, orEq bool) (WBool, error) {
	switch v := w2.(type) {
	case WNum:
		if orEq {
			return WBool(w <= v), nil
		}
		return WBool(w < v), nil
	default:
		var operator string
		if orEq {
			operator = smE
		} else {
			operator = sm
		}
		err := opError(w, v, operator)
		return false, err
	}
}

// Gr (see Sm)
// a >= b <==> !(a < b)
// a > b <==> !(a <= b)
func (w WNum) Gr(w2 WType, orEq bool) (WBool, error) {
	smRes, err := w.Sm(w2, !orEq)
	if err != nil {
		var operator string
		if orEq {
			operator = grE
		} else {
			operator = gr
		}
		return false, opError(w, w2, operator)
	}
	return !smRes, nil
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

// Sm returns true if w is smaller than w2, false else, returns an error if the
// 2 are of different types
func (w WString) Sm(w2 WType, orEq bool) (WBool, error) {
	switch v := w2.(type) {
	case WString:
		if orEq {
			return WBool(w <= v), nil
		}
		return WBool(w < v), nil
	default:
		var operator string
		if orEq {
			operator = smE
		} else {
			operator = sm
		}
		err := opError(w, v, operator)
		return false, err
	}
}

// Gr (see Sm)
// a >= b <==> !(a < b)
// a > b <==> !(a <= b)
func (w WString) Gr(w2 WType, orEq bool) (WBool, error) {
	smRes, err := w.Sm(w2, !orEq)
	if err != nil {
		var operator string
		if orEq {
			operator = grE
		} else {
			operator = gr
		}
		return false, opError(w, w2, operator)
	}
	return !smRes, nil
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

// Sm will always return false and an error for WBool as WBool has
// no order relation
func (w WBool) Sm(w2 WType, orEq bool) (WBool, error) {
	switch v := w2.(type) {
	default:
		var operator string
		if orEq {
			operator = smE
		} else {
			operator = sm
		}
		err := opError(w, v, operator)
		return false, err
	}
}

// Gr (see Sm)
// a >= b <==> !(a < b)
// a > b <==> !(a <= b)
func (w WBool) Gr(w2 WType, orEq bool) (WBool, error) {
	smRes, err := w.Sm(w2, !orEq)
	if err != nil {
		var operator string
		if orEq {
			operator = grE
		} else {
			operator = gr
		}
		return false, opError(w, w2, operator)
	}
	return !smRes, nil
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

// Sm returns true if w is smaller than w2, false else, returns an error if the
// 2 are of different types
func (w WList) Sm(w2 WType, orEq bool) (WBool, error) {
	switch v := w2.(type) {
	case WList:
		minLen := min(len(w), len(v))
		for i := 0; i < minLen; i++ {
			if !w[i].Equals(v[i]) {
				return w[i].Sm(v[i], orEq)
			}
		}
		if orEq {
			return len(w) <= len(v), nil
		}
		return len(w) < len(v), nil
	default:
		var operator string
		if orEq {
			operator = smE
		} else {
			operator = sm
		}
		err := opError(w, v, operator)
		return false, err
	}
}

// Gr (see Sm)
// a >= b <==> !(a < b)
// a > b <==> !(a <= b)
func (w WList) Gr(w2 WType, orEq bool) (WBool, error) {
	smRes, err := w.Sm(w2, !orEq)
	if err != nil {
		var operator string
		if orEq {
			operator = grE
		} else {
			operator = gr
		}
		return false, opError(w, w2, operator)
	}
	return !smRes, nil
}

func (w WList) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("[")
	for i, v := range w {
		buffer.WriteString(v.String())
		if i != len(w)-1 {
			buffer.WriteString(", ")
		}
	}
	buffer.WriteString("]")
	return buffer.String()
}

var (
	tab        = "\t"
	twoSpaces  = "  "
	fourSpaces = "    "
)

// Wmap is a naive implementation of a went "map" data structure
// a data structure that maps strings to other values in wentlang
type Wmap map[string]WType

// toString returns a string that is essentially a pretty-printed formatted Wmap
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

// Sm will always return false and an error for Wmap as Wmap has
// no order relation
func (w Wmap) Sm(w2 WType, orEq bool) (WBool, error) {
	switch v := w2.(type) {
	default:
		var operator string
		if orEq {
			operator = smE
		} else {
			operator = sm
		}
		err := opError(w, v, operator)
		return false, err
	}
}

// Gr (see Sm)
// a >= b <==> !(a < b)
// a > b <==> !(a <= b)
func (w Wmap) Gr(w2 WType, orEq bool) (WBool, error) {
	smRes, err := w.Sm(w2, !orEq)
	if err != nil {
		var operator string
		if orEq {
			operator = grE
		} else {
			operator = gr
		}
		return false, opError(w, w2, operator)
	}
	return !smRes, nil
}

func (w Wmap) String() string { return w.toString(0) }

// Helper functions

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
