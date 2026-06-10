package ast

import "fugu/pkg/token"

// обёртка для decl
type DeclStmt struct {
	Decl Declaration
}

func (ds *DeclStmt) Node()
func (ds *DeclStmt) Stmt()

type Variable struct {
	Token    token.Token // токен LET
	Name     *Identifier // имя переменной
	Type     Expression  // опциональный тип
	Value    Expression  // начальное значение
	Mutable  bool        // можно ли изменять
	Constant bool        // константа
	Exported bool        // публичная
	IsGlobal bool        // глобальная переменная

	Docs string
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

	Docs string
}

type StructField struct {
	Token token.Token // token.IDENTIFIER
	Name  *Identifier // имя поля
	Type  Expression  // тип поля
	// TODO: Tags  структуры
	Exported bool // публичное поле

	Docs string
}

func (s *Struct) Node()
func (s *Struct) Decl()

// struct Box<T> { value: T }
type GenericParam struct {
	Token      token.Token // token.IDENTIFIER
	Name       *Identifier // имя параметра (T)
	Constraint Expression  // ограничение (опционально)
}

// FN (metods struct) Name[generic](arg) -> (returns) {body}
type Functions struct {
	Token   token.Token // token.FN
	Name    *Identifier
	Body    *BlockStmt
	Generic []*GenericParam // параметры дженериков
}

type Return struct {
	L_PAREN token.Token

	R_PAREN token.Token
}

type BlockStmt struct {
	L_BRACE token.Token
	Smts    []Statement
	R_BRACE token.Token
}
