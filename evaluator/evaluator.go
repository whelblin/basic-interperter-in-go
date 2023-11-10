package evaluator

import (
	"fmt"
	"interperter/tokenizer"
	"os"
	"strconv"
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
var unary_operations = map[string]func(float32) float32{
	"-": func(f float32) float32 { return -f },
}
var env = map[string]interface{}{}

func Evalute(node []map[string]interface{}) {
	var index int = 0
	current_node := node[index]
	//fmt.Println("node:", current_node)
	if current_node["type"] == "program" {
		index += 1
		current_node = node[index]
		//fmt.Println("new node:", current_node)
		current_array := current_node["statements"].([]map[string]interface{})
		for _, v := range current_array {
			evalute_statement(v)
		}

	} else {
		for _, v := range node {
			evalute_statement(v)
		}
	}
}
func evalute_print(x map[string]interface{}) bool {
	result := evalute_statement(x)
	fmt.Printf("%v\n", result)
	return true
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
func evalute_assignment(name string, x map[string]interface{}) bool {
	//fmt.Println(x)
	evaluated_x := evalute_statement(x)
	env[name] = evaluated_x
	return true
}
func evalute_if(condition, then map[string]interface{}, else_statement any) any {
	con := evalate_condition(condition)
	//fmt.Println("Con", con)
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
func evalate_condition(condition map[string]interface{}) bool {
	if condition["type"] == "number" || condition["type"] == "identifier" {
		result := evalute_statement(condition)
		if result.(float64) > 0 {
			return true
		} else {
			return false
		}

	} else {
		return evalute_statement(condition).(bool)
	}
}
func evalute_statement(v map[string]interface{}) any {
	//fmt.Println(v)
	t := v["type"]
	//fmt.Println(t)
	if t != nil {
		if t == "print" {
			//fmt.Println("print")
			return evalute_print(v["expression"].(map[string]interface{}))
		} else if t == "binary" {
			return evalute_binary(v["operator"].(tokenizer.Token).Value, v["left"].(map[string]interface{}), v["right"].(map[string]interface{}))
		} else if t == "block" {
			//fmt.Println("block")
			Evalute(v["statements"].([]map[string]interface{}))
			return 1
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
		} else if t == "comparison" {
			evaluted_x, err_x := evalute_statement(v["left"].(map[string]interface{})).(float64)
			evaluted_y, err_y := evalute_statement(v["right"].(map[string]interface{})).(float64)
			if !err_x || !err_y {
				fmt.Println("invalid types")
				os.Exit(5)
			}
			return binary_comparisons[v["operator"].(tokenizer.Token).Value](evaluted_x, evaluted_y)
		}
	}
	//number
	if v["number"] != nil {
		//fmt.Println("number")
		i, _ := strconv.ParseFloat(v["number"].(string), 32)
		return i
	}
	if v["string"] != nil {
		return v["string"]
	}
	return -1
}
