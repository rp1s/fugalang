package action

import (
	. "fugu/pkg/token"
)

var ActionSrc = map[int]map[TokenKind]ActionStruct{
	0: { // стартовое
		GLITERAL: Sh(1),
		L_PAREN:  Sh(2),
	},

	1: { // после литерала
		GARITHMETIC: Red(1),
		R_PAREN:     Red(1),
		EOF:         Red(1),
	},

	2: { // после (
		GLITERAL: Sh(1),
		L_PAREN:  Sh(2),
	},

	3: { // после Выражение + - || eof
		INCREASE:  Red(3),
		DECREASE:  Red(3),
		MULTIPLY:  Sh(4),
		DIVIDE:    Sh(4),
		REMAINDER: Sh(4),
		DEGREE:    Sh(5),
		R_PAREN:   Red(3),
		EOF:       Red(3),
	},

	4: { // после * / % все трое на одном уровне
		INCREASE:  Red(4),
		DECREASE:  Red(4),
		MULTIPLY:  Sh(4),
		DIVIDE:    Sh(4),
		REMAINDER: Sh(4),
		DEGREE:    Sh(5),
		R_PAREN:   Red(4),
		EOF:       Red(4),
	},

	5: {
		INCREASE:  Red(5),
		DECREASE:  Red(5),
		MULTIPLY:  Red(5),
		DIVIDE:    Red(5),
		REMAINDER: Red(5),
		DEGREE:    Sh(5),
		R_PAREN:   Red(5),
		EOF:       Red(5),
	},
}

// ((2 + 4) - 3) + 2 * (1)
