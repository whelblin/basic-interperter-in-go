package main

import (
	"bufio"
	"fmt"
	evaluator "interperter/evaluator"
	"interperter/parser"
	"interperter/tokenizer"
	"os"

	"github.com/kr/pretty"
)

func main() {
	if len(os.Args) > 1 {
		args := os.Args[1:]
		if args != nil {
			dat, _ := os.ReadFile(args[0])
			fmt.Printf("")
			tokens, errCatch := tokenizer.Tokenize(string(dat))
			if errCatch != nil {
				fmt.Println("ERROR:", errCatch)
			}
			ast, errCatch := parser.Parse(tokens)
			if errCatch != nil {
				fmt.Println("ERROR:", errCatch)
			}
			evaluator.Evalute(ast)
		}
	} else {
		//runner()
		tokens, _ := tokenizer.Tokenize(`print 1 < 2;`)
		fmt.Println("Tokens", tokens)
		ast, _ := parser.Parse(tokens)
		fmt.Printf("%# v\n", pretty.Formatter(ast))
		evaluator.Evalute(ast)
	}
}

func runner() {
	fmt.Printf("starting\n")
	reader := bufio.NewReader(os.Stdin)
	for {
		var source_code string
		fmt.Printf(">> ")
		source_code, _ = reader.ReadString('\n')
		if source_code == "exit\n" {
			fmt.Println("Exiting...")
			os.Exit(0)
		}
		tokens, errCatch := tokenizer.Tokenize(source_code)
		if errCatch != nil {
			fmt.Println("ERROR:", errCatch)
		}
		ast, errCatch := parser.Parse(tokens)
		if errCatch != nil {
			fmt.Println("ERROR:", errCatch)
		}
		evaluator.Evalute(ast)
	}
}
