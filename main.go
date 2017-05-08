package mathparse

func NewParser(expression string) Parser {
	parse := Parser{}
	parse.ReadExpression(expression)
	return parse
}

func (p *Parser) FoundResult() bool {
	return len(p.tokens) <= 1 && p.tokens[0].Type == literal
}

func (p *Parser) GetValueResult() float64 {
	return p.tokens[0].ParseValue
}

func (p *Parser) GetExpressionResult() string {
	return getStringExpression(p.tokens)
}

func getStringExpression(set []Token) string {
	str := ""
	for _, tok := range set {
		switch tok.Type {
		case space:
			str += " "
		case literal:
			str += tok.Value
		case variable:
			str += tok.Value
		case operation:
			str += tok.Value
		case function:
			str += tok.Value + "(" + getStringExpression(tok.Children) + ")"
		case lparen:
			str += "(" + getStringExpression(tok.Children) + ")"
		case funcDelim:
			str += ","
		}
	}
	return str
}

func (p *Parser) GetTokens() []Token {
	return p.tokens
}
