package tokenizer

import (
	"fmt"
	"testing"
)

func Equal(a, b []token) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
func Test_simple_tokens(t *testing.T) {
	fmt.Println("testing simple tokens")
	examples := [3]string{"+", "-"}
	var testToken token
	for _, example := range examples {
		testToken = token{"binary_operator", example}
		if !Equal(Tokenize(example), []token{testToken}) {
			fmt.Println(Tokenize(example), testToken)
			t.Fatalf("Failed test")
		}
	}

}

/*
func Test_number_tokens(t *testing.T){
    fmt.Println("testing number tokens")
    for s := range [6]int {"1", "22", "12.1", "0", "12.", "123145"}{
        if tokenize(s) != [
            ["number", number(s)]
        ], f"Expected {[['number', s]]}, got {tokenize(s)}."
	}
}
*/
