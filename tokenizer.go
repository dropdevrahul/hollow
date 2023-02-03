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

func Var(v string, o int) Op {
	return Op{
		N: OP_VAR,
		O: o,
		V: v,
	}
}

func OpIf() Op {
	return Op{
		N: OP_IF,
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
		if op.N == OP_IF {
			stack = append(stack, i)
		} else if op.N == OP_END {
			// pop
			if_index := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			if p[if_index].N != OP_IF {
				log.Panic("No if block for end")
			}
			p[if_index].O = i
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
