package tokenizer

import (
	"fmt"
	"os"
	"regexp"
)

func Test() {
	fmt.Printf("testing token")
}

type token struct {
	name  string
	value string
}

var pattern = map[string]string{
	`^\s+`:         "space",
	`^\+|^\-`:      "binary_operator",
	`^\d+(\.\d*)?`: "number",
	`"([^"]|"")*"`: "string",
	//`.`:            "error",
}

// The lex/tokenize function

func Tokenize(characters string) []token {
	var testTokens = []token{}
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
		if v == "number" || v == "binary_operator" || v == "string" {
			testTokens = append(testTokens, token{v, substr[submatches[0]:submatches[1]]})
			continue
		}
		if v == "space" {
			continue
		}
		testTokens = append(testTokens, token{"", substr[submatches[0]:submatches[1]]})
	}
	fmt.Println(testTokens)
	return testTokens
	//tokens := [0]string{}
	//pos := 0

}

/*
def test_string_tokens():
    print("testing string tokens")
    for s in ['"example"', '"this is a longer example"', '"an embedded "" quote"']:
        # adjust for the embedded quote behaviour
        r = s[1:-1].replace('""','"')
        assert tokenize(s) == [
            ["string", r]
        ], f"Expected {[['string', r]]}, got {tokenize(s)}."

def test_identifier_tokens():
    print("testing identifier tokens")
    for s in ["x", "y", "z", "alpha", "beta", "gamma"]:
        assert tokenize(s) == [
            ["identifier", s]
        ], f"Expected {[['identifier', s]]}, got {tokenize(s)}."

def test_whitespace():
    print("testing whitespace")
    for s in ["1", "1  ", "  1", "  1  "]:
        assert tokenize(s) == [["number", 1]]

def test_multiple_tokens():
    print("testing multiple tokens")
    assert tokenize("1+2") == [["number", 1], "+", ["number", 2]]
    assert tokenize("1+2-3") == [["number", 1], "+", ["number", 2], "-", ["number", 3]]
    assert tokenize("3+4*(5-2)") == [
        ["number", 3],
        "+",
        ["number", 4],
        "*",
        "(",
        ["number", 5],
        "-",
        ["number", 2],
        ")",
    ]
    assert tokenize("3+4*(5-2)") == tokenize("3 + 4 * (5 - 2)")
    assert tokenize("3+4*(5-2)") == tokenize("  3  +  4 * (5 - 2)  ")
    assert tokenize("3+4*(5-2)") == tokenize(" 3 + 4 * (5 - 2) ")


def test_keywords():
    print("testing keywords")
    for keyword in ["print","if","else","while"]:
        assert tokenize(keyword) == [keyword]

*/
