package lexer

import (
	"fmt"
	"fugu/pkg/token"
	"testing"
)

func TestLexerComment(t *testing.T) {
	lex := New("// привет slava", "main.fg")
	tk := lex.NextToken()
	if tk.Kind == token.COMMENT {
		fmt.Println(tk.Kind.String())
		if lex.LiteralToken(tk) == "// привет slava" {
			fmt.Println(lex.LiteralToken(tk))
		} else {
			t.Fatal("Надо литерал '// привет slava' а получилось: ", lex.LiteralToken(tk))
		}
		return
	}
	t.Fatal("Надо токен комент а получилось: ", tk.Kind.String())
}
