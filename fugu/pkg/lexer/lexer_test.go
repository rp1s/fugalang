package lexer

import (
	"fugu/pkg/token"
	"testing"
)

func TestComment(t *testing.T) {
	// 1. Описываем структуру тест-кейса
	tests := []struct {
		name            string
		input           string
		checkSpaceTk    bool // проверка токена пробел следущим токеном
		expectedKind    token.TokenKind
		expectedLiteral string
	}{
		{
			name:            "Тест обработки однострочный комментарий",
			input:           "// привет slava",
			expectedKind:    token.COMMENT,
			expectedLiteral: " привет slava",
		},
		{
			name: "Тест обработки многострочного коментария",
			input: `/* 
Привет, это многострочный коммент для проверки корректной работы лексера) 
Пока расскажу вам про моего замечательного кота Фантика (полное имя Элефант). 
Он был породистым мейнкуном, я его очень люблю, но по состоянию здоровья 
моей семье пришлось отдать его знакомым :( 
*/ 
`,
			checkSpaceTk: true,
			expectedKind: token.M_COMMENT,
			expectedLiteral: ` 
Привет, это многострочный коммент для проверки корректной работы лексера) 
Пока расскажу вам про моего замечательного кота Фантика (полное имя Элефант). 
Он был породистым мейнкуном, я его очень люблю, но по состоянию здоровья 
моей семье пришлось отдать его знакомым :( 
`,
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
		if lit != tt.expectedLiteral {
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
		input        string
		expectedKind token.TokenKind
	}{
		// Операторы диапазонов
		{name: "Исключающий диапазон", input: "..", expectedKind: token.OP_RANGE},
		{name: "Включающий диапазон", input: "..=", expectedKind: token.RANGE_INCL},
		{name: "Полуоткрытый диапазон", input: "..<", expectedKind: token.RANGE_HALF_OPEN},

		// Операторы присваивания
		{name: "Присваивание с объявлением", input: ":=", expectedKind: token.APPROPRIATE},
		{name: "Обычное переопределение", input: "=", expectedKind: token.REDEFINITION},
		{name: "Уменьшение с присваиванием", input: "-=", expectedKind: token.A_DECREASE},
		{name: "Увеличение с присваиванием", input: "+=", expectedKind: token.A_INCREASE},
		{name: "Умножение с присваиванием", input: "*=", expectedKind: token.A_MULTIPLY},
		{name: "Деление с присваиванием", input: "/=", expectedKind: token.A_DIVIDE},
		{name: "Остаток с присваиванием", input: "%=", expectedKind: token.A_REMAINDER},
		{name: "Возведение в степень с присваиванием", input: "^=", expectedKind: token.A_DEGREE},

		// Логические операторы сравнения
		{name: "Равенство", input: "==", expectedKind: token.LIKEN},
		{name: "Неравенство", input: "!=", expectedKind: token.NOT_EQUAL},
		{name: "Меньше или равно", input: "<=", expectedKind: token.LESS_EQUAL},
		{name: "Больше или равно", input: ">=", expectedKind: token.GREATER_EQUAL},
		{name: "Меньше", input: "<", expectedKind: token.LESS},
		{name: "Больше", input: ">", expectedKind: token.GREATER},
		{name: "Логическое НЕ", input: "!", expectedKind: token.NOT},
		{name: "Логическое И", input: "&&", expectedKind: token.AND},
		{name: "Логическое ИЛИ", input: "||", expectedKind: token.OR},

		// Операторы арифметики
		{name: "Минус", input: "-", expectedKind: token.DECREASE},
		{name: "Плюс", input: "+", expectedKind: token.INCREASE},
		{name: "Умножение", input: "*", expectedKind: token.MULTIPLY},
		{name: "Деление", input: "/", expectedKind: token.DIVIDE},
		{name: "Остаток от деления", input: "%", expectedKind: token.REMAINDER},
		{name: "Степень", input: "^", expectedKind: token.DEGREE},

		// Побитовые операторы
		{name: "Битовый сдвиг влево", input: "<<", expectedKind: token.SHR_LESS},
		{name: "Битовый сдвиг вправо", input: ">>", expectedKind: token.SHR_GREATER},
		{name: "Побитовое НЕ", input: "~", expectedKind: token.BITWISE_NOT},

		// Операторы управления данными
		{name: "Лямбда", input: "=>", expectedKind: token.GOES_OVER},
		{name: "Пайплайн", input: "|>", expectedKind: token.PIPE},
		{name: "Тернарный оператор", input: "?:", expectedKind: token.DEFAULT},
		{name: "Безопасный вызов", input: "?.", expectedKind: token.SAFE_DOT},
		{name: "Взятие ссылки", input: "&", expectedKind: token.TAKE_LINK},

		// Операторы группировки
		{name: "Левая круглая скобка: ( ", input: "(", expectedKind: token.L_PAREN},
		{name: "Правая круглая скобка: ) ", input: ")", expectedKind: token.R_PAREN},
		{name: "Левая фигурная скобка: { ", input: "{", expectedKind: token.L_BRACE},
		{name: "Правая фигурная скобка: } ", input: "}", expectedKind: token.R_BRACE},
		{name: "Левая квадратная скобка: [ ", input: "[", expectedKind: token.L_BRACK},
		{name: "Правая квадратная скобка: ] ", input: "]", expectedKind: token.R_BRACK},

		// Операторы разделения
		{name: "Двоеточие", input: ":", expectedKind: token.COLON},
		{name: "Точка с запятой", input: ";", expectedKind: token.END},
		{name: "Запятая", input: ",", expectedKind: token.COMMA},
		{name: "Точка", input: ".", expectedKind: token.DOT},
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
		input           string
		expectedKind    token.TokenKind
		expectedLiteral string
	}{
		{name: "Идентификатор начинающийся числом", input: "10pixel", expectedKind: token.IDENTIFIER, expectedLiteral: "10pixel"},
		{name: "Идентификатор с цифрами и подчёркиванием", input: "2stack", expectedKind: token.IDENTIFIER, expectedLiteral: "2stack"},
		{name: "Обычный идентификатор с подчёркивания", input: "__init__", expectedKind: token.IDENTIFIER, expectedLiteral: "__init__"},

		{name: "Ключевое слово module", input: "module", expectedKind: token.MODULE, expectedLiteral: "module"},
		{name: "Ключевое слово fn", input: "fn", expectedKind: token.FN, expectedLiteral: "fn"},

		{name: "Целое число (INTEGER)", input: "12443", expectedKind: token.INTEGER, expectedLiteral: "12443"},
		{name: "Дробное число (FLOATING)", input: "12.3", expectedKind: token.FLOATING, expectedLiteral: "12.3"},
		{name: "Мнимое целое число (IMAGINARY)", input: "123i", expectedKind: token.IMAGINARY, expectedLiteral: "123i"},
		{name: "Мнимое дробное число (IMAGINARY)", input: "12.3i", expectedKind: token.IMAGINARY, expectedLiteral: "12.3i"},

		{name: "Обычная строка (STRING)", input: `"hello world"`, expectedKind: token.STRING, expectedLiteral: "hello world"},
		{name: "Сырая строка (RAW_STRING)", input: "`multiline code`", expectedKind: token.RAW_STRING, expectedLiteral: "multiline code"},
		{name: "Одиночный символ (CHARACTER)", input: "'я'", expectedKind: token.CHARACTER, expectedLiteral: "я"},
		{name: "Экранированный символ (CHARACTER)", input: "'\\n'", expectedKind: token.CHARACTER, expectedLiteral: "\\n"},
		{name: "Строка с интерполяцией (T_STRING)", input: `"status: ${code}"`, expectedKind: token.T_STRING, expectedLiteral: "status: ${code}"},
	}

	for _, tt := range tests {
		lex := New(tt.input, "main.fg")
		tk := lex.NextToken()

		if tk.Kind != tt.expectedKind {
			t.Errorf("[%s] Неверный тип токена. Ожидался: %s, получен: %s",
				tt.name, tt.expectedKind.String(), tk.Kind.String())
			continue
		}

		if tk.Literal(lex) != tt.expectedLiteral {
			t.Errorf("[%s] Неверный литерал. Ожидался: %q, получен: %q",
				tt.name, tt.expectedLiteral, tk.Literal(lex))
			continue
		}
	}
}

func TestLexerStabilization(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		expectedKind token.TokenKind
	}{
		{
			name:         "cтабилизация после незакрытого многострочного комментария",
			input:        "/* незакрытый комментарий \n fn main() {}",
			expectedKind: token.FN,
		},
		{
			name: "cтабилизация после незакрытой обычной строки",
			input: `"незакрытая строка
if x == 5 {}`,
			expectedKind: token.IF,
		},
		{
			name:         "cтабилизация после незакрытой сырой строки",
			input:        "`незакрытая сырая строка \n else { return }",
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
