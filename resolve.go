package mathparse

import (
	"math"
	"strconv"
)

func (p *Parser) Resolve() {
	// parenthases
	// exponents/roots
	// multiplication/division
	// addition/subtraction
	// functions
	// repeat
	p.tokens = parseExpression(p.tokens)
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

func parseExpression(set []Token) []Token {
	mod := false
	if set[0].Type == function || set[0].Type == lparen {
		set[0].Children = parseExpression(set[0].Children)
	}
	for i := 1; i < len(set)-1; i++ {
		if set[i].Type == operation {
			if (set[i].Value == "^") && (set[i-1].Type == literal && set[i+1].Type == literal) {
				mod = true
				set[i-1].ParseValue = math.Pow(set[i-1].ParseValue, set[i+1].ParseValue)
				set[i-1].Value = strconv.FormatFloat(set[i-1].ParseValue, 'f', -1, 64)
				set = append(set[:i], set[i+2:]...)
				i--
			}
		}
	}
	for i := 1; i < len(set)-1; i++ {
		if set[i].Type == operation {
			if (set[i].Value == "*" || set[i].Value == "/") && (set[i-1].Type == literal && set[i+1].Type == literal) {
				mod = true
				if set[i].Value == "*" {
					set[i-1].ParseValue = set[i-1].ParseValue * set[i+1].ParseValue
				} else if set[i].Value == "/" {
					set[i-1].ParseValue = set[i-1].ParseValue / set[i+1].ParseValue
				}
				set[i-1].Value = strconv.FormatFloat(set[i-1].ParseValue, 'f', -1, 64)
				set = append(set[:i], set[i+2:]...)
				i--
			}
		}
	}

	for i := 1; i < len(set)-1; i++ {
		if set[i].Type == operation {
			if (set[i].Value == "+" || set[i].Value == "-") && (set[i-1].Type == literal && set[i+1].Type == literal) {
				mod = true
				if set[i].Value == "+" {
					set[i-1].ParseValue = set[i-1].ParseValue + set[i+1].ParseValue
				} else if set[i].Value == "-" {
					set[i-1].ParseValue = set[i-1].ParseValue - set[i+1].ParseValue
				}
				set[i-1].Value = strconv.FormatFloat(set[i-1].ParseValue, 'f', -1, 64)
				set = append(set[:i], set[i+2:]...)
				i--
			}
		}
	}

	// functions
	for i := range set {
		if set[i].Type == lparen {
			if len(set[i].Children) == 1 {
				set[i] = set[i].Children[0]
			}
		}
		if set[i].Type == function {
			mod = true
			switch set[i].Value {
			case "sin":
				set[i] = newToken(literal, strconv.FormatFloat(math.Sin(set[i].Children[0].ParseValue), 'f', -1, 64))
			case "cos":
				set[i] = newToken(literal, strconv.FormatFloat(math.Cos(set[i].Children[0].ParseValue), 'f', -1, 64))
			case "tan":
				set[i] = newToken(literal, strconv.FormatFloat(math.Tan(set[i].Children[0].ParseValue), 'f', -1, 64))
			case "abs":
				set[i] = newToken(literal, strconv.FormatFloat(math.Abs(set[i].Children[0].ParseValue), 'f', -1, 64))
			case "log":
				set[i] = newToken(literal, strconv.FormatFloat(math.Log(set[i].Children[0].ParseValue), 'f', -1, 64))
			case "sqrt":
				set[i] = newToken(literal, strconv.FormatFloat(math.Sqrt(set[i].Children[0].ParseValue), 'f', -1, 64))
			case "max":
				set[i] = newToken(literal, strconv.FormatFloat(math.Max(set[i].Children[0].ParseValue, set[i].Children[2].ParseValue), 'f', -1, 64))
			case "min":
				set[i] = newToken(literal, strconv.FormatFloat(math.Min(set[i].Children[0].ParseValue, set[i].Children[2].ParseValue), 'f', -1, 64))
			case "mod":
				set[i] = newToken(literal, strconv.FormatFloat(math.Mod(set[i].Children[0].ParseValue, set[i].Children[2].ParseValue), 'f', -1, 64))
			case "pow":
				set[i] = newToken(literal, strconv.FormatFloat(math.Pow(set[i].Children[0].ParseValue, set[i].Children[2].ParseValue), 'f', -1, 64))
			}
		}
	}

	if len(set) > 1 && mod {
		set = parseExpression(set)
	}

	return set
}
