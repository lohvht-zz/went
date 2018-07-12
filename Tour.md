# A tour of \<nondescript\>

## Hello World!
The typical "Hello World" is written as follows:
```
echo('Hello nondescript')
```

## Data Types
nondescript supports values of familiar types like `int`, `float`, `string` and `bool`. nondescript supports `lists` in the form of `int` indexed ordered lists as well as `maps` in the form of key value pairs (with keys being `strings`).

`maps` and `lists` in nondescript may store any type as well.
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
A comment in nondescript starts with `//` ending with the end of line, or has the form of `/* anything and everything inluding newlines! */`. All comments are ignored by the interpreter

## Variable assignments
Variables are assigned in the following way:
```
message = 'Hello world'
```

You can assign multiple variables at once by separating them with semicolons `;`.

```
a = 42; b = 'This is a string'; c = false
```

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

## Defining functions
functions are defined using the `func` keyword.
```
func sum(a, b) {
  return a + b
}
```

Funtion parameters are separated by commas `,` and enclosed in parenthesis `(a, b)`, each parameter is defined by its name. The `return` keyword returns the function call.

Functions are able to return multiple values, by specifying the number of return types after the parenthesis, or by stating the types of the return results in parenthesis.
```
a = [1, 2, "string1"]
func unpack(array): 3 {
  return a[0], a[1], a[2]
}
val1, val2, val3 = unpack(a)

// Alternatively, if you want to enforce typing
func unpack2(array): int; int; string; {
  return a[0], a[1], a[2]
}
val4, val5, val6 = unpack2(a)
```

## Operators
Binary operators `+`, `-`, `*`, `/`, `%` to perform operation between 2 `numbers`. Unary operations `+`, `-`, `!`.

Comparators such as `<`, `>`, `==`, `>=`, `<=`, `!=` can be used to compare things in nondescript, and it will always try to evaluate based on the value of the constructs that its comparing.

Comparators like "not in" or "in" checks if object exists
```
echo(1 + 2 < 4) // true
arr1 = [1, 2, 3, 4]; 
```
