package tokenizer

import (
	"fmt"
	"os"
	"testing"
)

/*
*
helper function used to test if two tokens are the same
*
*/
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

/*
*
helper function that asserts if two functions are the same
*
*/
func Assert(s string, testToken []Token) {
	i, _ := Tokenize(s)
	if !Equal(i, testToken) {
		fmt.Println(i, testToken)
		os.Exit(1)
	}
}

func Test_simple_tokens(t *testing.T) {
	examples := [2]string{"+", "-"}
	var testToken Token
	for _, example := range examples {
		testToken = Token{"binary_operator", example}
		i, _ := Tokenize(example)
		if !Equal(i, []Token{testToken}) {
			i, _ := Tokenize(example)
			fmt.Println(i, testToken)
			t.Fatalf("Failed test")
		}
	}

}

func Test_number_tokens(t *testing.T) {
	for _, s := range [6]string{"1", "22", "12.1", "0", "12.", "123145"} {
		testToken := []Token{{"number", s}}
		Assert(s, testToken)
	}
}

func Test_string_tokens(t *testing.T) {
	for _, s := range [4]string{`"example"`, `"this is a longer example"`, `"an embedded "`, `" quote"`} {
		testToken := []Token{{"string", s}}
		Assert(s, testToken)
	}
}

func Test_identifier_tokens(t *testing.T) {
	for _, s := range [6]string{"x", "y", "z", "alpha", "beta", "gamma"} {
		testToken := []Token{{"identifier", s}}
		Assert(s, testToken)
	}
}

func Test_whitespace(t *testing.T) {
	for _, s := range [4]string{"1", "1  ", "  1", "  1  "} {
		testToken := []Token{{"number", "1"}}
		Assert(s, testToken)
	}
}

func Test_multiple_tokens(t *testing.T) {
	Assert(("1+2"), []Token{{"number", "1"}, {"binary_operator", "+"}, {"number", "2"}})
	Assert(("1+2-3"), []Token{{"number", "1"}, {"binary_operator", "+"}, {"number", "2"}, {"binary_operator", "-"}, {"number", "3"}})
	Assert(("3+4*(5-2)"), []Token{
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
