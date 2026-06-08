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

func TestLexerOpRange(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		expectedKind token.TokenKind
	}{
		{
			name:         "Тест 1",
			input:        "..",
			expectedKind: token.OP_RANGE,
		},
		{
			name:         "Тест 2",
			input:        "..=",
			expectedKind: token.RANGE_INCL,
		},
		{
			name:         "Тест 3",
			input:        "..<",
			expectedKind: token.RANGE_HALF_OPEN,
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
	}
}
