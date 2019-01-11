package lexer

import (
	"testing"

	"github.com/lohvht/went/lang/token"
)

type testHandler struct{ errors token.ErrorList }

func initTestHandler(name, input string) (*Lexer, *testHandler) {
	th := &testHandler{}
	eh := func(filename string, pos token.Pos, msg string) { th.errors.Add(filename, pos, msg) }
	l := New(name, input, eh)
	return l, th
}

// makeToken creates a Token given a Type and a string denoting its value
func makeToken(typ token.Type, value string) token.Token {
	return token.Token{Type: typ, Value: value}
}

// makeName is a helper method that creates an identifier with the string value
func makeName(value string) token.Token { return makeToken(token.NAME, value) }

func getKeywordTypToStr() map[token.Type]string {
	m := make(map[token.Type]string, len(token.Keywords))
	for k, v := range token.Keywords {
		m[v] = k
	}
	return m
}

var (
	tknEOF   = makeToken(token.EOF, "")
	tknDot   = makeToken(token.DOT, ".")
	tknLR    = makeToken(token.LROUND, "(")
	tknRR    = makeToken(token.RROUND, ")")
	tknLC    = makeToken(token.LCURLY, "{")
	tknRC    = makeToken(token.RCURLY, "}")
	tknLS    = makeToken(token.LSQUARE, "[")
	tknRS    = makeToken(token.RSQUARE, "]")
	tknColon = makeToken(token.COLON, ":")
	tknSemi  = makeToken(token.SEMICOLON, ";")
	tknComma = makeToken(token.COMMA, ",")
	// Operators
	// Arithmetic Operators
	tknPlus = makeToken(token.PLUS, "+")
	tknMin  = makeToken(token.MINUS, "-")
	tknDiv  = makeToken(token.DIV, "/")
	tknMult = makeToken(token.MULT, "*")
	tknMod  = makeToken(token.MOD, "%")
	// Assignment Operators
	tknAss     = makeToken(token.ASSIGN, "=")
	tknPlusAss = makeToken(token.PLUSASSIGN, "+=")
	tknMinAss  = makeToken(token.MINUSASSIGN, "-=")
	tknDivAss  = makeToken(token.DIVASSIGN, "/=")
	tknMultAss = makeToken(token.MULTASSIGN, "*=")
	tknModAss  = makeToken(token.MODASSIGN, "%=")
	// Comparison Operators
	tknEql  = makeToken(token.EQ, "==")
	tknNEql = makeToken(token.NEQ, "!=")
	tknGr   = makeToken(token.GR, ">")
	tknSm   = makeToken(token.SM, "<")
	tknGrEq = makeToken(token.GREQ, ">=")
	tknSmEq = makeToken(token.SMEQ, "<=")
	// Logical Operators
	tknLogicN = makeToken(token.LOGICALNOT, "!")
	tknOr     = makeToken(token.LOGICALOR, "||")
	tknAnd    = makeToken(token.LOGICALAND, "&&")

	// keywords
	keywordTypToStr = getKeywordTypToStr()
	tknClass        = makeToken(token.CLASS, keywordTypToStr[token.CLASS])
	tknSuper        = makeToken(token.SUPER, keywordTypToStr[token.SUPER])
	tknSelf         = makeToken(token.SELF, keywordTypToStr[token.SELF])
	tknFuncDef      = makeToken(token.FUNC, keywordTypToStr[token.FUNC])
	tknIf           = makeToken(token.IF, keywordTypToStr[token.IF])
	tknElse         = makeToken(token.ELSE, keywordTypToStr[token.ELSE])
	tknElseIf       = makeToken(token.ELIF, keywordTypToStr[token.ELIF])
	tknFor          = makeToken(token.FOR, keywordTypToStr[token.FOR])
	tknNull         = makeToken(token.NULL, keywordTypToStr[token.NULL])
	tknF            = makeToken(token.FALSE, keywordTypToStr[token.FALSE])
	tknT            = makeToken(token.TRUE, keywordTypToStr[token.TRUE])
	tknWhile        = makeToken(token.WHILE, keywordTypToStr[token.WHILE])
	tknReturn       = makeToken(token.RETURN, keywordTypToStr[token.RETURN])
	tknIn           = makeToken(token.IN, keywordTypToStr[token.IN])
	tknBreak        = makeToken(token.BREAK, keywordTypToStr[token.BREAK])
	tknCont         = makeToken(token.CONT, keywordTypToStr[token.CONT])
	tknVar          = makeToken(token.VAR, keywordTypToStr[token.VAR])
)

type lexTestcase struct {
	name   string
	input  string
	tokens []token.Token
}

var lexTests = []lexTestcase{
	// Positive Test Cases
	{"empty",
		"",
		[]token.Token{tknEOF},
	},
	{"line comment",
		"//Hi",
		[]token.Token{tknEOF},
	},
	{"line comment with \\n",
		"//Hello world\n",
		[]token.Token{tknEOF},
	},
	{"2 line comments with \\r\\n",
		"//Hello world\r\n//Howdy do",
		[]token.Token{tknEOF},
	},
	{"multiline comment",
		`/* This should be a comment
		more paragraphs
		and it should be parsed correctly
		*/
		x = 3.123
		`,
		[]token.Token{makeName("x"), tknAss, makeToken(token.FLOAT, "3.123"), tknSemi, tknEOF},
	},
	{"division parse",
		`x = 1.2 /* 2 *// 2
		`,
		[]token.Token{makeName("x"), tknAss, makeToken(token.FLOAT, "1.2"),
			tknDiv, makeToken(token.INT, "2"), tknSemi, tknEOF,
		},
	},
	{"keywords",
		"func if else elif for null false true while return break continue in var class super self",
		[]token.Token{tknFuncDef, tknIf, tknElse, tknElseIf, tknFor, tknNull, tknF, tknT,
			tknWhile, tknReturn, tknBreak, tknCont, tknIn, tknVar, tknClass, tknSuper,
			tknSelf, tknEOF,
		},
	},
	{"arithmetic operators",
		"+ - / * %",
		[]token.Token{tknPlus, tknMin, tknDiv, tknMult, tknMod, tknEOF},
	},
	{"assignment operators",
		"= += -= /= *= %=",
		[]token.Token{tknAss, tknPlusAss, tknMinAss, tknDivAss, tknMultAss, tknModAss, tknEOF},
	},
	{"comparison and logical operators",
		"== != > < >= <= ! || &&",
		[]token.Token{tknEql, tknNEql, tknGr, tknSm, tknGrEq, tknSmEq, tknLogicN,
			tknOr, tknAnd, tknEOF,
		},
	},
	{"identifiers and dots",
		"x.y.z+n.q.w()",
		[]token.Token{makeName("x"), tknDot, makeName("y"), tknDot, makeName("z"), tknPlus,
			makeName("n"), tknDot, makeName("q"), tknDot, makeName("w"),
			tknLR, tknRR, tknEOF,
		},
	},
	{"numbers (int, float, hexadecimal, octals)",
		"123 .345 1.234 0x1237A 012374",
		[]token.Token{makeToken(token.INT, "123"), makeToken(token.FLOAT, ".345"),
			makeToken(token.FLOAT, "1.234"), makeToken(token.INT, "0x1237A"),
			makeToken(token.INT, "012374"), tknEOF,
		},
	},
	// TODO: Make error test cases work
	// // Error Test Cases
	// {"single | error",
	// 	"x | y",
	// 	[]token.Token{makeName("x"), makeError(`expected Token U+007C '|'`)},
	// },
	// {"single & error",
	// 	"x & y",
	// 	[]token.Token{makeName("x"), makeError(`expected Token U+0026 '&'`)},
	// },
	// {"typo right bracket )",
	// 	"x + ) y",
	// 	[]token.Token{makeName("x"), tknPlus, makeError(`unexpected right bracket U+0029 ')'`)},
	// },
	// {"extra right bracket )",
	// 	"(x + 1)) * y",
	// 	[]token.Token{tknLR, makeName("x"), tknPlus, makeToken(token.INT, "1"),
	// 		tknRR, makeError(`unexpected right bracket U+0029 ')'`),
	// 	},
	// },
	// {"extra right brace bracket }",
	// 	"if x == 012.999 { return y }}",
	// 	[]token.Token{tknIf, makeName("x"), tknEql, makeToken(token.FLOAT, "012.999"),
	// 		tknLC, tknReturn, makeName("y"), tknSemi, tknRC,
	// 		makeError(`unexpected right bracket U+007D '}'`),
	// 	},
	// },
	// {"extra right square bracket ]",
	// 	"[x, 2, w]]",
	// 	[]token.Token{tknLS, makeName("x"), tknComma, makeToken(token.INT, "2"),
	// 		tknComma, makeName("w"), tknRS, makeError(`unexpected right bracket U+005D ']'`),
	// 	},
	// },
	// {"unclosed left bracket",
	// 	"(x+y)*((1/1.324)%4",
	// 	[]token.Token{tknLR, makeName("x"), tknPlus, makeName("y"), tknRR, tknMult,
	// 		tknLR, tknLR, makeToken(token.INT, "1"), tknDiv,
	// 		makeToken(token.FLOAT, "1.324"), tknRR, tknMod, makeToken(token.INT, "4"),
	// 		makeError(`unclosed left bracket: U+0028 '('`),
	// 	},
	// },
	// {"unclosed multiline comment",
	// 	`/* This is an unclosed comment!
	// 	/`,
	// 	[]token.Token{makeError("Multiline comment is not closed")},
	// },
}

func TestLex(t *testing.T) {
	for _, testcase := range lexTests {
		// TODO: Fix test cases
		outputTokens, _ := collect(&testcase)
		if !equal(outputTokens, testcase.tokens, false) {
			t.Errorf("%s: got\n\t%+v\nexpected\n\t%v", testcase.name, outputTokens, testcase.tokens)
		}
	}
}

// Helper Methods to check equality for tests and collect tokens

// collect gathers the emitted items into a Token slice
func collect(tc *lexTestcase) (tkns []token.Token, errs token.ErrorList) {
	l, th := initTestHandler(tc.name, tc.input)
	for {
		tkn := l.Scan()
		tkns = append(tkns, tkn)
		if tkn.Type == token.EOF {
			break
		}
	}
	errs = th.errors
	return
}

func equal(tknLst1, tknLst2 []token.Token, checkPos bool) bool {
	if len(tknLst1) != len(tknLst2) {
		return false
	}
	for k := range tknLst1 {
		tkn1 := tknLst1[k]
		tkn2 := tknLst2[k]
		switch {
		case tkn1.Type != tkn2.Type,
			tkn1.Value != tkn2.Value && !(tkn1.Type == token.SEMICOLON && tkn2.Type == token.SEMICOLON),
			checkPos && tkn1.Pos == tkn2.Pos:
			return false
		}
	}
	return true
}
