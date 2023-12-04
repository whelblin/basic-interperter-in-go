package evaluator

import (
	"interpreter/parser"
	"interpreter/tokenizer"
	"testing"

	"github.com/stretchr/testify/assert"
)

func evaluateText(text string) any {
	tokens, _ := tokenizer.Tokenize(text)
	ast, _ := parser.Parse(tokens)
	return first(Evaluate(ast))
}
func Test_binary_operations(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(binary_operations["+"](3.5, 5), 8.5, "should be equal")
	assert.Equal(unary_operations["-"](3.5), -3.5, "should be equal")

}

func Test_print_statements(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(evaluateText("print (2);"), []interface{}{2.0, "\n"}, "should be equal")
	assert.Equal(evaluateText("print (2 + 4);"), []interface{}{6.0, "\n"}, "should be equal")
	assert.Equal(evaluateText("print (-5);"), []interface{}{-5.0, "\n"}, "should be equal")
	assert.Equal(evaluateText(`print ("hello");`), []interface{}{`hello`, "\n"}, "should be equal")
	assert.Equal(evaluateText("print ((2 * 2) + 5);"), []interface{}{9.0, "\n"}, "should be equal")
	assert.Equal(evaluateText("print( 2 + 2 * 5);"), []interface{}{12.0, "\n"}, "should be equal")

}

func Test_print_statements_with_end(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(evaluateText("print (2);"), []interface{}{2.0, "\n"}, "should be equal")
	assert.Equal(evaluateText("print (2 + 4);"), []interface{}{6.0, "\n"}, "should be equal")
	assert.Equal(evaluateText("print (-5);"), []interface{}{-5.0, "\n"}, "should be equal")
	assert.Equal(evaluateText(`print ("hello");`), []interface{}{`hello`, "\n"}, "should be equal")
	assert.Equal(evaluateText("print ((2 * 2) + 5);"), []interface{}{9.0, "\n"}, "should be equal")
	assert.Equal(evaluateText("print( 2 + 2 * 5);"), []interface{}{12.0, "\n"}, "should be equal")

}
func Test_idenifier_statements(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(evaluateText("x = 5; print (x);"), []interface{}{5.0, "\n"}, "should be equal")
	assert.Equal(evaluateText("x = 5; print (x + 4);"), []interface{}{9.0, "\n"}, "should be equal")
	assert.Equal(evaluateText("x = -5; print (x);"), []interface{}{-5.0, "\n"}, "should be equal")
	assert.Equal(evaluateText("x = -5; print (x + 2);"), []interface{}{-3.0, "\n"}, "should be equal")
	assert.Equal(evaluateText("x = 5; print (x); x = 4;"), 4.0, "should be equal")
	assert.Equal(evaluateText(`x = 5; print (x); x = "hello";`), `hello`, "should be equal")
	assert.Equal(evaluateText(`x = 5; y = x + 3; print (x + y);`), []interface{}{13.0, "\n"}, "should be equal")
	assert.Equal(evaluateText(`x = 5; y = x + 3; print (x + y -5);`), []interface{}{8.0, "\n"}, "should be equal")
}
func Test_block_statement(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(evaluateText("{x = 5; print (x);}"), []interface{}{5.0, "\n"}, "should be equal")
	assert.Equal(evaluateText("{x = 5; print (x + 4);}"), []interface{}{9.0, "\n"}, "should be equal")
	assert.Equal(evaluateText("{x = -5; print (x);}"), []interface{}{-5.0, "\n"}, "should be equal")
	assert.Equal(evaluateText("{x = -5; print (x + 2);}"), []interface{}{-3.0, "\n"}, "should be equal")
	assert.Equal(evaluateText("{x = 5; print (x); x = 4;}"), 4.0, "should be equal")
	assert.Equal(evaluateText(`{x = 5; print (x); x = "hello";}`), `hello`, "should be equal")
	assert.Equal(evaluateText(`{x = 5; y = x + 3; print (x + y);}`), []interface{}{13.0, "\n"}, "should be equal")
	assert.Equal(evaluateText(`{x = 5; y = x + 3; print (x + y -5);}`), []interface{}{8.0, "\n"}, "should be equal")
}

func Test_if_statement(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(evaluateText("if(1){print (1);}"), []interface{}{1.0, "\n"}, "should be equal")
	assert.Equal(evaluateText("if(5){print(1);}"), []interface{}{1.0, "\n"}, "should be equal")
	assert.Equal(evaluateText("if(0){print (1);}"), nil, "should be equal")
	assert.Equal(evaluateText("if(-1){print (1);}"), nil, "should be equal")
	assert.Equal(evaluateText("x = 1; if(x){print (x);}"), []interface{}{1.0, "\n"}, "should be equal")
	assert.Equal(evaluateText("if(1){print (1);}else{print (2);}"), []interface{}{1.0, "\n"}, "should be equal")
	assert.Equal(evaluateText("if(0){print (1);}else{print (2);}"), []interface{}{2.0, "\n"}, "should be equal")
	assert.Equal(evaluateText("if(1 < 3){print (1);}else{print (2);}"), []interface{}{1.0, "\n"}, "should be equal")
	assert.Equal(evaluateText("if(3 > 5){print (1);}else{print (2);}"), []interface{}{2.0, "\n"}, "should be equal")
	assert.Equal(evaluateText("if(2 <= 2){print (1);}else{print (2);}"), []interface{}{1.0, "\n"}, "should be equal")
	assert.Equal(evaluateText("if(2 >= 1){print (1);}else{print (2);}"), []interface{}{1.0, "\n"}, "should be equal")
	assert.Equal(evaluateText("if(2 == 2){print (1);}else{print (2);}"), []interface{}{1.0, "\n"}, "should be equal")
	assert.Equal(evaluateText("if(2 != 2){print (1);}else{print (2);}"), []interface{}{2.0, "\n"}, "should be equal")
	assert.Equal(evaluateText(`if("a" != "a"){print (1);}else{print (2);}`), []interface{}{2.0, "\n"}, "should be equal")
	assert.Equal(evaluateText(`if("a" == "a"){print (1);}else{print (2);}`), []interface{}{1.0, "\n"}, "should be equal")

	assert.Equal(evaluateText(`if(1){if(0){
											print (1);
										}else{
											print (2);
											}
										}else{
											if(1){
											print (3);
											}else{
												print (4);
											}
										}`), []interface{}{2.0, "\n"}, "should be equal")
	assert.Equal(evaluateText(`if(1){if(1){
											print (1);
										}else{
											print (2);
											}
										}else{
											if(1){
											print (3);
											}else{
												print (4);
											}
										}`), []interface{}{1.0, "\n"}, "should be equal")
	assert.Equal(evaluateText(`if(0){if(0){
											print (1);
										}else{
											print (2);
											}
										}else{
											if(1){
											print (3);
											}else{
												print (4);
											}
										}`), []interface{}{3.0, "\n"}, "should be equal")
	assert.Equal(evaluateText(`if(0){if(0){
											print (1);
										}else{
											print (2);
											}
										}else{
											if(0){
											print (3);
											}else{
												print (4);
											}
										}`), []interface{}{4.0, "\n"}, "should be equal")

}

func Test_while_loop(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(evaluateText(`x = 10; while (x){x = x -1;}print (x);`), []interface{}{0.0, "\n"}, "should be equal")
	assert.Equal(evaluateText(`x = 10; while (x > 2){x = x -1;}print (x);`), []interface{}{2.0, "\n"}, "should be equal")

}
