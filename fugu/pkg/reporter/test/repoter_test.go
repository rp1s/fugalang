package reporter

import (
	"fugu/pkg/lexer"
	"fugu/pkg/reporter"
	"fugu/pkg/token"
	"testing"
)

// этот тест не имеет смысл я просто хотел посмотреть как работает вывод ошибок )
func TestWorkReporter(t *testing.T) {
	input := "let a: mut string"
	lex := lexer.New(input, "main.fg")

	tk := lex.NextToken()
	println(tk.Kind.String())

	if tk.Kind != token.EOF {
		for {
			tk := lex.NextToken()

			println(tk.Kind.String())
			if tk.Kind == token.EOF {
				break
			}
		}
	}
	rp := reporter.New(lex, "main.fg")
	rp.SendTk(reporter.TestError, tk)
}
