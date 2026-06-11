package reporter

import (
	"fmt"
	"fugu/pkg/lexer"
	"fugu/pkg/reporter"
	"fugu/pkg/token"
	"testing"
)

// этот тест не имеет смысл я просто хотел посмотреть как работает вывод ошибок )
func TestWorkReporter(t *testing.T) {
	input := []byte(`module main

fn main() {
	let a: mut string
}`)
	lex := lexer.New(input, "main.fg")
	var tks []token.Token

	for {
		if tk := lex.NextToken(); tk.Kind == token.EOF {
			tks = append(tks, tk)
			break
		} else {
			tks = append(tks, tk)
		}
	}

	tk := tks[1]

	fmt.Println(
		reporter.BoldCyan(
			string(tk.Literal(lex)),
		),
		tk.Kind.String(),
		tk.Start,
	)

	rp := reporter.New(lex, "main.fg")
	rp.SendTk(reporter.TestError, tk)
	rp.Close()
}
