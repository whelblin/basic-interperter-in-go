package parser

import (
	"interperter/tokenizer"
	"os"
)

var tokens []tokenizer.Token
var current_token_index int

func get_current_token() tokenizer.Token {
	if current_token_index < len(tokens) {
		return tokens[current_token_index]
	}
	return tokenizer.Token{"nil", "nil"}
}
func consume_token() {
	current_token_index += 1
}
func Parse(program_tokens []tokenizer.Token) []map[string]interface{} {
	current_token_index = 0
	tokens = program_tokens
	statements := []interface{}{}
	for get_current_token().Value != "nil" {
		statements = append(statements, parse_statement())
		//fmt.Println(statements)
	}
	//token = append(token, map[string]interface{}{"type": "program", "statements": statements})
	//fmt.Printf("\n")
	return []map[string]interface{}{{"type": "program"}, {"statements": statements}}
}

func parse_statement() interface{} {
	current_token := get_current_token()
	//fmt.Printf("%v\n", current_token)
	if current_token.Value == "print" {
		consume_token()
		expression := parse_expression()
		if get_current_token().Value != ";" {
			os.Exit(2)
		}
		consume_token()
		//fmt.Println("expression:", expression)
		return []map[string]interface{}{{"type": "print"}, {"expression": expression}}
	} else {
		consume_token()
	}
	return current_token_index
}

func parse_expression() any {
	left_term := parse_term()
	for get_current_token().Value == "+" || get_current_token().Value == "-" {
		operator := get_current_token()
		consume_token()
		right_term := parse_term()
		//left_term = [op, left_term, right_term]
		left_term = []map[string]interface{}{{"type": "binary"}, {"left": left_term}, {"operator": operator}, {"right": right_term}}
	}

	return left_term
}

// START HERE
func parse_term() any {
	left_factor := parse_factor()
	for get_current_token().Value == "*" || get_current_token().Value == "/" {
		operator := get_current_token()
		consume_token()
		right_factor := parse_factor()
		left_factor = []map[string]interface{}{{"type": "binary"}, {"left": left_factor}, {"operator": operator}, {"right": right_factor}}
	}
	return left_factor
}

func parse_factor() any {
	current_token := get_current_token()
	//fmt.Println("term: ", current_token)
	if current_token.Name == "number" {
		consume_token()
		return current_token.Value
	} else if current_token.Name == "binary_operator" {
		operator := get_current_token()
		consume_token()
		factor := parse_factor()
		return []map[string]interface{}{{"type": "unary"}, {"operator": operator}, {"expression": factor}}
	} else if current_token.Name == "left_parenthesis" {
		consume_token()
		expression := parse_expression()
		if get_current_token().Name != "right_parenthesis" {
			os.Exit(2)
		}
		consume_token()
		return expression
	} else {
		os.Exit(3)
	}
	return -1
}
