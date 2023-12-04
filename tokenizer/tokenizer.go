package tokenizer

import (
	"errors"
	"fmt"
	"regexp"
	"slices"
)

// class to store token data
// the name is the token type
// value is the value or name of the token
type Token struct {
	Name  string
	Value string
}

// patterns used to check the input data and assign tokens
var pattern = []map[string]string{
	{`^\s+|^\n|^\t`: "space"},
	{`^\+|^\-|^\*|^/`: "binary_operator"},
	{`^\d+(\.\d*)?`: "number"},
	{`^\n`: "newline"},
	{`^"([^"]|"")*"`: "string"},
	{`^<|^>`: "comparison"},
	{`^!`: "not"},
	{`^\(`: "left_parenthesis"},
	{`^\)`: "right_parenthesis"},
	{`^\{`: "left_curly_brace"},
	{`^\}`: "right_curly_brace"},
	{`^,`: "comma"},
	{`^\;`: "semicolon"},
	{`^print`: "print_statement"},
	{`^if`: "if"},
	{`^while`: "while"},
	{`^do`: "do"},
	{`^input`: "input"},
	{`^function`: "function"},
	{`^([a-zA-Z_][a-zA-Z0-9_]*)`: "identifier"}, //user variables and user fuction calls
	{`^={1}`: "assignment"},
}

// used to check to make sure the token type exisits

var tokenCheck = []string{"print_statement", "binary_operator", "number",
	"string", "comparison", "left_parenthesis",
	"right_parenthesis", "left_curly_brace",
	"right_curly_brace", "semicolon", "identifier", "assignment", "not", "comma", "print", "if", "while", "do", "function", "input", "newline"}

// The lex/tokenize function
/**
Goes though and finds the token assocated with the input text
goes character by character. Once it finds a match, add the token to the list
returns the list and any error codes
**/
func Tokenize(characters string) ([]Token, error) {
	var testTokens = []Token{}
	var i, v string
	pos := 0
	var submatches []int
	for pos < len(characters) {
		substr := characters[pos:]
		for num := range pattern {
			for i, v = range pattern[num] {
			} // gets the key and value
			re, _ := regexp.Compile(i)
			submatches = re.FindStringIndex(substr)
			// valid token will not be nil
			if submatches != nil {
				break
			}

		}
		// not allowed token
		if submatches == nil {
			return []Token{}, errors.New("not allowed token")
		}
		pos += submatches[1] // moves to get the next string
		// adds it to the token list if it is a valid word
		if slices.Contains(tokenCheck, v) {
			testTokens = append(testTokens, Token{v, substr[submatches[0]:submatches[1]]})
			continue
		} else if v == "space" {
			continue
		} else {
			fmt.Printf("Unknown type for %v", substr[submatches[0]:submatches[1]])
			return []Token{}, errors.New("unknown token type")
		}
	}
	return testTokens, nil
}
