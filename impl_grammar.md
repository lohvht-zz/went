
expr: orEval;
orEval: andEval ("||" orEval)*;
andEval: notTest ("&&" notTest)*;
notTest: "!" notTest | comparison;
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

atom: "[" [exprList] "]" | <!-- "{" mapList "}" | -->
  ID | NUM | STR | RAWSTR | "null" | "false" | "true";

exprList: expr ("," expr)* [","];
<!-- mapList: keyval ("," keyval)* [","];
keyval: ID ":" expr; -->