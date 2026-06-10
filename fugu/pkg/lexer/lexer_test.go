package lexer

import (
	"fugu/pkg/token"
	"testing"
)

func TestComment(t *testing.T) {
	// 1. Описываем структуру тест-кейса
	tests := []struct {
		name            string
		input           []byte
		checkSpaceTk    bool // проверка токена пробел следущим токеном
		expectedKind    token.TokenKind
		expectedLiteral []byte
	}{
		{
			name:            "Тест обработки однострочный комментарий",
			input:           []byte("// привет slava"),
			expectedKind:    token.COMMENT,
			expectedLiteral: []byte(" привет slava"),
		},
		{
			name: "Тест обработки многострочного коментария",
			input: []byte(`/* 
Привет, это многострочный коммент для проверки корректной работы лексера) 
Пока расскажу вам про моего замечательного кота Фантика (полное имя Элефант). 
Он был породистым мейнкуном, я его очень люблю, но по состоянию здоровья 
моей семье пришлось отдать его знакомым :( 
*/ 
`),
			checkSpaceTk: true,
			expectedKind: token.M_COMMENT,
			expectedLiteral: []byte(` 
Привет, это многострочный коммент для проверки корректной работы лексера) 
Пока расскажу вам про моего замечательного кота Фантика (полное имя Элефант). 
Он был породистым мейнкуном, я его очень люблю, но по состоянию здоровья 
моей семье пришлось отдать его знакомым :( 
`),
		},
	}

	for _, tt := range tests {
		lex := New(tt.input, "main.fg")
		tk := lex.NextToken()

		if tk.Kind != tt.expectedKind {
			t.Errorf("[%s] Неверный тип токена. Ожидался: %s, получен: %s",
				tt.name, tt.expectedKind.String(), tk.Kind.String())
			continue // переход к след тесту
		}

		lit := tk.Literal(lex)
		if string(lit) != string(tt.expectedLiteral) {
			t.Errorf("[%s] Неверный литерал.\nОжидался:\n%q\n\nПолучен:\n%q",
				tt.name, tt.expectedLiteral, lit)
		}

		if tt.checkSpaceTk {
			tk = lex.NextToken()
			if tk.Kind != token.SPACING {
				t.Errorf("[%s] Неверный тип токена. Ожидался: %s, получен: %s",
					tt.name, token.SPACING.String(), tk.Kind.String())
			}
		}

		tk = lex.NextToken()
		if tk.Kind != token.EOF {
			t.Errorf("[%s] Неверный тип токена. Ожидался: %s, получен: %s",
				tt.name, token.EOF.String(), tk.Kind.String())
		}

	}
}

func TestOperator(t *testing.T) {
	tests := []struct {
		name         string
		input        []byte
		expectedKind token.TokenKind
	}{
		// Операторы диапазонов
		{name: "Исключающий диапазон", input: []byte(".."), expectedKind: token.OP_RANGE},
		{name: "Включающий диапазон", input: []byte("..="), expectedKind: token.RANGE_INCL},
		{name: "Полуоткрытый диапазон", input: []byte("..<"), expectedKind: token.RANGE_HALF_OPEN},

		// Операторы присваивания
		{name: "Присваивание с объявлением", input: []byte(":="), expectedKind: token.APPROPRIATE},
		{name: "Обычное переопределение", input: []byte("="), expectedKind: token.REDEFINITION},
		{name: "Уменьшение с присваиванием", input: []byte("-="), expectedKind: token.A_DECREASE},
		{name: "Увеличение с присваиванием", input: []byte("+="), expectedKind: token.A_INCREASE},
		{name: "Умножение с присваиванием", input: []byte("*="), expectedKind: token.A_MULTIPLY},
		{name: "Деление с присваиванием", input: []byte("/="), expectedKind: token.A_DIVIDE},
		{name: "Остаток с присваиванием", input: []byte("%="), expectedKind: token.A_REMAINDER},
		{name: "Возведение в степень с присваиванием", input: []byte("^="), expectedKind: token.A_DEGREE},

		// Логические операторы сравнения
		{name: "Равенство", input: []byte("=="), expectedKind: token.LIKEN},
		{name: "Неравенство", input: []byte("!="), expectedKind: token.NOT_EQUAL},
		{name: "Меньше или равно", input: []byte("<="), expectedKind: token.LESS_EQUAL},
		{name: "Больше или равно", input: []byte(">="), expectedKind: token.GREATER_EQUAL},
		{name: "Меньше", input: []byte("<"), expectedKind: token.LESS},
		{name: "Больше", input: []byte(">"), expectedKind: token.GREATER},
		{name: "Логическое НЕ", input: []byte("!"), expectedKind: token.NOT},
		{name: "Логическое И", input: []byte("&&"), expectedKind: token.AND},
		{name: "Логическое ИЛИ", input: []byte("||"), expectedKind: token.OR},

		// Операторы арифметики
		{name: "Минус", input: []byte("-"), expectedKind: token.DECREASE},
		{name: "Плюс", input: []byte("+"), expectedKind: token.INCREASE},
		{name: "Умножение", input: []byte("*"), expectedKind: token.MULTIPLY},
		{name: "Деление", input: []byte("/"), expectedKind: token.DIVIDE},
		{name: "Остаток от деления", input: []byte("%"), expectedKind: token.REMAINDER},
		{name: "Степень", input: []byte("^"), expectedKind: token.DEGREE},

		// Побитовые операторы
		{name: "Битовый сдвиг влево", input: []byte("<<"), expectedKind: token.SHR_LESS},
		{name: "Битовый сдвиг вправо", input: []byte(">>"), expectedKind: token.SHR_GREATER},
		{name: "Побитовое НЕ", input: []byte("~"), expectedKind: token.BITWISE_NOT},

		// Операторы управления данными
		{name: "Лямбда", input: []byte("=>"), expectedKind: token.GOES_OVER},
		{name: "Пайплайн", input: []byte("|>"), expectedKind: token.PIPE},
		{name: "Тернарный оператор", input: []byte("?:"), expectedKind: token.DEFAULT},
		{name: "Безопасный вызов", input: []byte("?."), expectedKind: token.SAFE_DOT},
		{name: "Взятие ссылки", input: []byte("&"), expectedKind: token.TAKE_LINK},

		// Операторы группировки
		{name: "Левая круглая скобка: ( ", input: []byte("("), expectedKind: token.L_PAREN},
		{name: "Правая круглая скобка: ) ", input: []byte(")"), expectedKind: token.R_PAREN},
		{name: "Левая фигурная скобка: { ", input: []byte("{"), expectedKind: token.L_BRACE},
		{name: "Правая фигурная скобка: } ", input: []byte("}"), expectedKind: token.R_BRACE},
		{name: "Левая квадратная скобка: [ ", input: []byte("["), expectedKind: token.L_BRACK},
		{name: "Правая квадратная скобка: ] ", input: []byte("]"), expectedKind: token.R_BRACK},

		// Операторы разделения
		{name: "Двоеточие", input: []byte(":"), expectedKind: token.COLON},
		{name: "Точка с запятой", input: []byte(";"), expectedKind: token.END},
		{name: "Запятая", input: []byte(","), expectedKind: token.COMMA},
		{name: "Точка", input: []byte("."), expectedKind: token.DOT},
	}

	for _, tt := range tests {
		lex := New(tt.input, "main.fg")
		tk := lex.NextToken()

		if tk.Kind != tt.expectedKind {
			t.Errorf("[%s] Неверный тип токена. Ожидался: %s, получен: %s",
				tt.name, tt.expectedKind.String(), tk.Kind.String())
			continue
		}
	}
}

func TestLiteral(t *testing.T) {
	tests := []struct {
		name            string
		input           []byte
		expectedKind    token.TokenKind
		expectedLiteral []byte
	}{
		{name: "Идентификатор начинающийся числом", input: []byte("10pixel"), expectedKind: token.IDENTIFIER, expectedLiteral: []byte("10pixel")},
		{name: "Идентификатор с цифрами и подчёркиванием", input: []byte("2stack"), expectedKind: token.IDENTIFIER, expectedLiteral: []byte("2stack")},
		{name: "Обычный идентификатор с подчёркивания", input: []byte("__init__"), expectedKind: token.IDENTIFIER, expectedLiteral: []byte("__init__")},

		{name: "Ключевое слово module", input: []byte("module"), expectedKind: token.MODULE, expectedLiteral: []byte("module")},
		{name: "Ключевое слово fn", input: []byte("fn"), expectedKind: token.FN, expectedLiteral: []byte("fn")},

		{name: "Целое число (INTEGER)", input: []byte("12443"), expectedKind: token.INTEGER, expectedLiteral: []byte("12443")},
		{name: "Дробное число (FLOATING)", input: []byte("12.3"), expectedKind: token.FLOATING, expectedLiteral: []byte("12.3")},
		{name: "Мнимое целое число (IMAGINARY)", input: []byte("123i"), expectedKind: token.IMAGINARY, expectedLiteral: []byte("123i")},
		{name: "Мнимое дробное число (IMAGINARY)", input: []byte("12.3i"), expectedKind: token.IMAGINARY, expectedLiteral: []byte("12.3i")},

		{name: "Обычная строка (STRING)", input: []byte(`"hello world"`), expectedKind: token.STRING, expectedLiteral: []byte("hello world")},
		{name: "Сырая строка (RAW_STRING)", input: []byte("`multiline code`"), expectedKind: token.RAW_STRING, expectedLiteral: []byte("multiline code")},
		{name: "Одиночный символ (CHARACTER)", input: []byte("'я'"), expectedKind: token.CHARACTER, expectedLiteral: []byte("я")},
		{name: "Экранированный символ (CHARACTER)", input: []byte("'\\n'"), expectedKind: token.CHARACTER, expectedLiteral: []byte("\\n")},
		{name: "Строка с интерполяцией (T_STRING)", input: []byte(`"status: ${code}"`), expectedKind: token.T_STRING, expectedLiteral: []byte("status: ${code}")},
	}

	for _, tt := range tests {
		lex := New(tt.input, "main.fg")
		tk := lex.NextToken()

		if tk.Kind != tt.expectedKind {
			t.Errorf("[%s] Неверный тип токена. Ожидался: %s, получен: %s",
				tt.name, tt.expectedKind.String(), tk.Kind.String())
			continue
		}

		if string(tk.Literal(lex)) != string(tt.expectedLiteral) {
			t.Errorf("[%s] Неверный литерал. Ожидался: %q, получен: %q",
				tt.name, tt.expectedLiteral, tk.Literal(lex))
			continue
		}
	}
}

func TestLexerStabilization(t *testing.T) {
	tests := []struct {
		name         string
		input        []byte
		expectedKind token.TokenKind
	}{
		{
			name:         "cтабилизация после незакрытого многострочного комментария",
			input:        []byte("/* незакрытый комментарий \n fn main() {}"),
			expectedKind: token.FN,
		},
		{
			name: "cтабилизация после незакрытой обычной строки",
			input: []byte(`"незакрытая строка
if x == 5 {}`),
			expectedKind: token.IF,
		},
		{
			name:         "cтабилизация после незакрытой сырой строки",
			input:        []byte("`незакрытая сырая строка \n else { return }"),
			expectedKind: token.ELSE,
		},
	}

	for _, tt := range tests {
		l := New(tt.input, "main.fg")

		firstTok := l.NextToken()
		if firstTok.Kind != token.ILLEGAL {
			t.Fatalf("[%s] Первый токен обязан быть ILLEGAL. Получен: %s", tt.name, firstTok.Kind.String())
		}

		secondTok := l.NextToken()
		if secondTok.Kind != tt.expectedKind {
			t.Errorf("[%s] Неверный тип токена после стабилизации. Ожидался: %s, получен: %s",
				tt.name, tt.expectedKind.String(), secondTok.Kind.String())
		}

	}
}
