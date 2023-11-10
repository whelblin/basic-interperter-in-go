package evaluator

import (
	"fmt"
	"testing"
)

func Test_binary_operations(t *testing.T) {
	fmt.Printf("%v\n", (binary_operations["+"](3.5, 5)))
	fmt.Printf("%v\n", (unary_operations["-"](3.5)))

}
