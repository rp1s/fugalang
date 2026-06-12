package action

import (
	. "fugu/pkg/token"
)

var ActionSrc = map[int]map[TokenKind]ActionStruct{
	// State 0 - стартовое состояние
	0: {
		GLITERAL: Sh(1),
		L_PAREN:  Sh(2),
	},

	// Работа с выражениями

	1: {
		GARITHMETIC: Red(1),
		R_PAREN:     Red(1),
		EOF:         Red(1),
	},
	2: {
		GLITERAL: Sh(1),
		L_PAREN:  Sh(2),
	},
	7: {
		GARITHMETIC: Red(2),
		R_PAREN:     Red(2),
		EOF:         Red(2),
	},
	8: {
		INCREASE:  Red(3),
		DECREASE:  Red(3),
		MULTIPLY:  Sh(5),
		DIVIDE:    Sh(6),
		REMAINDER: Sh(13),
		DEGREE:    Sh(14),
		R_PAREN:   Red(3),
		EOF:       Red(3),
	},
	9: {
		INCREASE:  Red(4),
		DECREASE:  Red(4),
		MULTIPLY:  Sh(5),
		DIVIDE:    Sh(6),
		REMAINDER: Sh(13),
		DEGREE:    Sh(14),
		R_PAREN:   Red(4),
		EOF:       Red(4),
	},
	10: {
		INCREASE:  Red(5),
		DECREASE:  Red(5),
		MULTIPLY:  Red(5),
		DIVIDE:    Red(5),
		REMAINDER: Red(5),
		DEGREE:    Sh(14),
		R_PAREN:   Red(5),
		EOF:       Red(5),
	},
	11: {
		INCREASE:  Red(6),
		DECREASE:  Red(6),
		MULTIPLY:  Red(6),
		DIVIDE:    Red(6),
		REMAINDER: Red(6),
		DEGREE:    Sh(14),
		R_PAREN:   Red(6),
		EOF:       Red(6),
	},
	12: {
		EOF: Acc(),
	},
	13: {
		GLITERAL: Sh(1),
		L_PAREN:  Sh(2),
	},
	14: {
		GLITERAL: Sh(1),
		L_PAREN:  Sh(2),
	},
	15: {
		INCREASE:  Red(7),
		DECREASE:  Red(7),
		MULTIPLY:  Red(7),
		DIVIDE:    Red(7),
		REMAINDER: Red(7),
		DEGREE:    Sh(14),
		R_PAREN:   Red(7),
		EOF:       Red(7),
	},
}

// ((2 + 4) - 3) + 2 * (1)
