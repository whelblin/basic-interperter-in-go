package parser

import (
	"errors"
	"interperter/tokenizer"
	"os"
)

var tokens []tokenizer.Token
var current_token_index int

/*
@input: void
@return: a Token type
@info:
  - returns the current token by its index in the array
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

/*
@input: an array of Tokens given by the tokenizer function
@return: the ast in the form of a array of maped values by their type, error code
@info:
  - parses each token in the input array and adds it to the ast which is all wrapped
    under the statements section

@uses:
  - parse_statement(): to parse each statements
  - get_current_token(): gets the current token in the array
*/
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

/*
@input: void
@return: a sub tree of the ast, error code
@info:

  - finds the type of the current token and builds the sub tree
    for that statement
  - might consume multiple tokens for each statement
  - checks for invaild syntax and returns an error if needed

@uses:

  - parse_expression(): to parse expressions
  - parse_statement(): to parse statements within if,else,while loops
  - parse_parameters(): to parse function paramaters
  - get_current_token(): gets the current token in the array
  - consume_token(): moves to the next token
  - parse_block(): parses a block of code surrounded in {}
*/
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
			return map[string]interface{}{"type": "assignment", "name": name, "expression": expression}, nil
		}

	} else {
		return map[string]interface{}{}, errors.New("invalid statement")
	}
}

/*
@input: void
@output: the subtree for the ast for the block of statments, error code
@info:
  - goes through each token and parses each statements
  - each iteration through the loop might consume multiple tokens

@uses:
  - consume_token(): moves to the next token
  - parse_statement(): parses each statements starting from the current token
  - get_current_token(): gets the current token in the array
*/
func parse_block() map[string]interface{} {
	consume_token()
	statements := []map[string]interface{}{}
	for get_current_token().Value != "}" {
		result, _ := parse_statement()
		statements = append(statements, result)
	}
	consume_token()
	return map[string]interface{}{"type": "block", "statements": statements}
}

/*
@input: void
@output: the sub tree for the expression
@info
  - parses the left term first
  - if there is an operator, it will also parse the right term and return the tree
    --order of operations: *,/,+,-, <, <=, >, >=, ==, !=
    --the * and / are done in the parse_term to ensure it is done first
  - if there is no valid operator it will just return the left term

@uses:
  - consume_token(): moves to the next token
  - get_current_token(): gets the current token in the array
  - parse_term(): used to parse the left and right terms
*/
func parse_expression() map[string]interface{} {
	left_term := parse_term()

	for get_current_token().Value == "+" || get_current_token().Value == "-" {
		operator := get_current_token()
		consume_token()
		right_term := parse_term()
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

/*
@input: void
@output: the sub tree for the term
@info
  - parses the left tfactor first
  - if there is an operator, it will also parse the right factor and return the tree
    --order of operations: *,/,+,-, <, <=, >, >=, ==, !=
    --the * and / are done here instead of parse_term so they are done first
  - if there is no valid operator it will just return the left term

@uses:
  - consume_token(): moves to the next token
  - get_current_token(): gets the current token in the array
  - parse_factor(): used to parse the left and right factors
*/
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

/*
@input: void
@output: the sub tree for the factor
@info
  - factor can be a:
    -- number
    -- binary_operator, which means it is a unary operation because the operator is the left factor
    -- left_parenthesis for an expression
    -- string
    -- identifier, which can be a user varaible or function call

@uses:
  - consume_token(): moves to the next token
  - get_current_token(): gets the current token in the array
  - parse_function_call(): used to parse a function call
*/
func parse_factor() map[string]interface{} {
	current_token := get_current_token()
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
		return map[string]interface{}{current_token.Name: current_token.Value}
	} else {
		os.Exit(2)
	}
	return map[string]interface{}{}
}

/*
@input: void
@output: the array for the parameters(where each element is a sub tree of the parameter), error code
@info
  - parses each parameter and adds it to the array of parameters

@uses:
  - consume_token(): moves to the next token
  - get_current_token(): gets the current token in the array
*/
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

/*
@input: void
@output: the array for the arguments(where each element is a sub tree of the argument), error code
@info
  - parses each argument and adds it to the array of arguments

@uses:
  - consume_token(): moves to the next token
  - get_current_token(): gets the current token in the array
*/
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

/*
@input: function name
@output: subtree for the function call, error code
@info
  - parses each argument and returns the sub tree with the function name, and arguments

@uses:
  - consume_token(): moves to the next token
  - parse_arguments: parses the arguments and returns an array
*/
func parse_function_call(identifier string) (map[string]interface{}, error) {
	consume_token() // "("
	arguments, _ := parse_arguments()
	consume_token()
	return map[string]interface{}{"type": "function_call", "name": identifier, "arguments": arguments}, nil

}
