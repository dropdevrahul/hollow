package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

type Opcode int

const (
	OP_PUSH = iota
	OP_PLUS
	OP_DUMP
	OP_COUNT
)

var Stack []int

type Op struct {
	N Opcode // name
	O int    // operads  currently only have 1 for push
}

func Push(x int) Op {
	return Op{
		N: OP_PUSH,
		O: x,
	}
}

func Plus() Op {
	return Op{
		N: OP_PLUS,
	}
}

func Dump() Op {
	return Op{
		N: OP_DUMP,
	}
}

func Compile(program []Op) {
	b, err := os.ReadFile("dump.hollow")
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Create("output.asm")
	if err != nil {
		log.Fatal(err)
	}

	f.WriteString("%define SYS_EXIT 60\n\n")
	f.WriteString("section .text\n\n")
	f.Write(b)
	f.WriteString("global _start\n")
	f.WriteString("_start:\n")

	for _, op := range program {
		switch op.N {
		case OP_PUSH:
			f.WriteString(fmt.Sprintf("    push %d\n", op.O))
		case OP_DUMP:
			f.WriteString("    pop rdi\n")
			f.WriteString("    call dump\n")
		}
	}

	f.WriteString("    mov rax, SYS_EXIT\n")
	f.WriteString("    mov rdi, 0\n")
	f.WriteString("    syscall\n")
	f.Close()

	o, err := exec.Command("nasm", "-felf64", "output.asm").Output()
	if err != nil {
		log.Fatal(o, err)
	}

	o, err = exec.Command("ld", "-o", "output", "output.o").Output()
	if err != nil {
		log.Fatal(o, err)
	}

}

func main() {
	log.SetOutput(os.Stdout)
	// Our language consists of a runtime stack we can modify
	// using various instructions
	program := []Op{
		Push(64),
		Dump(),
		Push(78778),
		Dump(),
	}

	Compile(program)
}
