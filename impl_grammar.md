


<!-- Compound expressions, TODO: replace expr with the statements -->
<!-- For loops and while loops are expressions that evaluate to something
  upon ending the loops,
  1. If loop has no "break" and "default" block => returns null
  2. If loop has "break" but no "default" block, returns the "break" expression,
    else returns the falsy value for the "break" expression (TODO: Define a behaviour for the case where "break" is not defined, what does the expression evaluate?)
  3. If loop has "default" block with no "break", returns the "default" block     expression
  4. else, if loop has both "break" and "default", return their respective values, however the types of both these values must be equal
  (If in the case of multiple breaks, all evaluated values must be the same type)
 -->
whileExpr: "while" expr blockExpr "default" blockExpr
forExpr: "for" exprList "in" expr blockExpr "default" blockExpr
blockExpr: "{" (expr ";")* "}" <!-- TODO: ASI for blockExpr => Follow golang's ASI rule 1 https://golang.org/ref/spec#Semicolons -->

expr: orEval | ifExpr;
ifExpr: "if" expr blockExpr ("elif" expr blockExpr)* ["else" blockExpr];

orEval: andEval ("||" orEval)*;
andEval: notEval ("&&" notEval)*;
notEval: "!" notEval | comparison;
comparison: smExpr (compOp smExpr)*;
compOp: "==" | "!=" | "<" | ">" | "<=" | ">=" | ["!"] "in";

smExpr: term (("+" | "-") term)*;
term: factor (("*" | "/" | "%") factor)*;
factor: ("+" | "-") factor | atomExpr;

atomExpr: atom trailer*;
trailer: "(" [argList] ")" | "[" slice "]" | "." ID;
slice: expr | [expr] ":" [expr] [":" [expr]];
argList: arg ("," arg)* [","];
arg: expr | ID "=" expr;

<!-- a[1], a[1:2], a[:2], a[1:] -->
<!-- "{" mapList "}" | -->
atom: "[" [exprList] "]" | "(" expr ")" |
  ID | NUM | STR | RAWSTR | "null" | "false" | "true";

exprList: expr ("," expr)* [","];
<!-- mapList: keyval ("," keyval)* [","];
keyval: ID ":" expr; -->