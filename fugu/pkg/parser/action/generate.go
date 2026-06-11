//go:build ignore

package main

import (
	"fugu/pkg/parser/action"
	"fugu/pkg/token"
	"os"
)

func main() {
	content := action.GenerateActionTable(action.ActionSrc, int(token.EndToken))
	os.WriteFile("action_gen.go", []byte(content), 0644)
}
