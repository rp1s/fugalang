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

	pos int
}

func New(input []byte, fileName string) *Parser {
	lex := lexer.New(input, fileName)

	pars := &Parser{
		Tokens: make([]token.Token, 0),

		lex: lex,
	}
	pars.report = pars.lex.Report()
	if pars.report.IsUse {
		return nil
	}

	pars.advance().advance()
	if pars.curToken.Kind == token.EOF {
		pars.report.SendTk(reporter.LexerNoClosing, pars.curToken)
		return pars // TODO ошибку выкинуть
	}

	return pars
}

func (ps *Parser) advance() *Parser {
	tk := ps.lex.NextToken()
	ps.curToken = ps.peekToken
	ps.peekToken = tk

	if len(ps.Tokens) == 0 || ps.Tokens[len(ps.Tokens)-1].Kind != token.EOF {
		ps.Tokens = append(ps.Tokens, ps.curToken)
	}
	ps.pos++

	return ps
}
