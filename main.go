package mathparse

func NewParser(expression string) Parser {
	parse := Parser{}
	parse.ReadExpression(expression)
	return parse
}
