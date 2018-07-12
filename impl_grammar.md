
or_test: and_test ("||" or_test)*
and_test: not_test ("&&" not_test)*
not_test: "!" not_test | comparison;
comparison: expr (comp_op expr)*;

expr: term (("+" | "-") term)*;
term: factor (("*" | "/" | "%") factor)*;
factor: ("+" | "-") factor | atom_expr;

atom_expr: atom trailer*;
<!-- Add array, maps into atom -->
atom: ID | NUM | STR | RAWSTR | "null" | "false" | "true";
comp_op: "==" | "!=" | "<" | ">" | "<=" | ">=" | ["!"] "in";
