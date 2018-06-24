package utils

import (
	"fmt"
	"testing"
)

var tokenNames = map[tokenType]string{
	tokenError:        "error",
	tokenBool:         "boolean",
	tokenEquals:       "=",
	tokenDoubleEquals: "==",
	tokenNotEquals:    "!=",
	tokenLogicalNot:   "!",
	tokenEOF:          "EOF",
	tokenProperty:     "property",
	tokenIdentifier:   "identifier",
	tokenLeftParen:    "(",
	tokenRightParan:   ")",
	tokenLeftBrace:    "{",
	tokenRightBrace:   "}",
	tokenNumber:       "number",
	tokenOp:           "mathOp",
	tokenSpace:        "space",
	tokenNewline:      "newline",
	tokenQuotedString: "string",
	tokenRawString:    "raw string",
	tokenOr:           "or",
	tokenAnd:          "and",

	// Keywords
	tokenFunctionDef: "func",
	tokenVar:         "var",
	tokenIf:          "if",
	tokenElse:        "else",
	tokenElseIf:      "elseIf",
	tokenFor:         "for",
	tokenRange:       "range",
	tokenNull:        "null",
	// tokenReturn
	// tokenIn
	// tokenWhile
}

func (i tokenType) String() string {
	s := tokenNames[i]
	if s == "" {
		return fmt.Sprintf("item%d", int(i))
	}
	return s
}

// makeToken creates a token given a tokenType and a string denoting its value
func makeToken(typ tokenType, value string) token {
	return token{typ: typ, value: value}
}

var (
	tknEOF     = makeToken(tokenEOF, "")
	tknNL      = makeToken(tokenNewline, "\n")
	tknEqls    = makeToken(tokenEquals, "=")
	tknPlus    = makeToken(tokenOp, "+")
	tknMinus   = makeToken(tokenOp, "-")
	tknDiv     = makeToken(tokenOp, "/")
	tknMult    = makeToken(tokenOp, "*")
	tknMod     = makeToken(tokenOp, "%")
	tknSpace   = makeToken(tokenSpace, " ")
	tknFuncDef = makeToken(tokenFunctionDef, "func")
	tknVar     = makeToken(tokenVar, "var")
	tknIf      = makeToken(tokenIf, "if")
	tknElse    = makeToken(tokenElse, "else")
	tknElseIf  = makeToken(tokenElseIf, "elseIf")
	tknFor     = makeToken(tokenFor, "for")
	tknRange   = makeToken(tokenRange, "range")
	tknNull    = makeToken(tokenNull, "null")
)

type lexTestcase struct {
	name   string
	input  string
	tokens []token
}

var lexTests = []lexTestcase{
	// Positive Test Cases
	{"empty", "", []token{tknEOF}},
	{"line comment", "//Hi", []token{tknEOF}},
	{"line comment with \\n", "//Hello world\n", []token{
		tknNL,
		tknEOF,
	}},
	{"2 line comments with \\r\\n", "//Hello world\r\n//Howdy do", []token{
		tknNL,
		tknEOF,
	}},
	{"multiline comment", "/* This should be a comment\n more paragraphs*/", []token{
		tknEOF,
	}},
	{"division parse", "var x = 1.2 /* 2 *// 2", []token{
		tknVar,
		tknSpace,
		makeToken(tokenIdentifier, "x"),
		tknSpace,
		tknEqls,
		tknSpace,
		makeToken(tokenNumber, "1.2"),
		tknSpace,
		tknDiv,
		tknSpace,
		makeToken(tokenNumber, "2"),
		tknEOF,
	}},
	// Error Test Cases
}

func TestLex(t *testing.T) {
	for _, testcase := range lexTests {
		outputTokens := collect(&testcase)
		if !equal(outputTokens, testcase.tokens, false) {
			t.Errorf("%s: got\n\t%+v\nexpected\n\t%v", testcase.name, outputTokens, testcase.tokens)
		}
	}
}

// Helper Methods to check equality for tests and collect tokens

// collect gathers the emitted items into a token slice
func collect(tc *lexTestcase) (tkns []token) {
	l := tokenise(tc.name, tc.input)
	for {
		tkn := l.nextToken()
		tkns = append(tkns, tkn)
		if tkn.typ == tokenEOF || tkn.typ == tokenError {
			break
		}
	}
	return
}

func equal(tknLst1, tknLst2 []token, checkPos bool) bool {
	if len(tknLst1) != len(tknLst2) {
		return false
	}
	for k := range tknLst1 {
		if tknLst1[k].typ != tknLst2[k].typ {
			return false
		}
		if tknLst1[k].value != tknLst2[k].value {
			return false
		}
		if checkPos && tknLst1[k].pos != tknLst2[k].pos {
			return false
		}
		if checkPos && tknLst1[k].line != tknLst2[k].line {
			return false
		}
	}
	return true
}
