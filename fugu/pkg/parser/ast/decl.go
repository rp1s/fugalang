package ast

import "fugu/pkg/token"

type Variable struct {
	Token    token.Token // токен LET
	Name     *Identifier // имя переменной
	Type     Expression  // опциональный тип
	Value    Expression  // начальное значение
	Mutable  bool        // можно ли изменять
	Constant bool        // константа
	Exported bool        // публичная
	IsGlobal bool        // глобальная переменная
}

func (v *Variable) Node()
func (v *Variable) Decl()

type Struct struct {
	Token  token.Token    // token.STRUCT
	Name   *Identifier    // имя структуры
	Fields []*StructField // поля структуры
	// TODO: Methods  []*FunctionDeclaration // методы структуры
	Exported bool            // публичная
	Generic  []*GenericParam // параметры дженериков
}

type StructField struct {
	Token token.Token // token.IDENTIFIER
	Name  *Identifier // имя поля
	Type  Expression  // тип поля
	// TODO: Tags  структуры
	Exported bool // публичное поле
}

func (s *Struct) Node()
func (s *Struct) Decl()
