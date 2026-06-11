package action

import (
	. "fugu/pkg/token"
)

var ActionSrc = &map[int]map[TokenKind]ActionStruct{
	// State 0 - стартовое состояние
	0: {
		GNUMBER: Sh(2),
		L_PAREN: Sh(3),
	},

	/*
		// State 1 - после полной обработки выражения (ожидаем EOF)
		1: {
			EOF: Acc(),
		},

		// State 2 - после числа (GNUMBER)
		2: {
			GARITHMETIC: Red(2), // E -> GNUMBER
			R_PAREN:     Red(2),
			EOF:         Red(2),
		},

		// State 3 - после '('
		3: {
			GNUMBER: Sh(2),
			L_PAREN: Sh(3),
		},

		// State 4 - после оператора (GARITHMETIC)
		4: {
			GNUMBER: Sh(2),
			L_PAREN: Sh(3),
		},

		// State 5 - после выражения внутри скобок (E)
		5: {
			GARITHMETIC: Sh(4),
			R_PAREN:     Sh(6),
			EOF:         Red(3), // если нужно
		},

		// State 6 - после ')'
		6: {
			GARITHMETIC: Red(3), // E -> ( E )
			R_PAREN:     Red(3),
			EOF:         Red(3),
		},

		// State 7 - после E op E (для редукции бинарной операции)
		7: {
			GARITHMETIC: Red(1), // E -> E GARITHMETIC E
			R_PAREN:     Red(1),
			EOF:         Red(1),
		},
	*/
}
