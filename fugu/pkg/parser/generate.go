//go:build ignore

package main

import (
	"fmt"
	"fugu/pkg/parser"
	"fugu/pkg/token"
	"os"
)

func main() {
	fnInit := parser.GenerateInitSlice(parser.ActionSrc, int(token.EndToken))
	content := fmt.Sprintf(`package parser

%s
`, fnInit)

	os.WriteFile("action_gen.go", []byte(content), 0644)
}
