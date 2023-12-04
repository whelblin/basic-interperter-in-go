# basic-interpreter-in-go

| [How to install and run the code](#Install-and-How-To-Run)
| [Conversion from the python to go ](#conversion-from-python)
| [Added features from the python version](#added-features-from-the-python-version)
| [Language Features](#language-features)
|
## Install and How To Run
### Install go
go to https://go.dev/doc/install to get the instructions on how to install go
### Run the code

#### With No Arguments
- no arguments runs the RELP
    - The RELP can be exited by typing exit
```
go run .
```
#### With Arguments
- you can also give it a file to run instead
```
go run . test.t
```
### Build the code
```
go build ./...
```
### Run the Test Code
```
go test ./...
```
## Conversion from Python
### Strongly typed
Go is a strongly typed language so it required that each function have a return type and all variables also be assigned a type
Go does have the interface type, which accepts any type, but requires a binding when it its needs to be a specific type
I also had to create a structure called Token to handle the tokens because arrays must be all the same value in go
### Error Handling
Go does not have the same ability to catch errors like python. Instead it uses a testing library to return error values within the parser and tokenizer. Every function returns a error value. It ca either be nil or an error. These error messages are passed up the stack until the Evaluate function returns it
## Added Features From the Python Version
### Print Function
The print function can take multiple values separated by a comma. It delimiter can also be changed with the end= variable
```
print("The value of 1+2 is: ", 1+2, end = ""); print("");
```
```
Output>>The value of 1+2 is: 3
```
### Do While Loops
do whiles are also possible within this implementation of the interpreter.
```
x = 10;
print("The value of x is: ", end = "");
do{
    print(x, end=" ");
    x = x-1;
}while(x > 0);
print("");
```
```
Output>>The value of x is: 10 9 8 7 6 5 4 3 2 1 
```
### String Comparison
Strings can also be compared within the expressions
You can use <, <=, >, >=, !=, and ==
```
x = "hello";
y = "Hello";
print(x == y);

char = "a";
if (char < "b"){
    print("a is less than b");
}
else{
    print("a is greater than b");
}

```
```
Output>>false
      >>a is less than b
```
## Language Features
### Print Values
print statements must be wrapped in () and can print multiple values separated by a comma. Each print statements ends in a newline by default
```
print (5);
print("hello");
print (1 < 3);
print ("a" < "b");
print("hello", "world");
print("The value of 1+2 is: ", 1+2);
```
You can change the delimiter by adding an end= argument to the print function
```
print("The value of 1+2 is: ", 1+2, end = ""); print("");
```

```
Output>>The value of 1+2 is: 3
```
### Variables
You can store numbers and strings. The numbers can be integers or floating points
```
x = 12;
y = "hello world";
z = 1.34;
```
True values are considered numbers with a value greater than 1
```
x = 1; // true
x = 0; // false
x = 3; //true
```
### If Statement
if statements condition must have () and the body must have {}
```
if (1 < 3){
    print 1;
}
```
else statements also exist
```
if (1 < 3){
    print 1;
}
else {
    print 3;
}
```
since number greater than 1 are true, you can use numbers as the condition
```
if (5){
    print 5;
}
```
### While Statements
while loops are also allowed. The condition must be in () and be followed by {}
```
i = 10;
while (i > 0){
    print i;
    i = i - 1;
}
```
This can be done easier with just using a number
```
x = 10;
while(x){
    print x;
    x = x - 1;
}
```
### Do While Statements
do while loops are also allowed. Conditions must be within () and the body must be contained by {}
```
i = 10;
do{
    print i;
    i = i - 1;
} while(i > 0);
```
