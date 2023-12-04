package parser

import (
	"fmt"
	"interperter/tokenizer"
	"os"
	"reflect"
	"testing"

	"github.com/kr/pretty"
)

func Equal(a, b []map[string]interface{}) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !reflect.DeepEqual(a[i], b[i]) {
			return false
		}
	}
	return true
}

func Assert(s, testToken []map[string]interface{}) {
	if !Equal(s, testToken) {
		(fmt.Printf("%# v\n %# v\n", pretty.Formatter(s), pretty.Formatter(testToken)))
		os.Exit(2)
	}
}
func Test_parse(t *testing.T) {
	program_tokens, _ := tokenizer.Tokenize("print (1);")
	ast, _ := Parse(program_tokens)
	//fmt.Printf("%# v\n", pretty.Formatter(ast))
	Assert(ast, []map[string]interface{}{
		{"type": "program"},
		{"statements": []map[string]interface{}{
			{
				"type": "print",
				"expression": []map[string]interface{}{
					{
						"type":  "number",
						"value": "1",
					},
				},
				"end": map[string]interface{}{
					"type":  "control",
					"value": "\n",
				},
			},
		},
		},
	},
	)

	program_tokens, _ = tokenizer.Tokenize("print (1); print (2 + 3);")
	ast, _ = Parse(program_tokens)
	Assert(ast, []map[string]interface{}{
		{"type": "program"},
		{"statements": []map[string]interface{}{
			{
				"type": "print",
				"expression": []map[string]interface{}{
					{
						"type":  "number",
						"value": "1",
					},
				},
				"end": map[string]interface{}{
					"type":  "control",
					"value": "\n",
				},
			},
			{
				"type": "print",

				"expression": []map[string]interface{}{
					{"type": "binary",
						"left": map[string]interface{}{

							"type":  "number",
							"value": "2",
						},
						"operator": tokenizer.Token{Name: "binary_operator", Value: "+"},
						"right": map[string]interface{}{

							"type":  "number",
							"value": "3",
						}},
				},
				"end": map[string]interface{}{
					"type":  "control",
					"value": "\n",
				},
			},
		},
		},
	},
	)
}
