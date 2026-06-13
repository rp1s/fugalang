package parser

import (
	lexer "fugu/pkg/lexer"
	"fugu/pkg/parser/action"
	"fugu/pkg/reporter"
	"fugu/pkg/token"
)

type Parser struct {
	Tokens []token.Token

	lex    *lexer.Lexer
	report *reporter.Reporter

	curToken token.Token

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
		pars.report.SendTk(reporter.ParserCantStartWork, pars.curToken)
		return pars
	}

	pars.advance()
	if pars.curToken.Kind == token.EOF {
		pars.report.SendTk(reporter.ParserCantStartWork, pars.curToken)
		return pars
	}

	return pars
}

func (ps *Parser) advance() *Parser {
	ps.curToken = ps.lex.NextToken()

	if len(ps.Tokens) == 0 || ps.Tokens[len(ps.Tokens)-1].Kind != token.EOF {
		ps.Tokens = append(ps.Tokens, ps.curToken)
	}
	ps.pos++

	return ps
}

func (ps *Parser) Run() {
	var state int = 0
	for {
		as := action.Action(state, ps.curToken.Kind)

		switch as.Typ {
		case action.Accept:
			break
		case action.Reduce:
			break
		case action.Shift:
			break
		case action.Error:
			break
		}
		ps.advance()
	}
}
