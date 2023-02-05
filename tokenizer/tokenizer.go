package tokenizer

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Tokenizer struct {
	Index int
}

func MakeBlocks(p Program) Program {
	// set blocks for jump operations like If etc
	stack := []int{}
	for i, op := range p {
		switch op.Code {
		// in case of if we jump to end in case of false
		case OP_IF:
			stack = append(stack, i)
		case OP_ELSE:
			stack = append(stack, i)
		case OP_END:
			// pop()
			if_index := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			if p[if_index].Code != OP_IF &&
				p[if_index].Code != OP_ELSE {
				log.Panicf("Invalid end token %d", p[if_index].Code)
			}
			p[if_index].JMP = i
		}
	}

	return p
}

func (p *Tokenizer) GetNextToken(line []rune, endChar rune) string {
	res := []rune{}

	for line[p.Index] == ' ' {
		p.Index += 1
	}

	for i := p.Index; i < len(line) && line[i] != endChar; i++ {
		res = append(res, line[i])
	}

	p.Index += len(res)
	fmt.Println(string(res))
	return string(res)
}

func (p *Tokenizer) LexLine(line []rune) Program {
	var program Program
	word := ""
	// find each col in line seperated by space
	p.Index = 0
	for p.Index < len(line) {
		if line[p.Index] == '"' {
			word = p.GetNextToken(line, '"')
		} else {
			word = p.GetNextToken(line, ' ')
		}

		token, _, err := p.Tokenize(line, word)
		if err != nil {
			log.Panicf("Err %s on line number", err)
		}

		program = append(program, token)
	}

	return program
}

func (p *Tokenizer) LexFile(fpath string) (Program, error) {
	program := Program{}
	file, err := os.Open(fpath)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		p := p.LexLine([]rune(line))
		program = append(program, p...)
	}

	program = MakeBlocks(program)
	return program, nil
}

func (p *Tokenizer) Tokenize(l []rune, w string) (token Token, ahead string, err error) {
	switch w {
	case S_OP_PLUS, S_OP_DIV, S_OP_MUL, S_OP_SUB, S_OP_MOD:
		next := p.GetNextToken(l, ' ')
		t, err := strconv.Atoi(next)
		if err != nil {
			return token, next, err
		}

		s, ok := OP_SYMBOLS[w]
		if !ok {
			log.Panicf("No valid symbol for %s", w)
		}

		token = Arithmatic(s, t)

		return token, next, nil
	case S_OP_EQUALS, S_OP_GTE, S_OP_LTE, S_OP_GT, S_OP_LT:
		next := p.GetNextToken(l, ' ')
		t, err := strconv.Atoi(next)
		if err != nil {
			return token, next, err
		}

		OP, ok := OP_SYMBOLS[w]
		if !ok {
			return token, next, errors.New("Not a valid symbol " + w)
		}

		ins, ok := CMP_INS[OP]
		if !ok {
			return token, next, errors.New("Not a valid symbol " + w)
		}

		token = Cmp(OP, ins, t)

		return token, next, nil
	case S_OP_DUMP:
		token = Dump()
		return token, "", nil
	case S_OP_IF:
		// find end
		token = TokenIf()
		return token, "", nil
	case S_OP_ELSE:
		// find end
		token = TokenElse()
		return token, "", nil
	case S_OP_END:
		// find end
		token = TokenEnd()
		return token, "", nil
	case S_OP_MEM:
		token = Mem()
		return token, "", nil
	case S_OP_MEM_STORE:
		next := p.GetNextToken(l, ' ')
		t, err := strconv.Atoi(next)
		if err != nil {
			return token, "", err
		}

		token = MemStore(t)

		return token, next, nil
	case S_OP_MEM_LOAD:
		token = MemLoad()

		return token, "", nil
	default:
		t, err := strconv.Atoi(w)
		if err != nil {
			return token, "", err
		}
		token = Push(t)

		return token, "", nil
	}

	return token, "", nil
}
