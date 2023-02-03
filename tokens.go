package main

type Opcode int

const (
	OP_PUSH = iota
	OP_PLUS
	OP_SUB
	OP_DUMP
	OP_VAR
	OP_EQUALS
	OP_IF
	OP_END
	OP_COUNT
)

const (
	S_OP_PLUS   = "+"
	S_OP_SUB    = "-"
	S_OP_DUMP   = "."
	S_OP_EQUALS = "=="
	S_OP_IF     = "if"
	S_OP_END    = "end"
)

var Stack []int

type Op struct {
	N Opcode // name
	O int    // operand  currently only have 1 for push
	V string // variable name
}

type Program []Op
