package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func Compile(program []Op, output string) {
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

	for i, op := range program {
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
		case OP_MUL:
			f.WriteString(fmt.Sprintf("    push %d\n", op.O))
			f.WriteString("    pop rcx\n")
			f.WriteString("    pop rax\n")
			f.WriteString("    mul rcx\n")
			f.WriteString("    xor rcx, rcx\n") // contains higher bits
			f.WriteString("    push rax\n")     // contains lower bits
		case OP_DIV:
			f.WriteString(fmt.Sprintf("    push %d\n", op.O))
			f.WriteString("    xor rdx, rdx\n")
			f.WriteString("    pop rcx\n") // first we pop divisor  dividend/divisor
			f.WriteString("    pop rax\n") // div takes two register as dividend rdx:rax
			f.WriteString("    div rcx\n")
			f.WriteString("    push rax\n") // rax has quotient
		case OP_MOD:
			f.WriteString(fmt.Sprintf("    push %d\n", op.O))
			f.WriteString("    xor rdx, rdx\n")
			f.WriteString("    pop rcx\n") // first we pop divisor  dividend/divisor
			f.WriteString("    pop rax\n") // div takes two register as dividend rdx:rax
			f.WriteString("    div rcx\n")
			f.WriteString("    push rcx\n") //rdx has remainder
		case OP_DUMP:
			f.WriteString("    pop rdi\n")
			f.WriteString("    call dump\n")
		case OP_IF:
			f.WriteString("    pop rax\n")
			// set zf = 1 if rax is 0
			f.WriteString("    test rax, rax\n")
			// jump to jmp address block if zf == 1 i.e condition is false
			f.WriteString(fmt.Sprintf("    jz addr_%d\n", op.JMP))
		case OP_ELSE:
			f.WriteString(fmt.Sprintf("    jmp addr_%d\n", op.JMP))
		case OP_END:
			f.WriteString(fmt.Sprintf("addr_%d:\n", i))
		case OP_EQUALS, OP_GTE, OP_LTE, OP_LT, OP_GT:
			if op.INS == "" {
				log.Fatalf("No ins for Comparison operator %d", op.N)
			}
			f.WriteString(fmt.Sprintf("    push %d\n", op.O))
			f.WriteString("    mov rcx, 0\n")
			f.WriteString("    mov rdx, 1\n")
			f.WriteString("    pop rax\n")
			f.WriteString("    pop rbx\n")
			f.WriteString("    cmp rax, rbx\n")
			f.WriteString(fmt.Sprintf("    %s rcx, rdx\n", op.INS))
			f.WriteString("    push rcx\n")
		}
	}

	f.WriteString("    mov rax, SYS_EXIT\n")
	f.WriteString("    mov rdi, 0\n")
	f.WriteString("    syscall\n")

	// add section .data with variables
	f.WriteString("section .data\n")

	//for _, v := range vars {
	//	f.WriteString(fmt.Sprintf("  %s: dq  %d\n", v.V, v.O))
	//}

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
