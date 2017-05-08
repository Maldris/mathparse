# mathparse
golang library for parsing maths expression strings

## Purpose
Built for a personal project, to help simplify use of maths in text/template, mainly to avoid doing everything with nested function calls in template. i.e.:
`slope(add(sub(1 3) 12) 33 div(24 3.976) 41)`
Which can now be written as
`slope((1-3)+12, 33, 24/3.976, 41)`
Where this library is used inside the slope function (or other function, its just an example)
Which I find will make larger, more complex expressions easier to read

## Disclaimer
Not horribly efficent but perfectly servicable, advice welcomed, but this isnt intended to be a publicly maintained library, purely for personal use and the learning experience.

## limitations
Any variables used in expressions can only be single letter, i.e.: `ma` is evaluated as `m*a`
This was chosen as it allowed parsing with implicit multiplication (useful as people somtimes forget the \*, being used to maths notation) its intended use case is in text/templates templating, in which case variables are substituted for values prior to this library being called
Most other use cases I am likely to use this for will be maths related, where there is a convention of single letter variable names

## Usage
### Parser
The core functionality is provided by the Parser object, if you have the expression ready, you can create the parser, and process your expression at the same time with
```
  expression := "89sin(1.57) + 2.2(31)/7"
  p := mathparse.NewParser(expression)
```
At this point the expression is tokenised and ready to parse
This process is separated, and the Token class exported, to allow people to build other resolving logic if they so wish.

To resolve the expression, call `Resolve`
```
  p.Resolve()
```

When this is done, there are two possible results, either the expression has resolved down to a single numeric value, and can be output as a float64, or (due to variables in use, or an unknown function) a potentially simplified expression string can be retreived.
To know which Option to use, check if the expression is a value with `FoundResult`, a return value of true means a float value can be retrieved, otherwise, an expression
```
  if p.FoundResult() {
    var result float64
    result = p.GetValueResult()
    log.Print(result)
  } else {
    var expression string
    expression = p.GetExpressionResult()
    log.Print(expression)
  }
```
Here `GetValueResult` retreives the float result of the expression, and `GetExpressionResult` will return the expressions simplified form


Its worth noting that the parser object is reusable if need be, if after parsing one expression you wish to parse another, simply load your next one with either `ReadExpression` or `ReadMultipartExpression`.
Each of which will read the expression in and tokenise it as `NewParser` does, but on an existing parser.
`ReadMultipartExpression` exists so that if you have an expression already in multiple parts (i.e. separated as function inputs to a text/template function call) the library is still simple to use. It will simple concatenate the expression segments, and proceed to attempt to resolve the resultant expression.

If multiple segment data is intended, but not desired as a string expression, the raw tokens can be retreived via `GetTokens`, which will return the raw token tree from the parser.
Useful if you have funciton aruements separated by commas, and you want to then take the result and pass into an external function.


### tokens
The structure of each token is quite simple:
```
  type Token struct {
    Type       TokenType
    Value      string
    ParseValue float64
    Children   []Token
  }
```
Type is the type of token, from an enum (see below), value, which is the raw string value of the token (i.e. "3.8", "+", "a", "sin", etc), ParseValue, which is the float of the value, if the type is a literal. And lastly, Children, which will contain child tokens nested under this token, which is only the case for functions, and Parenthesis.

TokenType may take on the following values:
```
const (
  undefined TokenType = iota // 0 - unknown token character
  space                      // 1 - space character, ignored
  literal                    // 2 - a literal, a number
  variable                   // 3 - variables
  operation                  // 4 - any of the following mathematical operations: * / + - ^
  function                   // 5 - a function, it will have the expression for its arguements as Child tokens
  lparen                     // 6 - opening parenthasis, will have the enclosed expression as Child tokens
  rparen                     // 7 - closing parenthasis, used internally, stripped in tree creation, used to mark the end of the current function or parenthasis
  funcDelim                  // 8 - delimits function arguements, doesnt do anything, but prevents, adjacent expressions being evaluated together
)
```
