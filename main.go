package main

import (
	"fmt"
	"interperter/parser"
	"interperter/tokenizer"

	"github.com/kr/pretty"
)

func main() {
	fmt.Printf("starting\n")
	//fmt.Println(tokenizer.Tokenize("print 1 + 3"))
	test := tokenizer.Tokenize("print 1*2; print 1+3;")
	b := (parser.Parse(test))
	fmt.Printf("%# v\n", pretty.Formatter(b))
	//test_number_tokens()
	//test_string_tokens()
	//test_identifier_tokens()
	//test_whitespace()
	//test_multiple_tokens()
	//test_keywords()
}
