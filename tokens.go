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
	OP_ELSE
	OP_COUNT
)

const (
	S_OP_PLUS   = "+"
	S_OP_SUB    = "-"
	S_OP_DUMP   = "."
	S_OP_EQUALS = "=="
	S_OP_IF     = "if"
	S_OP_END    = "end"
	S_OP_ELSE   = "else"
)

var Stack []int

type Op struct {
	N   Opcode // name
	O   int    // operand for the operation
	JMP int    // jmp address to be used by some ops like If/else/while etc
}

type Program []Op
