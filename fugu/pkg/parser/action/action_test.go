package action

import (
	"fugu/pkg/token"
	"testing"

	"github.com/k0kubun/pp/v3"
)

func TestActionMap(t *testing.T) {
	pp.Println(Action(1, token.GARITHMETIC))
}
