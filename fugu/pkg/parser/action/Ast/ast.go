package ast

type Arena struct {
	nodes   []Node
	strings []string
}

type StringID uint32
type NodeID uint32

type Node struct {
	Kind  NodeKind
	Flags uint16 // флаги: public, const, mut

	Data1 uint32
	Data2 uint32
	Data3 uint32
	Extra int64 // числовые значения
}

type NodeKind uint16

const (
	KindInvalid NodeKind = iota
)
