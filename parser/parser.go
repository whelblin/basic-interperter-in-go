package parser

import (
	"errors"
	"fmt"
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
func Parse(program_tokens []tokenizer.Token) ([]map[string]interface{}, error) {
	current_token_index = 0
	tokens = program_tokens
	var errors error
	var result map[string]interface{}
	statements := []map[string]interface{}{}
	for get_current_token().Value != "nil" {
		result, errors = parse_statement()
		statements = append(statements, result)
	}
	return []map[string]interface{}{{"type": "program"}, {"statements": statements}}, errors
}

func parse_statement() (map[string]interface{}, error) {
	current_token := get_current_token()
	if current_token.Value == "print" {
		consume_token()
		expression := parse_expression()
		if get_current_token().Value != ";" {
			return map[string]interface{}{}, errors.New("missing ;")
		}
		consume_token()
		return map[string]interface{}{"type": "print", "expression": expression}, nil
	} else if current_token.Value == "{" {
		return parse_block(), nil

	} else if current_token.Name == "identifier" {
		// if statement keyword
		if current_token.Value == "if" {
			consume_token()
			if get_current_token().Value != "(" {
				return map[string]interface{}{}, errors.New("expected '('")
			}
			consume_token()
			condition := parse_expression()
			if get_current_token().Value != ")" {
				return map[string]interface{}{}, errors.New("expected ')'")
			}

			consume_token()
			then_statement, _ := parse_statement()
			var else_statement interface{}
			if get_current_token().Value == "else" {
				consume_token()
				else_statement, _ = parse_statement()
			} else {
				else_statement = nil
			}
			return map[string]interface{}{"type": "if", "condition": condition, "then": then_statement, "else": else_statement}, nil

		} else if current_token.Value == "while" {
			consume_token()
			if get_current_token().Value != "(" {
				return map[string]interface{}{}, errors.New("expected '('")
			}
			consume_token()
			condition := parse_expression()
			if get_current_token().Value != ")" {
				return map[string]interface{}{}, errors.New("expected ')'")
			}
			consume_token()
			do_statement, _ := parse_statement()
			return map[string]interface{}{"type": "while",
				"condition": condition,
				"do":        do_statement}, nil
		} else { // user variable
			name := current_token.Value
			consume_token()
			if get_current_token().Value != "=" {
				return map[string]interface{}{}, errors.New("missing =")

			}
			consume_token()
			expression := parse_expression()
			if get_current_token().Value != ";" {
				return map[string]interface{}{}, errors.New("missing ;")
			}
			consume_token()
			//fmt.Println("expression:", expression)
			return map[string]interface{}{"type": "assignment", "name": name, "expression": expression}, nil
		}
	} else {
		return map[string]interface{}{}, errors.New("invalid statement")
	}
	//return map[string]interface{}{}, errors.New("invalid statement")
}

func parse_block() map[string]interface{} {
	consume_token()
	statements := []map[string]interface{}{}
	for get_current_token().Value != "}" {
		result, _ := parse_statement()
		statements = append(statements, result)
		//fmt.Println(statements)
	}
	consume_token()
	return map[string]interface{}{"type": "block", "statements": statements}
}

func parse_expression() map[string]interface{} {
	left_term := parse_term()
	fmt.Println(left_term)
	fmt.Println(get_current_token())
	for get_current_token().Value == "+" || get_current_token().Value == "-" {
		operator := get_current_token()
		consume_token()
		right_term := parse_term()
		//left_term = [op, left_term, right_term]
		left_term = map[string]interface{}{"type": "binary", "left": left_term, "operator": operator, "right": right_term}
	}
	for get_current_token().Value == "<" || get_current_token().Value == ">" ||
		get_current_token().Value == "<=" || get_current_token().Value == ">=" ||
		get_current_token().Value == "==" || get_current_token().Value == "!=" {
		operator := get_current_token()
		consume_token()
		right_term := parse_term()
		left_term = map[string]interface{}{"type": "comparison", "left": left_term, "operator": operator, "right": right_term}
	}

	return left_term
}

// START HERE
func parse_term() map[string]interface{} {
	left_factor := parse_factor()
	for get_current_token().Value == "*" || get_current_token().Value == "/" {
		operator := get_current_token()
		consume_token()
		right_factor := parse_factor()
		left_factor = map[string]interface{}{"type": "binary", "left": left_factor, "operator": operator, "right": right_factor}
	}
	return left_factor
}

func parse_factor() map[string]interface{} {
	current_token := get_current_token()
	//fmt.Println("term: ", current_token)
	if current_token.Name == "number" {
		consume_token()
		return map[string]interface{}{current_token.Name: current_token.Value}
	} else if current_token.Name == "binary_operator" {
		operator := get_current_token()
		consume_token()
		factor := parse_factor()
		return map[string]interface{}{"type": "unary", "operator": operator, "expression": factor}
	} else if current_token.Name == "left_parenthesis" {
		consume_token()
		expression := parse_expression()
		if get_current_token().Name != "right_parenthesis" {
			os.Exit(2)
		}
		consume_token()
		return expression
	} else if current_token.Name == "identifier" {
		consume_token()
		return map[string]interface{}{"type": "identifier", "name": current_token.Value}

	} else if current_token.Name == "string" {
		consume_token()
		return map[string]interface{}{current_token.Name: current_token.Value}
	} else {
		os.Exit(3)
	}
	return map[string]interface{}{}
}
