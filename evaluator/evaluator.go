package evaluator

import (
	"errors"
	"fmt"
	"interpreter/tokenizer"
	"reflect"
	"strconv"
	"strings"
)

// returns the first value
// used for multiple return value functions
func first(n any, _ error) any {
	return n
}

// binary operations for numbers
var binary_operations = map[string]func(float64, float64) float64{
	"+": func(x, y float64) float64 { return x + y },
	"-": func(x, y float64) float64 { return x - y },
	"*": func(x, y float64) float64 { return x * y },
	"/": func(x, y float64) float64 { return x / y },
}

// binary comparisons for numbers
var binary_comparisons = map[string]func(float64, float64) bool{
	"<":  func(x, y float64) bool { return x < y },
	"<=": func(x, y float64) bool { return x <= y },
	">":  func(x, y float64) bool { return x > y },
	">=": func(x, y float64) bool { return x >= y },
	"==": func(x, y float64) bool { return x == y },
	"!=": func(x, y float64) bool { return x != y },
}

// binary comparisons for strings
var binary_comparisons_string = map[string]func(string, string) bool{
	"<":  func(x, y string) bool { return x < y },
	"<=": func(x, y string) bool { return x <= y },
	">":  func(x, y string) bool { return x > y },
	">=": func(x, y string) bool { return x >= y },
	"==": func(x, y string) bool { return x == y },
	"!=": func(x, y string) bool { return x != y },
}
var unary_operations = map[string]func(float64) float64{
	"-": func(f float64) float64 { return -f },
}

// environment that maps a name to a value of any type
var env = map[string]interface{}{}

/*
@input: ast tree from the Parse function
@output: the last statement in the tree to be executed
@info:
  - goes through the ast tree and executes each statement
  - if statement sections out first statements and block statements

@uses:
  - evaluate_statement(): evaluates each sub tree under the statements tag
*/
func Evaluate(node []map[string]interface{}) (any, error) {
	var index int = 0
	var error_code error
	var result any
	if len(node) > 0 {
		current_node := node[index]
		if current_node["type"] == "program" {
			index += 1
			current_node = node[index]
			current_array := current_node["statements"].([]map[string]interface{})
			for _, v := range current_array {
				result, error_code = evaluate_statement(v)
				if error_code != nil {
					return result, error_code
				}
			}

		} else {
			for _, v := range node {
				result, error_code = evaluate_statement(v)
				if error_code != nil {
					return result, error_code
				}
			}
		}
	}
	return result, error_code
}

/*
@input: list of items to print, end delimiter used at the end of each print
@output: array of the outputted items. This includes the delimiter
@info:
  - prints out each item given to the print function with the delimiter
  - the default delimiter is \n

@uses:
  - evaluate_statement(): for the custom delimiter(should be a string)
*/
func evaluate_print(list []map[string]interface{}, end map[string]interface{}) ([]interface{}, error) {
	var result []interface{}
	delimiter := "\n"
	if end["type"] != "control" {
		delimiter = first(evaluate_statement(end)).(string)
	} else {
		delimiter = end["value"].(string)
	}
	for i, x := range list {
		item, error_code := evaluate_statement(x)
		if error_code != nil {
			return result, error_code
		}
		result = append(result, item)
		fmt.Printf("%v%v", result[i], delimiter)
	}
	result = append(result, delimiter)
	return result, nil
}

/*
@input: binary operation, lhs(number), rhs(number),
@output: floating point result of the operation
@info:
  - performs the binary operation on x and y and returns the result
  - evaluates the x and y so their can be nested operations

@uses:
  - binary_operations(): to preform the operation on the x and y
*/
func evaluate_binary(op string, x, y map[string]interface{}) (float64, error) {

	evaluated_x, err_x := evaluate_statement(x)
	evaluated_y, err_y := evaluate_statement(y)
	if err_x != nil {
		return binary_operations[op](evaluated_x.(float64), evaluated_y.(float64)), err_x
	} else if err_y != nil {
		return binary_operations[op](evaluated_x.(float64), evaluated_y.(float64)), err_y
	}
	return binary_operations[op](evaluated_x.(float64), evaluated_y.(float64)), nil

}

/*
@input: unary operation, lhs(number)
@output: floating point result of the operation
@info:
  - performs the unary operation on x and returns the result
  - evaluates the x so their can be nested operations

@uses:
  - unary_operations(): to preform the operation on the x7
*/
func evaluate_unary(op string, x map[string]interface{}) (float64, error) {

	evaluated_x, error_code := evaluate_statement(x)
	return unary_operations[op](evaluated_x.(float64)), error_code

}

/*
@input: variable name, value x
@output: the value of the new variable in the env
@info:
  - evaluates the value of x and stores it in the environment mapped with the name parameter
  - returning the value allows for compound assignment

@uses:
  - evaluate_statement(): to evaluate the value of x to allow for complex assignment
*/
func evaluate_assignment(name string, x map[string]interface{}) (any, error) {
	evaluated_x, error_code := evaluate_statement(x)
	env[name] = evaluated_x
	return env[name], error_code
}

/*
@input: condition, then statement, else statement
@output: evaluated statement(then, else or nil)
@info:
  - first it evaluates the condition and branches  to the then or else based on it
  - a then statement is required but an else statement is not
  - it will evaluate either the then or else but not both statements
  - else statement is any because it can be nil or a sub tree

@uses:
  - evaluate_condition(): to evaluate the condition
  - evaluate_statement(): to  evaluate the then or else statement
*/
func evaluate_if(condition, then map[string]interface{}, else_statement any) (any, error) {
	con, error_code := evaluate_condition(condition)
	if error_code != nil {
		return con, error_code
	}
	if then == nil {
		return nil, errors.New("expected then statement")
	}
	if con {
		return evaluate_statement(then)
	} else {
		if else_statement != nil {
			return evaluate_statement(else_statement.(map[string]interface{}))
		} else {
			return nil, nil
		}
	}

}

/*
@input: condition, then statement
@output: the last evaluated evaluated statement or nil
@info:
  - first it evaluates the condition and loops based on it
  - it will evaluate the condition first and then evaluate the then statement
    the then statement will not be evaluated if the condition is false at first

@uses:
  - evaluate_condition(): to evaluate the condition
  - evaluate_statement(): to  evaluate the then statement
*/
func evaluate_while(condition, then_statement map[string]interface{}) (any, error) {
	var result interface{} = nil
	if then_statement == nil {
		return nil, errors.New("expected then statement")
	}
	value, error_code := evaluate_condition(condition)
	for value {
		if error_code != nil {
			return condition, error_code
		}
		result = first(evaluate_statement(then_statement))
		value, error_code = evaluate_condition(condition)
	}
	return result, error_code
}

/*
@input: condition, then statement
@output: the last evaluated evaluated statement
@info:
  - first it evaluates the then statement, then will loop based on the condition
  - the then statement will always be evaluated at least once and will always return a value

@uses:
  - evaluate_condition(): to evaluate the condition
  - evaluate_statement(): to  evaluate the then statement
*/
func evaluate_do(condition, then_statement map[string]interface{}) (any, error) {
	if then_statement == nil {
		return nil, errors.New("expected then statement")
	}
	result, error_code := (evaluate_statement(then_statement))
	value, con_error_code := evaluate_condition(condition)
	for value {
		if error_code != nil {
			return then_statement, error_code
		}
		if con_error_code != nil {
			return condition, con_error_code
		}
		result, error_code = evaluate_statement(then_statement)
		value, con_error_code = evaluate_condition(condition)
	}
	if error_code == nil {
		return result, con_error_code
	}
	return result, error_code
}

/*
@input: condition to be evaluated
@output: the evaluated condition (true or false)
@info:
  - it checks if the type of condition is a comparison.
    -- if so, returns the statement based on a boolean value
    -- otherwise, the condition is a number and checks if it is greater than 0 for true

@uses:
  - evaluate_statement(): to  evaluate the then statement
*/
func evaluate_condition(condition map[string]interface{}) (bool, error) {
	if condition["type"] == "comparison" {
		x, e := evaluate_statement(condition)
		return x.(bool), e

	}
	x, e := evaluate_statement(condition)
	return x.(float64) > 0, e

}

/*
@IN PROGRESS: TRUE
@input: function name, parameter array, and body
@output: the function assigned in the env
@info:
  - Still need to code

@uses:
*/
func function_declaration(name string, parameters []string, body map[string]interface{}) any {
	return 1

}

/*
@input: sub tree of the ast for a statement
@output: the last evaluated evaluated statement
@info:
  - every sub tree has a tag of type which has the type of node it is
  - depending on its type, it calls the corresponding function to handle the statement

@uses:
  - evaluate_print()
  - evaluate_binary()
  - evaluate_unary()
  - function_declaration()
  - evaluate_condition()
  - evaluate_do()
  - Evaluate()
  - evaluate_assignment()
  - evaluate_if()
  - evaluate_while()
  - evaluate_statement()
*/
func evaluate_statement(v map[string]interface{}) (any, error) {
	//fmt.Println(v)
	t := v["type"]
	//fmt.Println(t)
	if t != nil {
		if t == "control" {
			return v["value"], nil
		} else if t == "number" {
			//fmt.Println("number")
			i, _ := strconv.ParseFloat(v["value"].(string), 32)
			return i, nil
		} else if t == "string" {
			return strings.Trim(v["value"].(string), string(v["value"].(string)[0])), nil

		} else if t == "print" {
			//fmt.Println("print")
			return evaluate_print(v["expression"].([]map[string]interface{}), v["end"].(map[string]interface{}))
		} else if t == "binary" {
			return evaluate_binary(v["operator"].(tokenizer.Token).Value, v["left"].(map[string]interface{}), v["right"].(map[string]interface{}))
		} else if t == "unary" {
			return evaluate_unary(v["operator"].(tokenizer.Token).Value, v["expression"].(map[string]interface{}))

		} else if t == "block" {
			return Evaluate(v["statements"].([]map[string]interface{}))
		} else if t == "assignment" {
			return evaluate_assignment(v["name"].(string), v["expression"].(map[string]interface{}))
		} else if t == "identifier" {
			name := v["name"]
			return env[name.(string)], nil
		} else if t == "if" {
			return evaluate_if(v["condition"].(map[string]interface{}), v["then"].(map[string]interface{}), v["else"])
		} else if t == "while" {
			return evaluate_while(v["condition"].(map[string]interface{}), v["do"].(map[string]interface{}))
		} else if t == "do_statement" {
			return evaluate_do(v["condition"].(map[string]interface{}), v["do"].(map[string]interface{}))
		} else if t == "comparison" {
			evaluated_x, err_x := evaluate_statement(v["left"].(map[string]interface{}))
			evaluated_y, err_y := evaluate_statement(v["right"].(map[string]interface{}))
			if err_x != nil {
				return nil, errors.New("invalid lhs for comparison")
			}
			if err_y != nil {
				return nil, errors.New("invalid type for comparison")
			}
			if reflect.TypeOf(evaluated_x) == reflect.TypeOf(1.0) && reflect.TypeOf(evaluated_y) == reflect.TypeOf(1.0) { // numbers
				return binary_comparisons[v["operator"].(tokenizer.Token).Value](evaluated_x.(float64), evaluated_y.(float64)), nil
			} else if reflect.TypeOf(evaluated_x) == reflect.TypeOf("1") && reflect.TypeOf(evaluated_y) == reflect.TypeOf("1") { // strings
				return binary_comparisons_string[v["operator"].(tokenizer.Token).Value](evaluated_x.(string), evaluated_y.(string)), nil
			} else {
				return false, errors.New("invalid type for comparison")
			}
		} else if t == "function" {
			var parameters []string
			for _, p := range v["parameters"].([]map[string]interface{}) {
				parameters = append(parameters, p["name"].(string))
			}
			fmt.Println(parameters)
			return function_declaration(v["name"].(string), parameters, v["body"].(map[string]interface{})), nil
		}
	}
	//number
	return -1, nil
}
