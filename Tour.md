# A tour of \<`Went`\>

## Hello World!
The typical "Hello World" is written as follows:
```
echo('Hello `Went`')
```

## Data Types
`Went` supports values of familiar types like `int`, `float`, `string` and `bool`. `Went` supports `lists` in the form of `int` indexed ordered lists as well as `maps` in the form of key value pairs (with keys being `strings`).

`maps` and `lists` in `Went` may store the above types as well as other `maps` or `lists` as well.
<!-- as well as `structs`, which can hold functions called `methods` as well as `properties` that can hold any of the data types defined above. -->

## Zero Values
Here are the Zero values of the following data types:
- `int`: `0`
- `float`: `0.0`
- `string`: `''`
- `bool`: `false`
- `list`: `[]`
- `map`: `{}`
- `null`: `null`

`null` is a special case.

When tested for truth value using conditionals such as during `if` or `while`, or as an operand of boolean operation, the above values will always evaluate as `false`.


## Comments
A comment in `Went` starts with `//` ending with the end of line, or has the form of `/* anything and everything inluding newlines! */`. All comments are ignored by the interpreter

## Variable assignments
Variables are assigned in the following way:
```
message = 'Hello world'
```

You can assign multiple variables at once by separating them with semicolons `;`.

```
a = 42; b = 'This is a string'; c = false
```

## Statements


## Automatic Semicolon Insertion


## Conditionals
Conditions are expressed using an `if`, `elif`, `else` *expression*.

```
a = 0
if a > 0 {
  echo('Greater than 0')
} elif a < 0 {
  echo('Smaller than 0')
} else {
  echo('Perfectly balanced, as all things should be')
}
```

Condition must evaluate to a value of type `bool`, the `else` and `elif` blocks may be omitted if needed.

Since the conditional is an *expression*, it returns a value as well. Thus there is no ternary operators `(condition) ? then : else`. The code below illustrates so.
```
// Traditionally
max = a
if a < b { max = b }

// Using else
max
if a > b {
  max = a
} else {
  max = b
}

// As expression
max = if a > b { a } else { b }
```

`if` branches can be blocks, and the last expression is the value of a block.
```
a = 1
b = 2
max = if a < b {
  echo("Choose b")
  // Some more statements ...
  b
} elif a > b {
  echo("Choose a")
  // Some more statements ...
  a
} else {
  0
}
```

When using `if` as an *expression*, take note that you should also implement the `else` clause as well, else the value will default to the *zero value* of the type of the final expression within the `if` block.

## Loops
A `for` loop iterates through items in a collection, such as `lists` or `maps`. The syntax is as follows:
```
for item in collection {
  echo(item)
}
```

If you require access to the `key` or `index` for `maps` or `lists` respectively, you may use the following as well.

```
for item, i in collection {
  echo("The item is", item, " with key/index ", i)
}
```

To iterate over a range of numbers, use normal C style iteration
```
for n=0; n<10; n++ {
  echo(n)
}
```

You can also use `while` loops to execute expressions as long as a condition holds true.
```
a = 100
while a != 1 {
  echo(a)
  a--
}
```

For both `for` and `while` loops, the usual `break` and `continue` can also still be used.

### Break returns
For all loops, we can also specify return values, thus `for` and `while` loops are also expressions.

```
array1 = [1, 3, 4, 6, 8, 20, 30]
foundItem = for item in collection {
  if item > 20 {
    break item
  }
}
```

Since some loops can run to completion without breaking, the default return value for a loop construct is the zero-value of the types of the returned values. 

For example if we have `break intValue, boolValue`, where `intValue` is of type `int` and `boolValue` is of type `bool`, for a loop that will terminate normally then the corresponding return for values will be `0, false`.

If we would like to assign values for a loop-expression that finishes normally, we can use a `default` block.

```
map = {key1: "Hello World", key2: "How are you doing?", key3: -4}
foundVal = for key, val in map {
  if val == "Woopty Doo" {
    break val
  }
} default {
  "Not found"
}
```

In the case above, the loop does not terminate with a break, and so foundVal will be set to `"Not found"`.

The `default` block works the same way as most blocks in `Went`, where blocks are expressions that return the value of their last statement.


## Defining functions
functions are defined using the `func` keyword.
```
func sum(a, b) {
  return a + b
}
```

Funtion parameters are separated by commas `,` and enclosed in parenthesis `(a, b)`, each parameter is defined by its name. The `return` keyword returns the function call.

Functions are able to return multiple values, as well.
```
a = [1, 2, "string1"]
func unpack(array) {
  return a[0], a[1], a[2]
}
val1, val2, val3 = unpack(a)
```

## Operators
Binary operators `+`, `-`, `*`, `/`, `%` to perform operation between 2 `numbers`. Unary operations `+`, `-`, `!`.

Comparators such as `<`, `>`, `==`, `>=`, `<=`, `!=` can be used to compare things in `Went`, and it will always try to evaluate based on the value of the constructs that its comparing.

```
echo(1 + 2 < 4) // true
```

Comparators like "not in" or "in" checks if a value exists inside an array
```
arr1 = [1, 2, 3, 4]
if 3 in arr1 {
  echo('Hooray, found 3!') // will echo
}

if -1 not in arr1 {
  echo('Boo, can't take your negativity here!')
}
```

# Execution Model

## Structure of a program
A *`went`* program is constructed from code `block`s. A `block` is a piece of `went` program text that is executed as a single unit. The following are considered `block`s: a function body, a module.

Each command typed interactively (through the `interpreter` shell) is a block. A script file (i.e. a file given as an input from the command line argument to the `interpreter`) is a code block (NOTE: if a file has a live command typed, it will pause the execution in the same scope)

A code block is executed in an *execution frame*. A frame contains some administrative information (used for debugging) and determines where and how execution continues after the code block's execution has completed.

## Naming and Binding
### Nam