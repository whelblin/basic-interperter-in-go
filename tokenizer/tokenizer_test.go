package tokenizer

import (
	"fmt"
	"os"
	"testing"
)

func Equal(a, b []Token) bool {
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

func Assert(s, testToken []Token) {
	if !Equal(s, testToken) {
		fmt.Println(s, testToken)
		os.Exit(2)
	}
}
func Test_simple_tokens(t *testing.T) {
	examples := [2]string{"+", "-"}
	var testToken Token
	for _, example := range examples {
		testToken = Token{"binary_operator", example}
		if !Equal(Tokenize(example), []Token{testToken}) {
			fmt.Println(Tokenize(example), testToken)
			t.Fatalf("Failed test")
		}
	}

}

func Test_number_tokens(t *testing.T) {
	for _, s := range [6]string{"1", "22", "12.1", "0", "12.", "123145"} {
		testToken := Token{"number", s}
		if !Equal(Tokenize(s), []Token{testToken}) {
			fmt.Println(Tokenize(s), testToken)
			t.Fatalf("Failed test")
		}
	}
}

func Test_string_tokens(t *testing.T) {
	for _, s := range [4]string{`"example"`, `"this is a longer example"`, `"an embedded "`, `" quote"`} {
		testToken := Token{"string", s}
		if !Equal(Tokenize(s), []Token{testToken}) {
			fmt.Println(Tokenize(s), testToken)
			t.Fatalf("Failed test")
		}
	}
}

func Test_identifier_tokens(t *testing.T) {
	for _, s := range [6]string{"x", "y", "z", "alpha", "beta", "gamma"} {
		testToken := Token{"identifier", s}
		if !Equal(Tokenize(s), []Token{testToken}) {
			fmt.Println(Tokenize(s), testToken)
			t.Fatalf("Failed test")
		}
	}
}

func Test_whitespace(t *testing.T) {
	for _, s := range [4]string{"1", "1  ", "  1", "  1  "} {
		testToken := Token{"number", "1"}
		if !Equal(Tokenize(s), []Token{testToken}) {
			fmt.Println(Tokenize(s), testToken)
			t.Fatalf("Failed test")
		}
	}
}

func Test_multiple_tokens(t *testing.T) {
	Assert(Tokenize("1+2"), []Token{{"number", "1"}, {"binary_operator", "+"}, {"number", "2"}})
	Assert(Tokenize("1+2-3"), []Token{{"number", "1"}, {"binary_operator", "+"}, {"number", "2"}, {"binary_operator", "-"}, {"number", "3"}})
	Assert(Tokenize("3+4*(5-2)"), []Token{
		{"number", "3"},
		{"binary_operator", "+"},
		{"number", "4"},
		{"binary_operator", "*"},
		{"left_parenthesis", "("},
		{"number", "5"},
		{"binary_operator", "-"},
		{"number", "2"},
		{"right_parenthesis", ")"}})
	/*
	  assert tokenize("3+4*(5-2)") == tokenize("3 + 4 * (5 - 2)")
	  assert tokenize("3+4*(5-2)") == tokenize("  3  +  4 * (5 - 2)  ")
	  assert tokenize("3+4*(5-2)") == tokenize(" 3 + 4 * (5 - 2) ")
	*/
}

/*
def test_keywords():
    print("testing keywords")
    for keyword in ["print","if","else","while"]:
        assert tokenize(keyword) == [keyword]

*/
