package tokenizer

type TokenCode int

const (
	OP_PUSH = iota
	OP_PLUS
	OP_SUB
	OP_MUL
	OP_DIV
	OP_MOD
	OP_DUMP
	OP_VAR
	OP_EQUALS
	OP_GTE
	OP_GT
	OP_LTE
	OP_LT
	OP_IF
	OP_END
	OP_ELSE
	OP_MEM
	OP_MEM_STORE
	OP_MEM_LOAD
	OP_COUNT
)

const (
	S_OP_PLUS      = "+"
	S_OP_SUB       = "-"
	S_OP_MUL       = "*"
	S_OP_DIV       = "/"
	S_OP_MOD       = "%"
	S_OP_DUMP      = "dump"
	S_OP_EQUALS    = "=="
	S_OP_LTE       = "<="
	S_OP_LT        = "<"
	S_OP_GT        = ">"
	S_OP_GTE       = ">="
	S_OP_IF        = "if"
	S_OP_END       = "end"
	S_OP_MEM       = "mem"
	S_OP_MEM_STORE = ","
	S_OP_MEM_LOAD  = "."
	S_OP_ELSE      = "else"
)

// Token Type
type TokenType int

const (
	TOKEN_BI = iota
	TOKEN_LITERAL
	TOKEN_UNI
)

var OP_SYMBOLS map[string]TokenCode = map[string]TokenCode{
	S_OP_PLUS:   OP_PLUS,
	S_OP_SUB:    OP_SUB,
	S_OP_MUL:    OP_MUL,
	S_OP_DIV:    OP_DIV,
	S_OP_MOD:    OP_MOD,
	S_OP_DUMP:   OP_DUMP,
	S_OP_EQUALS: OP_EQUALS,
	S_OP_GTE:    OP_GTE,
	S_OP_GT:     OP_GT,
	S_OP_LT:     OP_LT,
	S_OP_LTE:    OP_LTE,
	S_OP_IF:     OP_IF,
	S_OP_END:    OP_END,
	S_OP_ELSE:   OP_ELSE,
}

var CMP_INS map[TokenCode]string = map[TokenCode]string{
	OP_EQUALS: "cmove",
	OP_LTE:    "cmovbe",
	OP_LT:     "cmovb",
	OP_GTE:    "cmovae",
	OP_GT:     "cmova",
}

var Stack []int

type Token struct {
	Type TokenType // type of token
	Code TokenCode // name
	O    int       // operand for the operation
	JMP  int       // jmp address to be used by some ops like If/else/while etc
	INS  string    // some instructions can contains their assembly instruction in case multiple op use same set of operations e.g comparison ==, <= etc
}

type Program []Token

func Push(x int) Token {
	return Token{
		Code: OP_PUSH,
		Type: TOKEN_UNI,
		O:    x,
	}
}

func Arithmatic(op TokenCode, x int) Token {
	return Token{
		Code: op,
		Type: TOKEN_BI,
		O:    x,
	}
}

func Dump() Token {
	return Token{
		Type: TOKEN_UNI,
		Code: OP_DUMP,
	}

}

func Cmp(op TokenCode, ins string, x int) Token {
	return Token{
		Code: op,
		Type: TOKEN_BI,
		O:    x,
		INS:  ins,
	}
}

func Gte(x int) Token {
	return Token{
		Type: TOKEN_BI,
		Code: OP_GTE,
		O:    x,
	}
}

func Lte(x int) Token {
	return Token{
		Type: TOKEN_BI,
		Code: OP_LTE,
		O:    x,
	}
}

func TokenIf() Token {
	return Token{
		Type: TOKEN_LITERAL,
		Code: OP_IF,
	}
}

func TokenElse() Token {
	return Token{
		Type: TOKEN_LITERAL,
		Code: OP_ELSE,
	}
}

func TokenEnd() Token {
	return Token{
		Type: TOKEN_LITERAL,
		Code: OP_END,
	}
}

func Mem() Token {
	return Token{
		Type: TOKEN_LITERAL,
		Code: OP_MEM,
	}
}

func MemStore(x int) Token {
	return Token{
		Type: TOKEN_LITERAL,
		Code: OP_MEM_STORE,
		O:    x,
	}
}

func MemLoad() Token {
	return Token{
		Type: TOKEN_LITERAL,
		Code: OP_MEM_LOAD,
	}
}
