package ast

import "fugu/pkg/token"

//
// Справка:
// Expressions - Выражения (значение, которое вычисляется и возвращает результат)
// Statements - Инструкции (действие, которое что-то делает, но не возвращает значение)
// Declarations - Объявления (создание новой сущности: переменную, функцию)
//

type Node interface {
	Node()
}

type Statement interface {
	Node
	Stmt()
}

type Expression interface {
	Node
	Expr()
}

type Declaration interface {
	Node
	Decl()
}

type Program struct {
	Declaration []Declaration
}

type Identifier struct {
	Token token.Token
}
