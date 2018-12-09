package token

import (
	"testing"
)

// makeToken creates a Token given a Type and a string denoting its value
func makeToken(typ Type, value string) Token { return Token{Type: typ, Value: value} }

// makeName is a helper method that creates an identifier with the string value
func makeName(value string) Token { return makeToken(NAME, value) }

// makeError is a helper method that creates an error with the string value
func makeError(value string) Token { return makeToken(ERROR, value) }

var (
	tknEOF   = makeToken(EOF, "")
	tknDot   = makeToken(DOT, ".")
	tknLR    = makeToken(LROUND, tokenTypes[LROUND])
	tknRR    = makeToken(RROUND, tokenTypes[RROUND])
	tknLC    = makeToken(LCURLY, tokenTypes[LCURLY])
	tknRC    = makeToken(RCURLY, tokenTypes[RCURLY])
	tknLS    = makeToken(LSQUARE, tokenTypes[LSQUARE])
	tknRS    = makeToken(RSQUARE, tokenTypes[RSQUARE])
	tknColon = makeToken(COLON, tokenTypes[COLON])
	tknSemi  = makeToken(SEMICOLON, tokenTypes[SEMICOLON])
	tknComma = makeToken(COMMA, tokenTypes[COMMA])
	// Operators
	// Arithmetic Operators
	tknPlus = makeToken(PLUS, tokenTypes[PLUS])
	tknMin  = makeToken(MINUS, tokenTypes[MINUS])
	tknDiv  = makeToken(DIV, tokenTypes[DIV])
	tknMult = makeToken(MULT, tokenTypes[MULT])
	tknMod  = makeToken(MOD, tokenTypes[MOD])
	// Assignment Operators
	tknAss     = makeToken(ASSIGN, tokenTypes[ASSIGN])
	tknPlusAss = makeToken(PLUSASSIGN, tokenTypes[PLUSASSIGN])
	tknMinAss  = makeToken(MINUSASSIGN, tokenTypes[MINUSASSIGN])
	tknDivAss  = makeToken(DIVASSIGN, tokenTypes[DIVASSIGN])
	tknMultAss = makeToken(MULTASSIGN, tokenTypes[MULTASSIGN])
	tknModAss  = makeToken(MODASSIGN, tokenTypes[MODASSIGN])
	// Comparison Operators
	tknEql  = makeToken(EQ, tokenTypes[EQ])
	tknNEql = makeToken(NEQ, tokenTypes[NEQ])
	tknGr   = makeToken(GR, tokenTypes[GR])
	tknSm   = makeToken(SM, tokenTypes[SM])
	tknGrEq = makeToken(GREQ, tokenTypes[GREQ])
	tknSmEq = makeToken(SMEQ, tokenTypes[SMEQ])
	// Logical Operators
	tknLogicN = makeToken(LOGICALNOT, tokenTypes[LOGICALNOT])
	tknOr     = makeToken(LOGICALOR, tokenTypes[LOGICALOR])
	tknAnd    = makeToken(LOGICALAND, tokenTypes[LOGICALAND])

	// keywords
	tknFuncDef = makeToken(FUNC, tokenTypes[FUNC])
	tknIf      = makeToken(IF, tokenTypes[IF])
	tknElse    = makeToken(ELSE, tokenTypes[ELSE])
	tknElseIf  = makeToken(ELIF, tokenTypes[ELIF])
	tknFor     = makeToken(FOR, tokenTypes[FOR])
	tknNull    = makeToken(NULL, tokenTypes[NULL])
	tknF       = makeToken(FALSE, tokenTypes[FALSE])
	tknT       = makeToken(TRUE, tokenTypes[TRUE])
	tknWhile   = makeToken(WHILE, tokenTypes[WHILE])
	tknReturn  = makeToken(RETURN, tokenTypes[RETURN])
	tknIn      = makeToken(IN, tokenTypes[IN])
	tknBreak   = makeToken(BREAK, tokenTypes[BREAK])
	tknCont    = makeToken(CONT, tokenTypes[CONT])
	tknVar     = makeToken(VAR, tokenTypes[VAR])
)

type lexTestcase struct {
	name   string
	input  string
	tokens []Token
}

var lexTests = []lexTestcase{
	// Positive Test Cases
	{"empty",
		"",
		[]Token{tknEOF},
	},
	{"line comment",
		"//Hi",
		[]Token{tknEOF},
	},
	{"line comment with \\n",
		"//Hello world\n",
		[]Token{tknEOF},
	},
	{"2 line comments with \\r\\n",
		"//Hello world\r\n//Howdy do",
		[]Token{tknEOF},
	},
	{"multiline comment",
		`/* This should be a comment
		more paragraphs
		and it should be parsed correctly
		*/
		x = 3.123
		`,
		[]Token{makeName("x"), tknAss, makeToken(FLOAT, "3.123"), tknSemi, tknEOF},
	},
	{"division parse",
		`x = 1.2 /* 2 *// 2
		`,
		[]Token{makeName("x"), tknAss, makeToken(FLOAT, "1.2"),
			tknDiv, makeToken(INT, "2"), tknSemi, tknEOF,
		},
	},
	{"keywords",
		"func if else elif for null false true while return break continue in var",
		[]Token{tknFuncDef, tknIf, tknElse, tknElseIf, tknFor, tknNull, tknF, tknT,
			tknWhile, tknReturn, tknBreak, tknCont, tknIn, tknVar, tknEOF,
		},
	},
	{"arithmetic operators",
		"+ - / * %",
		[]Token{tknPlus, tknMin, tknDiv, tknMult, tknMod, tknEOF},
	},
	{"assignment operators",
		"= += -= /= *= %=",
		[]Token{tknAss, tknPlusAss, tknMinAss, tknDivAss, tknMultAss, tknModAss, tknEOF},
	},
	{"comparison and logical operators",
		"== != > < >= <= ! || &&",
		[]Token{tknEql, tknNEql, tknGr, tknSm, tknGrEq, tknSmEq, tknLogicN,
			tknOr, tknAnd, tknEOF,
		},
	},
	{"identifiers and dots",
		"x.y.z+n.q.w()",
		[]Token{makeName("x"), tknDot, makeName("y"), tknDot, makeName("z"), tknPlus,
			makeName("n"), tknDot, makeName("q"), tknDot, makeName("w"),
			tknLR, tknRR, tknEOF,
		},
	},
	// Error Test Cases
	{"single | error",
		"x | y",
		[]Token{makeName("x"), makeError(`expected Token U+007C '|'`)},
	},
	{"single & error",
		"x & y",
		[]Token{makeName("x"), makeError(`expected Token U+0026 '&'`)},
	},
	{"typo right bracket )",
		"x + ) y",
		[]Token{makeName("x"), tknPlus, makeError(`unexpected right bracket U+0029 ')'`)},
	},
	{"extra right bracket )",
		"(x + 1)) * y",
		[]Token{tknLR, makeName("x"), tknPlus, makeToken(INT, "1"),
			tknRR, makeError(`unexpected right bracket U+0029 ')'`),
		},
	},
	{"extra right brace bracket }",
		"if x == 012.999 { return y }}",
		[]Token{tknIf, makeName("x"), tknEql, makeToken(FLOAT, "012.999"),
			tknLC, tknReturn, makeName("y"), tknSemi, tknRC,
			makeError(`unexpected right bracket U+007D '}'`),
		},
	},
	{"extra right square bracket ]",
		"[x, 2, w]]",
		[]Token{tknLS, makeName("x"), tknComma, makeToken(INT, "2"),
			tknComma, makeName("w"), tknRS, makeError(`unexpected right bracket U+005D ']'`),
		},
	},
	{"unclosed left bracket",
		"(x+y)*((1/1.324)%4",
		[]Token{tknLR, makeName("x"), tknPlus, makeName("y"), tknRR, tknMult,
			tknLR, tknLR, makeToken(INT, "1"), tknDiv,
			makeToken(FLOAT, "1.324"), tknRR, tknMod, makeToken(INT, "4"),
			makeError(`unclosed left bracket: U+0028 '('`),
		},
	},
	{"unclosed multiline comment",
		`/* This is an unclosed comment!
		/`,
		[]Token{makeError("Multiline comment is not closed")},
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

// collect gathers the emitted items into a Token slice
func collect(tc *lexTestcase) (tkns []Token) {
	l := Tokenise(tc.name, tc.input)
	for {
		tkn := l.Next()
		tkns = append(tkns, tkn)
		if tkn.Type == EOF || tkn.Type == ERROR {
			break
		}
	}
	return
}

func equal(tknLst1, tknLst2 []Token, checkPos bool) bool {
	if len(tknLst1) != len(tknLst2) {
		return false
	}
	for k := range tknLst1 {
		tkn1 := tknLst1[k]
		tkn2 := tknLst2[k]
		switch {
		case tkn1.Type != tkn2.Type,
			tkn1.Value != tkn2.Value && !(tkn1.Type == SEMICOLON && tkn2.Type == SEMICOLON),
			checkPos && tkn1.Pos == tkn2.Pos:
			return false
		}
	}
	return true
}
