package parser

import (
	"errors"
	"fmt"
	"interperter/tokenizer"
	"os"
)

var tokens []tokenizer.Token
var current_token_index int

/*
*
returns the current token by its index in the array
*
*/
func get_current_token() tokenizer.Token {
	if current_token_index < len(tokens) {
		return tokens[current_token_index]
	}
	return tokenizer.Token{Name: "nil", Value: "nil"}
}

/*
*
adds one to the current token index to move to the next token
*
*/
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
	if current_token.Name == "print_statement" {
		consume_token()
		if get_current_token().Value != "(" {
			return map[string]interface{}{}, errors.New("expected '('")
		}
		consume_token()
		var expression map[string]interface{}
		var print_items []map[string]interface{}
		end := map[string]interface{}{"string": "\n"}
		for get_current_token().Value != ")" {
			if get_current_token().Value != "," {
				if get_current_token().Value == "end" {
					consume_token()
					if get_current_token().Value != "=" {
						return map[string]interface{}{}, errors.New("expected '='")
					}
					consume_token()
					end = parse_expression()

				} else {
					expression = parse_expression()
					print_items = append(print_items, expression)
				}
			} else {
				consume_token()
			}
		}
		consume_token()
		if get_current_token().Value != ";" {
			return map[string]interface{}{}, errors.New("missing ;")
		}
		consume_token()
		fmt.Println(print_items)
		return map[string]interface{}{"type": "print", "expression": print_items, "end": end}, nil
	} else if current_token.Value == "{" {
		return parse_block(), nil

	} else if current_token.Value == "if" {
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
	} else if current_token.Value == "do" {
		consume_token()
		if get_current_token().Value != "{" {
			return map[string]interface{}{}, errors.New("expected '{'")
		}
		do := parse_block()
		fmt.Println(do)
		if get_current_token().Value != "while" {
			return map[string]interface{}{}, errors.New("expected 'while'")
		}
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
		if get_current_token().Value != ";" {
			return map[string]interface{}{}, errors.New("missing ;")
		}
		consume_token()
		return map[string]interface{}{"type": "do_statement", "do": do, "condition": condition}, nil
	} else if current_token.Value == "function" {
		consume_token()
		name := get_current_token().Value
		consume_token()
		if get_current_token().Value != "(" {
			return map[string]interface{}{}, errors.New("expected '('")
		}
		consume_token()
		parameters, _ := parse_parameters()
		consume_token()
		body := parse_block()
		return map[string]interface{}{"type": "function", "name": name, "parameters": parameters, "body": body}, nil

	} else if current_token.Name == "identifier" {
		name := current_token.Value
		consume_token()
		if get_current_token().Value == "(" { // function call
			consume_token()
			var parameters []map[string]interface{}
			for get_current_token().Value != ")" {
				if get_current_token().Value != "," {
					parameters = append(parameters, map[string]interface{}{get_current_token().Name: get_current_token().Value})
				}
				consume_token()
			}
			consume_token()
			if get_current_token().Value != ";" {
				return map[string]interface{}{}, errors.New("missing ;")
			}
			consume_token()
			return map[string]interface{}{"type": "function_call", "name": name, "parameters": parameters}, nil
		} else {
			// user variable

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

	for get_current_token().Value == "+" || get_current_token().Value == "-" {
		operator := get_current_token()
		consume_token()
		right_term := parse_term()
		//left_term = [op, left_term, right_term]
		left_term = map[string]interface{}{"type": "binary", "left": left_term, "operator": operator, "right": right_term}
	}
	for get_current_token().Value == "<" || get_current_token().Value == ">" {
		operator := get_current_token()
		consume_token()
		if get_current_token().Value == "=" {
			operator.Value += "="
			consume_token()
		}
		right_term := parse_term()
		left_term = map[string]interface{}{"type": "comparison", "left": left_term, "operator": operator, "right": right_term}
	}

	for get_current_token().Value == "=" || get_current_token().Value == "!" {
		operator := get_current_token()
		consume_token()
		if get_current_token().Value != "=" {
			os.Exit(2)
		}
		operator.Value += get_current_token().Value
		consume_token()
		right_term := parse_term()
		left_term = map[string]interface{}{"type": "comparison", "left": left_term, "operator": operator, "right": right_term}

	}

	return left_term
}

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
			os.Exit(4)
		}
		consume_token()
		return expression
	} else if current_token.Name == "identifier" {
		identifier := current_token.Value
		consume_token()
		if get_current_token().Value == "(" {
			parse_function_call(identifier)
		} else {
			return map[string]interface{}{"type": "identifier", "name": identifier}
		}

	} else if current_token.Name == "string" {
		consume_token()
		fmt.Println("string", current_token.Value)
		if current_token.Value == "\"\\n\"" {
			return map[string]interface{}{"control": "\n"}
		} else {
			return map[string]interface{}{current_token.Name: current_token.Value}
		}
	} else if current_token.Name == "newline" {
		consume_token()
		return map[string]interface{}{"control": "\n"}
	} else {
		os.Exit(2)
	}
	return map[string]interface{}{}
}
func parse_parameters() ([]map[string]interface{}, error) {
	var parameters []map[string]interface{}
	for get_current_token().Value != ")" {
		parameters = append(parameters, parse_expression())
		if get_current_token().Value != ")" {
			if get_current_token().Value != "," {
				return []map[string]interface{}{}, errors.New("expected ', or )'")
			} else {
				consume_token()
			}
		}

	}
	return parameters, nil
}
func parse_arguments() ([]map[string]interface{}, error) {
	var arguments []map[string]interface{}
	for get_current_token().Value != ")" {
		arguments = append(arguments, parse_expression())
		if get_current_token().Value != ")" {
			if get_current_token().Value != "," {
				return []map[string]interface{}{}, errors.New("expected ', or )'")
			} else {
				consume_token()
			}
		}

	}
	return arguments, nil
}
func parse_function_call(identifier string) (map[string]interface{}, error) {
	consume_token() // "("
	arguments, _ := parse_arguments()
	consume_token()
	return map[string]interface{}{"type": "function_call", "name": identifier, "arguments": arguments}, nil

}
