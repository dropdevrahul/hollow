package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
)

type Opcode int

const (
	OP_PUSH = iota
	OP_PLUS
	OP_SUB
	OP_DUMP
	OP_VAR
	OP_EQUALS
	OP_COUNT
)

const (
	S_OP_PLUS   = "+"
	S_OP_SUB    = "-"
	S_OP_DUMP   = "."
	S_OP_EQUALS = "=="
)

var Stack []int

type Op struct {
	N Opcode // name
	O int    // operand  currently only have 1 for push
	V string // variable name
}

type Program []Op

func Push(x int) Op {
	return Op{
		N: OP_PUSH,
		O: x,
	}
}

func Plus(x int) Op {
	return Op{
		N: OP_PLUS,
		O: x,
	}
}

func Sub(x int) Op {
	return Op{
		N: OP_SUB,
		O: x,
	}
}

func Dump() Op {
	return Op{
		N: OP_DUMP,
	}
}

func Equals(x int) Op {
	return Op{
		N: OP_EQUALS,
	}
}

func Var(v string, o int) Op {
	return Op{
		N: OP_VAR,
		O: o,
		V: v,
	}
}

func Tokenize(fpath string) (Program, error) {
	file, err := os.Open(fpath)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	split := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		advance, token, err = bufio.ScanWords(data, atEOF)
		return
	}
	// Set the split function for the scanning operation.
	scanner.Split(split)
	p := Program{}
	for scanner.Scan() {
		w := scanner.Text()
		switch w {
		case S_OP_PLUS:
			scanner.Scan()
			w := scanner.Text()
			t, err := strconv.Atoi(w)
			if err != nil {
				return p, err
			}
			o := Plus(t)
			p = append(p, o)
		case S_OP_SUB:
			scanner.Scan()
			w := scanner.Text()
			t, err := strconv.Atoi(w)
			if err != nil {
				return p, err
			}
			o := Sub(t)
			p = append(p, o)
		case S_OP_EQUALS:
			// get operand
			scanner.Scan()
			w := scanner.Text()
			t, err := strconv.Atoi(w)
			if err != nil {
				return p, err
			}
			o := Equals(t)
			p = append(p, o)
		case S_OP_DUMP:
			o := Dump()
			p = append(p, o)
		default:
			t, err := strconv.Atoi(w)
			if err != nil {
				return p, err
			}
			o := Push(t)
			p = append(p, o)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("Invalid input: %s", err)
	}

	return p, nil
}

func Compile(program []Op, output string) {
	var vars []Op
	b, err := os.ReadFile("hbin/dump.hbin")
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Create(output + ".asm")
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
		case OP_PLUS:
			f.WriteString(fmt.Sprintf("    push %d\n", op.O))
			f.WriteString("    pop rax\n")
			f.WriteString("    pop rbx\n")
			f.WriteString("    add rax, rbx\n")
			f.WriteString("    push rax\n")
		case OP_SUB:
			f.WriteString(fmt.Sprintf("    push %d\n", op.O))
			f.WriteString("    pop rax\n")
			f.WriteString("    pop rbx\n")
			f.WriteString("    sub rbx, rax\n")
			f.WriteString("    push rbx\n")
		case OP_DUMP:
			f.WriteString("    pop rdi\n")
			f.WriteString("    call dump\n")
		case OP_VAR:
			vars = append(vars, op)
		case OP_EQUALS:
			f.WriteString("    mov rcx, 0\n")
			f.WriteString("    mov rdx, 1\n")
			f.WriteString("    pop rax\n")
			f.WriteString("    pop rbx\n")
			f.WriteString("    cmp rax, rbx\n")
			f.WriteString("    cmovae rcx, rdx\n")
			f.WriteString("    push rcx\n")
		}
	}

	f.WriteString("    mov rax, SYS_EXIT\n")
	f.WriteString("    mov rdi, 0\n")
	f.WriteString("    syscall\n")

	// add section .data with variables
	f.WriteString("section .data\n")

	for _, v := range vars {
		f.WriteString(fmt.Sprintf("  %s: dq  %d\n", v.V, v.O))
	}

	f.Close()

	o, err := exec.Command("nasm", "-felf64", output+".asm").Output()
	if err != nil {
		log.Fatal(o, err)
	}

	o, err = exec.Command("ld", "-o", output, output+".o").Output()
	if err != nil {
		log.Fatal(o, err)
	}

}

func main() {
	log.SetOutput(os.Stdout)
	// Our language consists of a runtime stack we can modify
	// using various instructions
	output := flag.String("o", "output", "output file name")

	flag.Parse()

	program, err := Tokenize(flag.Args()[0])
	if err != nil {
		log.Println(err)
	}
	Compile(program, *output)
}
