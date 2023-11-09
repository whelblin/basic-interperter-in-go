package tokenizer

import (
	"fmt"
	"os"
	"regexp"
	"slices"
)

type Token struct {
	Name  string
	Value string
}

var pattern = map[string]string{
	`^\s+`:                      "space",
	`^\+|^\-|^\*|^/`:            "binary_operator",
	`^\d+(\.\d*)?`:              "number",
	`^"([^"]|"")*"`:             "string",
	`^<|^>|^<=|^>=|^==|^!=`:     "comparison",
	`^\(`:                       "left_parenthesis",
	`^\)`:                       "right_parenthesis",
	`^\{`:                       "left_curly_brace",
	`^\}`:                       "right_curly_brace",
	`^\;`:                       "semicolon",
	`^([a-zA-Z_][a-zA-Z0-9_]*)`: "identifier",
	`^=`:                        "assignment",
	//`^print`:                  "print_statement",
}

var tokenCheck = []string{"print_statement", "binary_operator", "number",
	"string", "comparison", "left_parenthesis",
	"right_parenthesis", "left_curly_brace",
	"right_curly_brace", "semicolon", "identifier", "assignment"}

// The lex/tokenize function

func Tokenize(characters string) []Token {
	var testTokens = []Token{}
	fmt.Println("tokenizing", characters)
	var i, v string
	pos := 0
	var submatches []int
	for pos < len(characters) {
		substr := characters[pos:]
		//fmt.Println("substr:", substr)
		for i, v = range pattern {
			re, _ := regexp.Compile(i)
			submatches = re.FindStringIndex(substr)
			//fmt.Println("\tsubmatches", submatches, submatches == nil)

			if !(submatches == nil) {
				break
			}
		}
		// not allowed token
		if submatches == nil {
			fmt.Printf("Not allowed token: %v at pos %v >> ", characters, pos)
			os.Exit(1)
		}
		pos += submatches[1]
		//fmt.Printf("substr: %v\n", substr)
		if slices.Contains(tokenCheck, v) {
			testTokens = append(testTokens, Token{v, substr[submatches[0]:submatches[1]]})
			continue
		} else if v == "space" {
			continue
		} else {
			fmt.Printf("Unknown type for %v", substr[submatches[0]:submatches[1]])
			os.Exit(1)
		}
		testTokens = append(testTokens, Token{"", substr[submatches[0]:submatches[1]]})
	}
	//fmt.Println(testTokens)
	return testTokens
	//tokens := [0]string{}
	//pos := 0

}
