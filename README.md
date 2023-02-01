## Hollow
A simple programming language written in Golang

Stack based with simple instructions, result will always be on top of stack you can use dump(`.`) command to display the result to stdout

* + : 100 + 56  
* - : 100 - 200

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





