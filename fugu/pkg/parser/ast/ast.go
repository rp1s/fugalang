package ast

import "fugu/pkg/token"

//
//
// Expressions - Выражения
// Statements - Инструкции
// Declarations - Объявления
//
//

type Node interface {
	Literal(source token.OutLiteral) string
}

type Statement interface {
	Node
	stmt()
}

type Expression interface {
	Node
	expr()
}

type Declaration interface {
	Node
	decl()
}

type Program struct {
	Declaration []Declaration
}

func (p *Program) Literal(source token.OutLiteral) string {
	if len(p.Declaration) > 0 {
		return p.Declaration[0].Literal(source)
	}
	return ""
}
