package parser

import (
	lexer "fugu/pkg/lexer"
	"fugu/pkg/reporter"
	"fugu/pkg/token"
)

type Parser struct {
	Tokens []token.Token

	lex    *lexer.Lexer
	report *reporter.Reporter

	curToken  token.Token
	peekToken token.Token
}

func New(input, fileName string) *Parser {
	lex := lexer.New(input, fileName)

	pars := &Parser{
		Tokens: make([]token.Token, 0),

		lex: lex,
	}
	pars.report = pars.lex.Report()
	return pars
}
