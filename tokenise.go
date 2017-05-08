package mathparse

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Parser struct {
	letterBuffer string
	numberBuffer string
	expression   string
	tokens       []Token
}

type Token struct {
	Type       TokenType
	Value      string
	ParseValue float64
	Children   []Token
}

type TokenType uint

const (
	undefined TokenType = iota // 0
	space                      // 1
	literal                    // 2
	variable                   // 3
	operation                  // 4
	function                   // 5
	lparen                     // 6
	rparen                     // 7
	funcDelim                  // 8
)

// "89sin(45) + 2.2x/7"

func (p *Parser) ReadExpression(str string) {
	p.expression = str
	p.tokens = []Token{}
	p.tokenise()
}

func (p *Parser) ReadMultipartExpression(str []string) {
	p.expression = strings.Join(str, " ")
	p.tokens = []Token{}
	p.tokenise()
}

func (p *Parser) tokenise() {
	dumpLetter := func(p *Parser) {
		for i, ch := range p.letterBuffer {
			p.tokens = append(p.tokens, newToken(variable, string(ch)))
			if i < len(p.letterBuffer)-1 {
				p.tokens = append(p.tokens, newToken(operation, "*"))
			}
		}
		p.letterBuffer = ""
	}
	dumpNumber := func(p *Parser) {
		if len(p.numberBuffer) > 0 {
			p.tokens = append(p.tokens, newToken(literal, p.numberBuffer))
			p.numberBuffer = ""
		}
	}
	for _, ch := range p.expression {
		switch getTokenType(ch) {
		case space:
			continue
		case literal:
			p.numberBuffer += string(ch)
		case variable:
			if len(p.numberBuffer) > 0 {
				dumpNumber(p)
				p.tokens = append(p.tokens, newToken(operation, "*"))
			}
			p.letterBuffer += string(ch)
		case operation:
			dumpNumber(p)
			dumpLetter(p)
			p.tokens = append(p.tokens, newToken(operation, string(ch)))
		case lparen:
			if len(p.numberBuffer) > 0 {
				dumpNumber(p)
				p.tokens = append(p.tokens, newToken(operation, "*"))
			}
			if len(p.letterBuffer) > 0 {
				p.tokens = append(p.tokens, newToken(function, p.letterBuffer))
				p.letterBuffer = ""
			}
			p.tokens = append(p.tokens, newToken(lparen, "("))
		case rparen:
			dumpLetter(p)
			dumpNumber(p)
			p.tokens = append(p.tokens, newToken(rparen, ")"))
		case funcDelim:
			dumpNumber(p)
			dumpLetter(p)
			p.tokens = append(p.tokens, newToken(funcDelim, ","))
		}
	}

	if len(p.numberBuffer) > 0 {
		dumpNumber(p)
	}
	if len(p.letterBuffer) > 0 {
		dumpLetter(p)
	}

	p.tokens, _ = buildTree(p.tokens)
}

func buildTree(set []Token) ([]Token, int) {
	toks := []Token{}
	for i := 0; i < len(set); i++ {
		tok := set[i]
		switch tok.Type {
		case function:
			child, offset := buildTree(set[i+2:])
			tok.Children = child
			i += 2 + offset
			toks = append(toks, tok)
		case lparen:
			child, offset := buildTree(set[i+1:])
			tok.Children = child
			toks = append(toks, tok)
			i += 1 + offset
		case rparen:
			// toks = append(toks, tok)
			return toks, i
		default:
			toks = append(toks, tok)
		}
	}
	return toks, len(set)
}

func newToken(typ TokenType, value string) Token {
	tok := Token{
		Type:  typ,
		Value: value,
	}
	if typ == literal {
		tok.ParseValue, _ = strconv.ParseFloat(value, 64)
	}
	return tok
}

func getTokenType(ch rune) TokenType {
	let := string(ch)
	if let == " " {
		return space
	} else if isDigit(let) {
		return literal
	} else if isLetter(let) {
		return variable
	} else if isOperator(let) {
		return operation
	} else if isOpenParen(let) {
		return lparen
	} else if isCloseParen(let) {
		return rparen
	}
	return undefined
}

func isComma(let string) bool {
	return let == ","
}

func isDigit(let string) bool {
	res, err := regexp.MatchString(`[0-9\.]`, let)
	if err != nil {
		fmt.Print(err)
	}
	return res
}

func isLetter(let string) bool {
	res, err := regexp.MatchString("[a-zA-Z]", let)
	if err != nil {
		fmt.Print(err)
	}
	return res
}

func isOperator(let string) bool {
	res, err := regexp.MatchString(`\*|\/|\+|\^|\-`, let)
	// "\x43|\x47|\x42|\x94|\x45"
	if err != nil {
		fmt.Print(err)
	}
	return res
}

func isOpenParen(let string) bool {
	return let == "("
}

func isCloseParen(let string) bool {
	return let == ")"
}
