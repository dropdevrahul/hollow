package main

import (
	"flag"
	"log"
	"os"
)

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
