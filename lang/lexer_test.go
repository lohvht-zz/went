package lang

import (
	"testing"
)

// makeToken creates a token given a tokenType and a string denoting its value
func makeToken(typ tokenType, value string) token {
	return token{typ: typ, value: value}
}

// makeName is a helper method that creates an identifier with the string value
func makeName(value string) token {
	return makeToken(tokenName, value)
}

// makeError is a helper method that creates an error with the string value
func makeError(value string) token {
	return makeToken(tokenError, value)
}

var (
	tknEOF   = makeToken(tokenEOF, "")
	tknDot   = makeToken(tokenDot, ".")
	tknLR    = makeToken(tokenLeftRound, "(")
	tknRR    = makeToken(tokenRightRound, ")")
	tknLC    = makeToken(tokenLeftCurly, "{")
	tknRC    = makeToken(tokenRightCurly, "}")
	tknLS    = makeToken(tokenLeftSquare, "[")
	tknRS    = makeToken(tokenRightSquare, "]")
	tknColon = makeToken(tokenColon, ":")
	tknSemi  = makeToken(tokenSemicolon, ";")
	tknComma = makeToken(tokenComma, ",")
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
	tknF       = makeToken(tokenFalse, "false")
	tknT       = makeToken(tokenTrue, "true")
	tknWhile   = makeToken(tokenWhile, "while")
	tknReturn  = makeToken(tokenReturn, "return")
	tknIn      = makeToken(tokenIn, "in")
	tknBreak   = makeToken(tokenBreak, "break")
	tknCont    = makeToken(tokenCont, "continue")
	tknVar     = makeToken(tokenVar, "var")
)

type lexTestcase struct {
	name   string
	input  string
	tokens []token
}

var lexTests = []lexTestcase{
	// Positive Test Cases
	{"empty",
		"",
		[]token{tknEOF},
	},
	{"line comment",
		"//Hi",
		[]token{tknEOF},
	},
	{"line comment with \\n",
		"//Hello world\n",
		[]token{tknEOF},
	},
	{"2 line comments with \\r\\n",
		"//Hello world\r\n//Howdy do",
		[]token{tknEOF},
	},
	{"multiline comment",
		`/* This should be a comment
		more paragraphs
		and it should be parsed correctly
		*/`,
		[]token{tknEOF},
	},
	{"division parse",
		"x = 1.2 /* 2 *// 2",
		[]token{makeName("x"), tknAss, makeToken(tokenNumber, "1.2"),
			tknDiv, makeToken(tokenNumber, "2"), tknEOF,
		},
	},
	{"keywords",
		"func if else elif for null false true while return break continue in var",
		[]token{tknFuncDef, tknIf, tknElse, tknElseIf, tknFor, tknNull, tknF, tknT,
			tknWhile, tknReturn, tknBreak, tknCont, tknIn, tknVar, tknEOF,
		},
	},
	{"arithmetic operators",
		"+ - / * %",
		[]token{tknPlus, tknMin, tknDiv, tknMult, tknMod, tknEOF},
	},
	{"assignment operators",
		"= += -= /= *= %=",
		[]token{tknAss, tknPlusAss, tknMinAss, tknDivAss, tknMultAss, tknModAss, tknEOF},
	},
	{"comparison and logical operators",
		"== != > < >= <= ! || &&",
		[]token{tknEql, tknNEql, tknGr, tknSm, tknGrEq, tknSmEq, tknLogicN,
			tknOr, tknAnd, tknEOF,
		},
	},
	{"identifiers and dots",
		"x.y.z+n.q.w()",
		[]token{makeName("x"), tknDot, makeName("y"), tknDot, makeName("z"), tknPlus,
			makeName("n"), tknDot, makeName("q"), tknDot, makeName("w"),
			tknLR, tknRR, tknEOF,
		},
	},
	// Error Test Cases
	{"single | error",
		"x | y",
		[]token{makeName("x"), makeError(`expected token U+007C '|'`)},
	},
	{"single & error",
		"x & y",
		[]token{makeName("x"), makeError(`expected token U+0026 '&'`)},
	},
	{"typo right bracket )",
		"x + ) y",
		[]token{makeName("x"), tknPlus, makeError(`unexpected right bracket U+0029 ')'`)},
	},
	{"extra right bracket )",
		"(x + 1)) * y",
		[]token{tknLR, makeName("x"), tknPlus, makeToken(tokenNumber, "1"),
			tknRR, makeError(`unexpected right bracket U+0029 ')'`),
		},
	},
	{"extra right brace bracket }",
		"if x == 1 { return y }}",
		[]token{tknIf, makeName("x"), tknEql, makeToken(tokenNumber, "1"),
			tknLC, tknReturn, makeName("y"), tknSemi, tknRC,
			makeError(`unexpected right bracket U+007D '}'`),
		},
	},
	{"extra right square bracket ]",
		"[x, 2, w]]",
		[]token{tknLS, makeName("x"), tknComma, makeToken(tokenNumber, "2"),
			tknComma, makeName("w"), tknRS, makeError(`unexpected right bracket U+005D ']'`),
		},
	},
	{"unclosed left bracket",
		"(x+y)*((1/1.324)%4",
		[]token{tknLR, makeName("x"), tknPlus, makeName("y"), tknRR, tknMult,
			tknLR, tknLR, makeToken(tokenNumber, "1"), tknDiv,
			makeToken(tokenNumber, "1.324"), tknRR, tknMod, makeToken(tokenNumber, "4"),
			makeError(`unclosed left bracket: U+0028 '('`),
		},
	},
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
			// Check due to Automatic semicolon insertion, some semicolon tokens may
			// contain values in the strings that correspond
			if tknLst1[k].typ != tokenSemicolon || tknLst2[k].typ != tokenSemicolon {
				return false
			}
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
