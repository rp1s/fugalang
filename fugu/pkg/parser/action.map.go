package parser

import (
	"fugu/pkg/reporter"
	. "fugu/pkg/token"
)

type ActionType int

const (
	_ ActionType = iota
	Shift
	Reduce
	Accept
	Error
)

type Action struct {
	Typ     ActionType
	Num     int
	ErrCode reporter.Code
}

// конструкторы
func Sh(state int) Action        { return Action{Typ: Shift, Num: state, ErrCode: reporter.NoError} }
func Red(rule int) Action        { return Action{Typ: Reduce, Num: rule, ErrCode: reporter.NoError} }
func Acc() Action                { return Action{Typ: Accept, ErrCode: reporter.NoError} }
func Err(e reporter.Code) Action { return Action{Typ: Error, ErrCode: e} }

// MAP ACTION

var action = map[int]map[TokenKind]Action{
	// Старт
	0: {
		INTEGER: Sh(5),
		L_PAREN: Sh(4),
	},
}
