package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

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

func Div(x int) Op {
	return Op{
		N: OP_DIV,
		O: x,
	}
}

func Mul(x int) Op {
	return Op{
		N: OP_MUL,
		O: x,
	}
}

func Mod(x int) Op {
	return Op{
		N: OP_MOD,
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
		O: x,
	}
}

func OpIf() Op {
	return Op{
		N: OP_IF,
	}
}

func OpElse() Op {
	return Op{
		N: OP_ELSE,
	}
}

func OpEnd() Op {
	return Op{
		N: OP_END,
	}
}

func MakeBlocks(p Program) Program {
	// set blocks for jump operations like If etc
	stack := []int{}
	for i, op := range p {
		switch op.N {
		// in case of if we jump to end in case of false
		case OP_IF:
			stack = append(stack, i)
		case OP_ELSE:
			stack = append(stack, i)
		case OP_END:
			// pop()
			if_index := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			if p[if_index].N != OP_IF &&
				p[if_index].N != OP_ELSE {
				log.Panicf("Invalid end token %d", p[if_index].N)
			}
			p[if_index].JMP = i
		}
	}

	return p
}

func Tokenize(fpath string) (Program, error) {
	file, err := os.Open(fpath)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	p := Program{}
	scanner := bufio.NewScanner(file)
	split := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		advance, token, err = bufio.ScanWords(data, atEOF)
		return
	}

	scanner.Split(split)

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
		case S_OP_DIV:
			scanner.Scan()
			w := scanner.Text()
			t, err := strconv.Atoi(w)
			if err != nil {
				return p, err
			}
			o := Div(t)
			p = append(p, o)
		case S_OP_MUL:
			scanner.Scan()
			w := scanner.Text()
			t, err := strconv.Atoi(w)
			if err != nil {
				return p, err
			}
			o := Mul(t)
			p = append(p, o)
		case S_OP_MOD:
			scanner.Scan()
			w := scanner.Text()
			t, err := strconv.Atoi(w)
			if err != nil {
				return p, err
			}
			o := Mod(t)
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
		case S_OP_IF:
			// find end
			o := OpIf()
			p = append(p, o)
		case S_OP_ELSE:
			// find end
			o := OpElse()
			p = append(p, o)
		case S_OP_END:
			// find end
			o := OpEnd()
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

	p = MakeBlocks(p)

	return p, nil
}
