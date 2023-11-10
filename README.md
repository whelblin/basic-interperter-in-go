# basic-interperter-in-go
## Language Features
### print values
currently you can only print one value for each print statement. Each print statements ends in a endline
```
print 5;
print "hello";
print 1 < 3;
print "a" < "b";
```
### variables
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
### if statement
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
### while statements
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
