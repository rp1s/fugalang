package lexer

import (
	"fugu/pkg/token"
	"testing"
)

func TestLexerComment(t *testing.T) {
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

func TestLexerOperator(t *testing.T) {
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
