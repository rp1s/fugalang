//go:generate go run generate.go
package parser

import (
	"fmt"
	"fugu/pkg/reporter"
	. "fugu/pkg/token"
)

var action []Action

type ActionType int

const (
	_ ActionType = iota

	Shift  // Сдвиг: поместить текущий токен в стек и перейти в новое состояние
	Reduce // Свёртка: заменить последовательность символов по правилу грамматики
	Accept // Успех: входная цепочка полностью распознана
	Error  // Ошибка: для текущего состояния и токена действие не определено
)

type Action struct {
	Typ     ActionType
	Num     int
	ErrCode reporter.Code
}

type Table struct {
	actions []Action
	tkCount int
}

func (t *Table) Get(state int, tk TokenKind) Action {
	fmt.Println(state*t.tkCount + int(tk))
	return t.actions[state*t.tkCount+int(tk.Group())]
}

// конструкторы
func Sh(state int) Action        { return Action{Typ: Shift, Num: state, ErrCode: reporter.NoError} }
func Red(rule int) Action        { return Action{Typ: Reduce, Num: rule, ErrCode: reporter.NoError} }
func Acc() Action                { return Action{Typ: Accept, ErrCode: reporter.NoError} }
func Err(e reporter.Code) Action { return Action{Typ: Error, ErrCode: e} }

func BuildActionSlice(src map[int]map[TokenKind]Action, tokenCount int) *Table {
	maxState := 0
	for s := range src {
		if s > maxState {
			maxState = s
		}
	}

	t := &Table{
		actions: make([]Action, (maxState+1)*tokenCount),
		tkCount: tokenCount,
	}

	for state, row := range src {
		base := state * tokenCount
		for tk, act := range row {
			t.actions[base+int(tk)] = act
		}
	}

	return t
}

func GenerateInitSlice(src map[int]map[TokenKind]Action, tokenCount int) string {
	maxState := 0
	for s := range src {
		if s > maxState {
			maxState = s
		}
	}

	size := (maxState + 1) * tokenCount

	table := make([]Action, size)

	for state, row := range src {
		base := state * tokenCount
		for tk, act := range row {
			table[base+int(tk)] = act
		}
	}

	out := "import \"fugu/pkg/reporter\"\n\n"
	out += "func InitActionSlice() {\n"
	out += "\taction = []Action{\n"

	for i, act := range table {
		out += fmt.Sprintf("\t\t%v, // %d\n", renderActionType(act), i)
	}

	out += "\t}\n"
	out += "}\n"

	return out
}

func renderActionType(a Action) string {
	switch a.Typ {
	case Shift:
		return fmt.Sprintf("Sh(%d)", a.Num)
	case Reduce:
		return fmt.Sprintf("Red(%d)", a.Num)
	case Accept:
		return "Acc()"
	case Error:
		return fmt.Sprintf("Err(%s)", a.ErrCode.Code())
	}
	return "Err(reporter.NoError)"
}

var ActionSrc = map[int]map[TokenKind]Action{
	// Точка старта
	0: {
		GNUMBER: Sh(2),
		L_PAREN: Sh(3),
	},
	1: {
		GARITHMETIC: Sh(4),
		EOF:         Acc(),
	},
	2: {
		GARITHMETIC: Red(6),
		R_PAREN:     Red(6),
		EOF:         Red(6),
	},
	3: {
		GNUMBER: Sh(2),
		L_PAREN: Sh(3),
	},
	4: {
		GNUMBER: Sh(2),
		L_PAREN: Sh(3),
	},
}
