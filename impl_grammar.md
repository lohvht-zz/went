


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

# Inputs

```
input: stmt* EOF

```

# Statements
```
stmt: (smallStmt | compoundStmt) ";";
smallStmt: exprStmt | loopCtlStmt | returnStmt;

exprStmt: exprList (augAssign exprList | ('=' exprList)*);
augAssign: "+=" | "-=" | "/=" | "*=" | "%=";

loopCtlStmt: breakStmt | continueStmt;
breakStmt: "break" [exprList];
continueStmt: "continue";
returnStmt: "return" [exprList];

compoundStmt: ifStmt | whileStmt | forStmt | funcDef;
ifStmt: "if" expr block ("elif" expr block)* ["else" block];
whileStmt: "while" expr block ["default" block];
forStmt: "for" exprList "in" expr block ["default" block];

funcDef: "func" NAME parameters ":" block;
parameters: "(" [argList] ")";
argList: NAME (',' NAME)*;

block: "{" stmt+ "}";
```

# Expressions
```
expr: orEval;
orEval: andEval ("||" orEval)*;
andEval: notEval ("&&" notEval)*;
notEval: "!" notEval | comparison;
comparison: smExpr (compOp smExpr)*;
compOp: "==" | "!=" | "<" | ">" | "<=" | ">=" | ["!"] "in";

smExpr: term (addOp term)*; 
term: factor (multOp factor)*;
factor: unaryOp factor | atomExpr;
addOp: "+" | "-";
multOp: "*" | "/" | "%";
unaryOp: "+" | "-";

atomExpr: atom trailer*;
trailer: "(" [argList] ")" | "[" slice "]" | "." NAME;
slice: orEval | [orEval] ":" [orEval] [":" [orEval]];
argList: arg ("," arg)* [","];
arg: orEval | NAME "=" orEval;

<!-- a[1], a[1:2], a[:2], a[1:] -->

// TODO: Map To be implemented
atom: "[" [exprList] "]" | "{" mapList "}" | "(" orEval ")" |
  NAME | NUM | STR | RAWSTR | "null" | "false" | "true";

// TODO: To be implemented
exprList: orEval ("," orEval)* [","];
mapList: keyval ("," keyval)* [","];
keyval: NAME ":" orEval;
```