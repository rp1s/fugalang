package token

type TokenKind uint16

const (
	_         TokenKind = iota
	ILLEGAL             // неизвестный токен
	COMMENT             // comment
	M_COMMENT           // /* comment */
	SPACING             // whitespace
	EOF

	// Группы
	GNUMBER
	GSTRING
	GARITHMETIC

	literals_start

	INTEGER    // 123
	IMAGINARY  // 123i
	FLOATING   // 12.3
	STRING     // "abc"
	T_STRING   // "hello ${x}"
	RAW_STRING // `abc`
	CHARACTER  // 'a'
	IDENTIFIER // myVar

	literals_end
	keyword_start

	// ключевые слова объявления модулей и использования
	MODULE
	USE

	// ключевые слова объявления
	FN
	LET
	MUT
	CONST
	TYPE
	ENUM
	STRUCT
	INTERFACE

	// значения
	NONE
	TRUE
	FALSE
	CHAN // канал

	// контрольные конструкции
	RETURN
	MATCH
	IF
	ELSE
	FOR
	RANGE
	CONTINUE
	BREAK

	// управление выполнением
	DEFER  // Отложенный вызов функции перед выходом из области видимости
	SELECT // Ожидание первого готового события из каналов/корутин

	// корутины
	ASYNC // Объявление асинхронной функции или блока
	AWAIT // Ожидание завершения асинхронной операции
	YIELD // Возврат промежуточного значения из генератора

	// нативные потоки
	SPAWN // Запуск функции в отдельном системном потоке ОС

	// групповое управление
	SYNC // Создание изолированной области для группы задач
	WAIT // Ожидание завершения ВСЕХ задач в текущем контексте
	HALT // Безопасная экстренная остановка ВСЕХ задач группы

	// низкоуровневые операции
	CLANG  // вставка C кода
	EXPORT // подключение внешней функции из динамической библиотеки
	EXTERN // объявление внешней функции без реализации (FFI)
	UNSAFE // блок небезопасного кода, требующий явного разрешения

	keyword_end
	operator_start

	// операторы группировки
	L_PAREN // (
	R_PAREN // )
	L_BRACE // {
	R_BRACE // }
	L_BRACK // [
	R_BRACK // ]

	// операторы присваевания
	APPROPRIATE  // :=
	REDEFINITION // =
	A_DECREASE   // -=
	A_INCREASE   // +=
	A_MULTIPLY   // *=
	A_DIVIDE     // /=
	A_REMAINDER  // %=
	A_DEGREE     // ^=

	// логические операторы сравнения
	LIKEN         // ==
	NOT_EQUAL     // !=
	LESS_EQUAL    // <=
	GREATER_EQUAL // >=
	LESS          // <
	GREATER       // >
	NOT           // !
	AND           // &&
	OR            // ||

	// операторы арифметики
	DECREASE  // -
	INCREASE  // +
	MULTIPLY  // *
	DIVIDE    // /
	REMAINDER // %
	DEGREE    // ^

	// операторы битовых сдвигов и побитовых операций
	SHR_LESS    // <<
	SHR_GREATER // >>
	BITWISE_NOT // ~

	// операторы диапазонов
	OP_RANGE        // ..   (Исключающий / Открытый)
	RANGE_INCL      // ..=  (Включающий / Закрытый)
	RANGE_HALF_OPEN // ..<  (Полуоткрытый)
	OP_ARRAY        // ...

	// операторы управления данных
	GOES_OVER //  =>
	OP_RETURN // ->
	PIPE      // |>
	DEFAULT   // ?:
	SAFE_DOT  // ?.
	TAKE_LINK // &

	// операторы разделения
	COLON // :
	END   // ;
	COMMA // ,
	DOT   // .

	operator_end

	EndToken
)

func (tk *TokenKind) Group() TokenKind {
	switch *tk {
	case INTEGER, IMAGINARY, FLOATING:
		return GNUMBER
	case STRING, T_STRING, RAW_STRING:
		return GSTRING
	case DECREASE, INCREASE, MULTIPLY, DIVIDE, REMAINDER, DEGREE:
		return GARITHMETIC
	default:
		return *tk
	}
}

func Group(tk TokenKind) TokenKind {
	return tk.Group()
}

type Token struct {
	Kind  TokenKind
	Pos   Position // начало токена
	Start int      // абсолютное смещение до начала токена
	End   int      // абсолютное смещение до конца токена
}

// структура указывающая позицию в коде токена
type Position struct {
	FileName string
	Line     int
	Column   int
	Offset   int
}

type OutLiteral interface {
	Input() *[]byte
}

func (tk Token) Literal(source OutLiteral) []byte {
	switch tk.Kind {
	case STRING:
		return (*source.Input())[tk.Start+1 : tk.End-1]
	case T_STRING:
		return (*source.Input())[tk.Start+1 : tk.End-1]
	case RAW_STRING:
		return (*source.Input())[tk.Start+1 : tk.End-1]
	case CHARACTER:
		return (*source.Input())[tk.Start+1 : tk.End-1]
	case COMMENT:
		return (*source.Input())[tk.Start+2 : tk.End]
	case M_COMMENT:
		return (*source.Input())[tk.Start+2 : tk.End-2]
	default:
		return (*source.Input())[tk.Start:tk.End]
	}
}

var keywords = map[string]TokenKind{
	// Объявление модулей и использование
	"module": MODULE,
	"use":    USE,

	// Объявления структур данных и переменных
	"fn":        FN,
	"let":       LET,
	"mut":       MUT,
	"const":     CONST,
	"type":      TYPE,
	"enum":      ENUM,
	"struct":    STRUCT,
	"interface": INTERFACE,

	// Встроенные значения
	"none":  NONE,
	"true":  TRUE,
	"false": FALSE,
	"chan":  CHAN,

	// Контрольные конструкции
	"return":   RETURN,
	"match":    MATCH,
	"if":       IF,
	"else":     ELSE,
	"for":      FOR,
	"range":    RANGE,
	"continue": CONTINUE,
	"break":    BREAK,

	// Управление выполнением
	"defer":  DEFER,
	"select": SELECT,

	// Асинхронность и корутины
	"async": ASYNC,
	"await": AWAIT,
	"yield": YIELD,

	// Системные потоки ОС
	"spawn": SPAWN,

	// Групповое управление задачами
	"sync": SYNC,
	"wait": WAIT,
	"halt": HALT,

	// Низкоуровневые операции и FFI
	"clang":  CLANG,
	"export": EXPORT,
	"extern": EXTERN,
	"unsafe": UNSAFE,
}

// проверяет является ли строка ключевым словом.
// Если да возвращает его тип, если нет возвращает простой IDENTIFIER.
func SearchKeyword(ident []byte) TokenKind {
	if kind, ok := keywords[string(ident)]; ok {
		return kind
	}
	return IDENTIFIER
}

func (tk TokenKind) String() string {
	switch tk {
	case ILLEGAL:
		return "ILLEGAL"
	case COMMENT:
		return "COMMENT"
	case SPACING:
		return "SPACING"
	case EOF:
		return "EOF"
	case INTEGER:
		return "INTEGER"
	case IMAGINARY:
		return "IMAGINARY"
	case FLOATING:
		return "FLOATING"
	case STRING:
		return "STRING"
	case T_STRING:
		return "T_STRING"
	case RAW_STRING:
		return "RAW_STRING"
	case CHARACTER:
		return "CHARACTER"
	case IDENTIFIER:
		return "IDENTIFIER"
	case MODULE:
		return "MODULE"
	case USE:
		return "USE"
	case FN:
		return "FN"
	case LET:
		return "LET"
	case MUT:
		return "MUT"
	case CONST:
		return "CONST"
	case TYPE:
		return "TYPE"
	case ENUM:
		return "ENUM"
	case STRUCT:
		return "STRUCT"
	case INTERFACE:
		return "INTERFACE"
	case NONE:
		return "NONE"
	case TRUE:
		return "TRUE"
	case FALSE:
		return "FALSE"
	case CHAN:
		return "CHAN"
	case RETURN:
		return "RETURN"
	case MATCH:
		return "MATCH"
	case IF:
		return "IF"
	case ELSE:
		return "ELSE"
	case FOR:
		return "FOR"
	case RANGE:
		return "RANGE"
	case CONTINUE:
		return "CONTINUE"
	case BREAK:
		return "BREAK"
	case DEFER:
		return "DEFER"
	case SELECT:
		return "SELECT"
	case ASYNC:
		return "ASYNC"
	case AWAIT:
		return "AWAIT"
	case YIELD:
		return "YIELD"
	case SPAWN:
		return "SPAWN"
	case SYNC:
		return "SYNC"
	case WAIT:
		return "WAIT"
	case HALT:
		return "HALT"
	case CLANG:
		return "CLANG"
	case EXPORT:
		return "EXPORT"
	case EXTERN:
		return "EXTERN"
	case UNSAFE:
		return "UNSAFE"
	case APPROPRIATE:
		return "APPROPRIATE"
	case REDEFINITION:
		return "REDEFINITION"
	case A_DECREASE:
		return "A_DECREASE"
	case A_INCREASE:
		return "A_INCREASE"
	case A_MULTIPLY:
		return "A_MULTIPLY"
	case A_DIVIDE:
		return "A_DIVIDE"
	case A_REMAINDER:
		return "A_REMAINDER"
	case LIKEN:
		return "LIKEN"
	case NOT_EQUAL:
		return "NOT_EQUAL"
	case LESS_EQUAL:
		return "LESS_EQUAL"
	case GREATER_EQUAL:
		return "GREATER_EQUAL"
	case LESS:
		return "LESS"
	case GREATER:
		return "GREATER"
	case NOT:
		return "NOT"
	case AND:
		return "AND"
	case OR:
		return "OR"
	case DECREASE:
		return "DECREASE"
	case INCREASE:
		return "INCREASE"
	case MULTIPLY:
		return "MULTIPLY"
	case DIVIDE:
		return "DIVIDE"
	case REMAINDER:
		return "REMAINDER"
	case DEGREE:
		return "DEGREE"
	case SHR_LESS:
		return "SHR_LESS"
	case SHR_GREATER:
		return "SHR_GREATER"
	case BITWISE_NOT:
		return "BITWISE_NOT"
	case TAKE_LINK:
		return "TAKE_LINK"
	case OP_RANGE:
		return "OP_RANGE"
	case RANGE_HALF_OPEN:
		return "RANGE_HALF_OPEN"
	case RANGE_INCL:
		return "RANGE_INCL"
	case GOES_OVER:
		return "GOES_OVER"
	case PIPE:
		return "PIPE"
	case DEFAULT:
		return "DEFAULT"
	case SAFE_DOT:
		return "SAFE_DOT"
	case L_PAREN:
		return "L_PAREN"
	case R_PAREN:
		return "R_PAREN"
	case L_BRACE:
		return "L_BRACE"
	case R_BRACE:
		return "R_BRACE"
	case L_BRACK:
		return "L_BRACK"
	case R_BRACK:
		return "R_BRACK"
	case COLON:
		return "COLON"
	case END:
		return "END"
	case COMMA:
		return "COMMA"
	case DOT:
		return "DOT"
	}
	return ""
}
