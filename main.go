package main

import (
	"fmt"
	"interperter/tokenizer"
)

func main() {
	fmt.Printf("starting\n")
	tokenizer.Tokenize("3 - 4")
	tokenizer.Tokenize(`"hello"`)
	//test_number_tokens()
	//test_string_tokens()
	//test_identifier_tokens()
	//test_whitespace()
	//test_multiple_tokens()
	//test_keywords()
}
