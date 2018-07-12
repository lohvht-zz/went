package lang

import (
	"testing"
)

// makeToken creates a token given a tokenType and a string denoting its value
func makeToken(typ tokenType, value string) token {
	return token{typ: typ, value: value}
}

var (
	tknEOF  = makeToken(tokenEOF, "")
	tknSemi = makeToken(tokenSemicolon, ";")
	// Operators
	// Arithmetic Operators
	tknPlus = makeToken(tokenPlus, "+")
	tknMin  = makeToken(tokenMinus, "-")
	tknDiv  = makeToken(tokenDiv, "/")
	tknMult = makeToken(tokenMult, "*")
	tknMod  = makeToken(tokenMod, "%")
	// Assignment Operators
	tknAss     = makeToken(tokenAssign, "=")
	tknPlusAss = makeToken(tokenPlusAssign, "+=")
	tknMinAss  = makeToken(tokenMinusAssign, "-=")
	tknDivAss  = makeToken(tokenDivAssign, "/=")
	tknMultAss = makeToken(tokenMultAssign, "*=")
	tknModAss  = makeToken(tokenModAssign, "%=")
	// Comparison Operators
	tknEql  = makeToken(tokenEquals, "==")
	tknNEql = makeToken(tokenNotEquals, "!=")
	tknGr   = makeToken(tokenGreater, ">")
	tknSm   = makeToken(tokenSmaller, "<")
	tknGrEq = makeToken(tokenGreaterEquals, ">=")
	tknSmEq = makeToken(tokenSmallerEquals, "<=")
	// Logical Operators
	tknLogicN = makeToken(tokenLogicalNot, "!")
	tknOr     = makeToken(tokenOr, "||")
	tknAnd    = makeToken(tokenAnd, "&&")

	// keywords
	tknFuncDef = makeToken(tokenFunc, "func")
	tknIf      = makeToken(tokenIf, "if")
	tknElse    = makeToken(tokenElse, "else")
	tknElseIf  = makeToken(tokenElseIf, "elif")
	tknFor     = makeToken(tokenFor, "for")
	tknNull    = makeToken(tokenNull, "null")
	tknWhile   = makeToken(tokenWhile, "while")
	tknReturn  = makeToken(tokenReturn, "return")
	tknIn      = makeToken(tokenIn, "in")
	tknBreak   = makeToken(tokenBreak, "break")
	tknCont    = makeToken(tokenCont, "continue")
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
		tknEOF,
	}},
	{"2 line comments with \\r\\n", "//Hello world\r\n//Howdy do", []token{
		tknEOF,
	}},
	{"multiline comment", "/* This should be a comment\n more paragraphs*/", []token{
		tknEOF,
	}},
	{"division parse", "x = 1.2 /* 2 *// 2", []token{
		makeToken(tokenIdentifier, "x"),
		tknAss,
		makeToken(tokenNumber, "1.2"),
		tknDiv,
		makeToken(tokenNumber, "2"),
		tknEOF,
	}},
	{"keywords", "func if else elif for null while return break continue in", []token{
		tknFuncDef,
		tknIf,
		tknElse,
		tknElseIf,
		tknFor,
		tknNull,
		tknWhile,
		tknReturn,
		tknBreak,
		tknCont,
		tknIn,
		tknEOF,
	}},
	{"arithmetic operators", "+ - / * %", []token{
		tknPlus,
		tknMin,
		tknDiv,
		tknMult,
		tknMod,
		tknEOF,
	}},
	{"assignment operators", "= += -= /= *= %=", []token{
		tknAss,
		tknPlusAss,
		tknMinAss,
		tknDivAss,
		tknMultAss,
		tknModAss,
		tknEOF,
	}},
	{"comparison and logical operators", "== != > < >= <= ! || &&", []token{
		tknEql,
		tknNEql,
		tknGr,
		tknSm,
		tknGrEq,
		tknSmEq,
		tknLogicN,
		tknOr,
		tknAnd,
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
