//go:generate go run generate.go
package action

import (
	"fmt"
	"fugu/pkg/reporter"
	. "fugu/pkg/token"
	"sort"
	"strings"
)

type ActionType int

const (
	_ ActionType = iota

	Shift  // Сдвиг: поместить текущий токен в стек и перейти в новое состояние
	Reduce // Свёртка: заменить последовательность символов по правилу грамматики
	Accept // Успех: входная цепочка полностью распознана
	Error  // Ошибка: для текущего состояния и токена действие не определено
)

type ActionStruct struct {
	Typ     ActionType
	Num     int
	ErrCode reporter.Code
}

type table struct {
	actions []ActionStruct // действия
	check   []int          // проверка состояния
	base    []int          // смещения для состояний
}

// конструкторы
func Sh(state int) ActionStruct {
	return ActionStruct{Typ: Shift, Num: state, ErrCode: reporter.NoError}
}
func Red(rule int) ActionStruct {
	return ActionStruct{Typ: Reduce, Num: rule, ErrCode: reporter.NoError}
}
func Acc() ActionStruct                { return ActionStruct{Typ: Accept, ErrCode: reporter.NoError} }
func Err(e reporter.Code) ActionStruct { return ActionStruct{Typ: Error, ErrCode: e} }

func BuildActionSlice(src *map[int]map[TokenKind]ActionStruct) *table {
	ms := 0
	for s := range *src {
		if s > ms {
			ms = s
		}
	}
	states := ms + 1

	type couple struct {
		tk  int
		act ActionStruct
	}
	rows := make([][]couple, states)
	for state, row := range *src {
		entries := make([]couple, 0, len(row))
		for tk, act := range row {
			entries = append(entries, couple{int(tk), act})
		}
		sort.Slice(entries, func(i, j int) bool { return entries[i].tk < entries[j].tk })
		rows[state] = entries
	}

	var actions []ActionStruct
	var check []int
	base := make([]int, states)
	for i := range base {
		base[i] = -1
	}

	isFree := func(idx int) bool {
		return idx >= len(check) || check[idx] == -1
	}

	order := make([]int, 0, states)
	for i := 0; i < states; i++ {
		if len(rows[i]) > 0 {
			order = append(order, i)
		}
	}
	sort.Slice(order, func(i, j int) bool { return len(rows[order[i]]) > len(rows[order[j]]) })

	for _, state := range order {
		entries := rows[state]
		for b := 0; ; b++ {
			ok := true
			for _, e := range entries {
				if !isFree(b + e.tk) {
					ok = false
					break
				}
			}
			if ok {
				base[state] = b
				maxIdx := b + entries[len(entries)-1].tk
				for len(actions) <= maxIdx {
					actions = append(actions, ActionStruct{Typ: Error, ErrCode: reporter.NoError})
					check = append(check, -1)
				}
				for _, e := range entries {
					idx := b + e.tk
					actions[idx] = e.act
					check[idx] = state
				}
				break
			}
		}
	}

	return &table{
		actions: actions,
		check:   check,
		base:    base,
	}
}

func GenerateActionTable(src *map[int]map[TokenKind]ActionStruct, tokenCount int) string {
	table := BuildActionSlice(src)

	var out strings.Builder
	out.WriteString("//! DO NOT EDIT\n")
	out.WriteString("package action\n\n")
	out.WriteString("import (\n")
	out.WriteString("\t\"fugu/pkg/reporter\"\n")
	out.WriteString("\t. \"fugu/pkg/token\"\n")
	out.WriteString(")\n\n")
	out.WriteString("var Actions = []ActionStruct{\n")
	for i, act := range table.actions {
		out.WriteString(fmt.Sprintf("\t%v, // %d\n", renderActionType(act), i))
	}
	out.WriteString("}\n\n")
	out.WriteString("var Check = []int{\n")
	for i, c := range table.check {
		out.WriteString(fmt.Sprintf("\t%d, // %d\n", c, i))
	}
	out.WriteString("}\n\n")
	out.WriteString("var Base = []int{\n")
	for i, b := range table.base {
		out.WriteString(fmt.Sprintf("\t%d, // state %d\n", b, i))
	}
	out.WriteString("}\n\n")
	out.WriteString("func Action(state int, tk TokenKind) ActionStruct {\n")
	out.WriteString("\tif state < 0 || state >= len(Base) {\n")
	out.WriteString("\t\treturn Err(reporter.NoError)\n")
	out.WriteString("\t}\n")
	out.WriteString("\tb := Base[state]\n")
	out.WriteString("\tif b < 0 {\n")
	out.WriteString("\t\treturn Err(reporter.NoError)\n")
	out.WriteString("\t}\n")
	out.WriteString("\tidx := b + int(tk)\n")
	out.WriteString("\tif idx >= 0 && idx < len(Actions) && Check[idx] == state {\n")
	out.WriteString("\t\treturn Actions[idx]\n")
	out.WriteString("\t}\n")
	out.WriteString("\treturn Err(reporter.NoError)\n")
	out.WriteString("}\n")
	return out.String()
}

func renderActionType(a ActionStruct) string {
	switch a.Typ {
	case Shift:
		return fmt.Sprintf("Sh(%d)", a.Num)
	case Reduce:
		return fmt.Sprintf("Red(%d)", a.Num)
	case Accept:
		return "Acc()"
	case Error:
		return fmt.Sprintf("Err(%s)", "reporter."+a.ErrCode.Code())
	}
	return "Err(reporter.NoError)"
}
