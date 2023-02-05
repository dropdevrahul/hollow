package main

import (
	"flag"
	"log"
	"os"

	"github.com/dropdevrahul/hollow/parser"
	"github.com/dropdevrahul/hollow/tokenizer"
)

func main() {
	log.SetOutput(os.Stdout)
	// Our language consists of a runtime stack we can modify
	// using various instructions
	output := flag.String("o", "output", "output file name")

	flag.Parse()

	t := tokenizer.Tokenizer{}

	program, err := t.LexFile(flag.Args()[0])
	if err != nil {
		log.Println(err)
	}

	parser.Compile(program, *output)
}
