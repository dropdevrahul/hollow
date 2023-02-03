## Hollow
A simple programming language written in Golang

Stack based with simple instructions, result will always be on top of stack you can use dump(`.`) command to display the result to stdout

```
100 => pushes the given int to stack 

+ => add given number to whatever is at top of stack and push result to stack e.g 100 + 50 pushes 150 to stack(orignal 100 will be replaced)

- => sub two numbers and push result to stack e.g 100 - 50 pushes result to stack

. => prints the current top of stack to stdout
== => compares a number with top of stack e.g  == 100 will return 0 if current number      is not equal to top of stack 1 otherwise

Please  note for all +, == , - operations only one operand is required i.e the right hand side of the operator but there should be a valid value on stack for this

15

+ 10 .

will result in 15 , however
+ 10 . will result in undefined behaviour

```

For an input file:

```
64 + 56 . 
65 + 100 .
200 - 100 .
65 == 100 .
100 - 200 .
100 - 200 == -100 .
```

Output will be:
```
120
165
100
0
-100
1
```

## Compiling and running
There is a sample `example/program.hollow` to check otherwise you can use ur own file with any name with syntax as shown in the file or above examples:
```
go run main.go -o myprogram example/program.hollow

./myprogram

```

For help

```
go run main.go -h
```





