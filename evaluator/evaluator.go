package evaluator

import (
	"fmt"
	"interperter/tokenizer"
	"os"
	"reflect"
	"strconv"
	"strings"
)

var binary_operations = map[string]func(float64, float64) float64{
	"+": func(x, y float64) float64 { return x + y },
	"-": func(x, y float64) float64 { return x - y },
	"*": func(x, y float64) float64 { return x * y },
	"/": func(x, y float64) float64 { return x / y },
}

var binary_comparisons = map[string]func(float64, float64) bool{
	"<":  func(x, y float64) bool { return x < y },
	"<=": func(x, y float64) bool { return x <= y },
	">":  func(x, y float64) bool { return x > y },
	">=": func(x, y float64) bool { return x >= y },
	"==": func(x, y float64) bool { return x == y },
	"!=": func(x, y float64) bool { return x != y },
}
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
var env = map[string]interface{}{}

func Evalute(node []map[string]interface{}) any {
	var index int = 0
	current_node := node[index]
	var result any
	//fmt.Println("node:", current_node)
	if current_node["type"] == "program" {
		index += 1
		current_node = node[index]
		//fmt.Println("new node:", current_node)
		current_array := current_node["statements"].([]map[string]interface{})
		for _, v := range current_array {
			result = evalute_statement(v)
		}

	} else {
		for _, v := range node {
			result = evalute_statement(v)
		}
	}
	return result
}
func evalute_print(list []map[string]interface{}, end map[string]interface{}) any {
	var result []interface{}
	delimter := "\n"
	if end["string"] != "\n" {
		delimter = evalute_statement(end).(string)
	}
	for i, x := range list {
		result = append(result, evalute_statement(x))
		//fmt.Println("result", result)
		fmt.Printf("%v%v", result[i], delimter)
	}
	result = append(result, delimter)
	return result
}
func evalute_binary(op string, x, y map[string]interface{}) float64 {

	evaluted_x, err_x := evalute_statement(x).(float64)
	evaluted_y, err_y := evalute_statement(y).(float64)
	if !err_x || !err_y {
		fmt.Println("invalid types")
		os.Exit(5)
	}
	return binary_operations[op](evaluted_x, evaluted_y)

}

func evalute_unary(op string, x map[string]interface{}) float64 {

	evaluted_x, err_x := evalute_statement(x).(float64)
	if !err_x {
		fmt.Println("invalid types")
		os.Exit(5)
	}
	return unary_operations[op](evaluted_x)

}
func evalute_assignment(name string, x map[string]interface{}) any {
	//fmt.Println(x)
	evaluated_x := evalute_statement(x)
	env[name] = evaluated_x
	return env[name]
}
func evalute_if(condition, then map[string]interface{}, else_statement any) any {
	con := evalate_condition(condition)
	if con {
		return evalute_statement(then)
	} else {
		if else_statement != nil {
			return evalute_statement(else_statement.(map[string]interface{}))
		} else {
			return nil
		}
	}

}

func evalute_while(condition, then_statement map[string]interface{}) any {
	var result interface{} = nil
	for evalate_condition(condition) {
		result = evalute_statement(then_statement)
	}
	return result
}
func evalute_do(condition, then_statement map[string]interface{}) any {
	var result interface{} = evalute_statement(then_statement)
	for evalate_condition(condition) {
		result = evalute_statement(then_statement)
	}
	return result
}
func evalate_condition(condition map[string]interface{}) bool {
	if condition["type"] == "comparison" {
		return evalute_statement(condition).(bool)
	}
	return evalute_statement(condition).(float64) > 0

}
func function_declaration(name string, parameters []string, body map[string]interface{}) any {
	return 1

}
func evalute_statement(v map[string]interface{}) any {
	//fmt.Println(v)
	t := v["type"]
	//fmt.Println(t)
	if t != nil {
		if t == "print" {
			//fmt.Println("print")
			return evalute_print(v["expression"].([]map[string]interface{}), v["end"].(map[string]interface{}))
		} else if t == "binary" {
			return evalute_binary(v["operator"].(tokenizer.Token).Value, v["left"].(map[string]interface{}), v["right"].(map[string]interface{}))
		} else if t == "unary" {
			return evalute_unary(v["operator"].(tokenizer.Token).Value, v["expression"].(map[string]interface{}))

		} else if t == "block" {
			//fmt.Println("block")
			return Evalute(v["statements"].([]map[string]interface{}))
		} else if t == "assignment" {
			return evalute_assignment(v["name"].(string), v["expression"].(map[string]interface{}))
		} else if t == "identifier" {
			name := v["name"]
			//fmt.Println("env", env[name.(string)])
			return env[name.(string)]
		} else if t == "if" {
			return evalute_if(v["condition"].(map[string]interface{}), v["then"].(map[string]interface{}), v["else"])
		} else if t == "while" {
			return evalute_while(v["condition"].(map[string]interface{}), v["do"].(map[string]interface{}))
		} else if t == "do_statement" {
			return evalute_do(v["condition"].(map[string]interface{}), v["do"].(map[string]interface{}))
		} else if t == "comparison" {
			evaluted_x := evalute_statement(v["left"].(map[string]interface{}))
			evaluted_y := evalute_statement(v["right"].(map[string]interface{}))
			if reflect.TypeOf(evaluted_x) == reflect.TypeOf(1.0) && reflect.TypeOf(evaluted_y) == reflect.TypeOf(1.0) { // numbers
				return binary_comparisons[v["operator"].(tokenizer.Token).Value](evaluted_x.(float64), evaluted_y.(float64))
			} else if reflect.TypeOf(evaluted_x) == reflect.TypeOf("1") && reflect.TypeOf(evaluted_y) == reflect.TypeOf("1") { // strings
				return binary_comparisons_string[v["operator"].(tokenizer.Token).Value](evaluted_x.(string), evaluted_y.(string))
			} else {
				fmt.Println("invalid types")
				os.Exit(5)
			}
		} else if t == "function" {
			var parameters []string
			for _, p := range v["parameters"].([]map[string]interface{}) {
				parameters = append(parameters, p["name"].(string))
			}
			fmt.Println(parameters)
			return function_declaration(v["name"].(string), parameters, v["body"].(map[string]interface{}))
		}
	}
	//number
	if v["number"] != nil {
		//fmt.Println("number")
		i, _ := strconv.ParseFloat(v["number"].(string), 32)
		return i
	}
	if v["string"] != nil {
		return strings.Trim(v["string"].(string), string(v["string"].(string)[0]))

	}
	if v["control"] != nil {
		return v["control"]
	}
	return -1
}
