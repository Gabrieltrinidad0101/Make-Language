
![Logo](https://github.com/Gabrieltrinidad0101/Make-Language/blob/master/logo.png)


# Make Language

Make Language is a powerful tool and framework that allows you to create your own programming language and has an excellent API which allows you to add more features.

## Features

- var
- string
- number
- const
- if
- elif
- else
- for
- while
- array
- functions
- class
- math operations

Any of these functionality can be modify using a simple conf.json file

## Language Reference

#### Var

```
  var a = 1
  var b = TRUE
  var c = FALSE
  var d = 1 >= 1
  var e = 1 == 1
  var f = 1 != 1
  var g = 1 <= 1
  print(a,b,c,d,e,f,g)
```

#### Const

```
  const a = 1
  a = 2 // error
```

#### String

```
  var a = "hello world"
  print(a.replace("e","u"))
  print(a +  " 123")
  print(a * 2)
  print(a.upper())
```

#### Number

```
  var a = 25
  print(a + 5) // 30
  print(a - 5) // 25
  print(a * 2) // 50
  print(a / 5) // 5
  print(a ~ 2) // 5
  print(100 ~ 2) // 10
```

#### If

```
  if(10 == 10){
    print("Ten is equal ten")
  }
```

#### Elif

```
  if(10 != 10){
    print("Ten is not equal ten")
  }elif(10 == 10){
    print("Ten is equal ten")
  }
```

#### Else

```
  if(10 != 10){
    print("Ten is not equal ten")
  }else{
    print("Ten is equal ten")
  }
```

#### While
 
```
var i = 0
while(i <= 10){
  if(i == 5){
    continue
  } elif(i == 8){
    break
  }
  
  print(i)
}
```

#### For
 
```
for(var i = 0; i < 10; ++i){
  if(i == 5){
    continue
  } elif(i == 8){
    break
  }
  print(i)
}
```
#### Functions

```
func a(a1,a2){
  return a1 + a2
}

print(a(1,2))
```



#### Class

```
  class Test {
    func a(){
      this.b("Hello ")
    }

    func b(hello){
      print(hello + "World")
    }

    func c(){
      this.z = "Hello world"
      print(this.z)
    }

    func d(){
      return 100
    }

    func e(){
      return this
    }
  }

  const test = new Test()
  test.a()
  print(test.b() == 100)
  print(test.e().b() == 100)
```
## API
create a json file inside of language syntax you can only modify the key if you change any key in the json file 
you are going to modify the language syntax

```
  {
    "numbers": {
        "1": "1",
        "2": "2",
        "3": "3",
        "4": "4",
        "5": "5",
        "6": "6",
        "7": "7",
        "8": "8",
        "9": "9",
        "0": "0"
    },
    "compares": {
        "==": "EQE",
        ">=": "GTE",
        "<=": "LTE",
        "!=": "NEQE",
        "++": "PLUS1",
        "--": "MINUS1"
    },
    "language_syntax": {
        "[": "LSQUAREBRACKET",
        "]": "RSQUAREBRACKET",
        "&&": "AND",
        "||": "OR",
        ">": "GT",
        "<": "LT",
        "!": "NEQ",
        ",": "COMMA",
        "=": "EQ",
        "+": "PLUS",
        "-": "MINUS",
        "*": "MUL",
        "/": "DIV",
        "(": "LPAREN",
        ")": "RPAREN",
        "^":  "POW",
        "~": "SQUARE_ROOT",
        "func": "FUNC",
        "while": "WHILE",
        "for": "FOR",
        "if": "IF",
        "elif": "ELIF",
        "else": "ELSE",
        "class": "CLASS",
        ".": "SPOT",
        "break": "BREAK",
        "continue": "CONTINUE",
        "new": "NEW",
        "{": "OPEN_CURLY_BRACE",
        "}": "CLOSING_CURLY_BRACE",
        "var": "VAR",
        "this": "THIS",
        "const": "CONST",
        "return": "RETURN",
        "\"": "STRING",
        "\n": "NEWLINE",
        "\r": "NEWLINE",
        ";": "SEMICOLON"
    },
    "functions": {
        "print": "print"
    },
    "scope": "global"
  }
```

## Custom operator
```
  package main

  import (
    "makeLanguages/src"
    "makeLanguages/src/features/booleans"
    "makeLanguages/src/interprete/interpreteStructs"
    "makeLanguages/src/utils"
  )

  func lessOrGreaterOne(value1 interpreteStructs.IBaseElement, value2 interpreteStructs.IBaseElement) interface{} {
    params := &[]interpreteStructs.IBaseElement{
      value1,
      value2,
    }
    utils.ValidateTypes(params, "Number", "Number")
    number1 := value1.GetValue().(float64)
    number2 := value2.GetValue().(float64)
    boolean := number1+1 == number2 || number1-1 == number2
    return booleans.NewBoolean(boolean)
  }

  func main() {
    makeLanguage := src.NewMakeLanguage("./conf.json", "./main.mkL")
    makeLanguage.AddOperetor("<1>", lessOrGreaterOne)
    makeLanguage.Run()
  }
```

```
  print(1 <1> 2) // true
  print(10 <1> 2) // false
```

## Custom functions

```
  package main

  import (
    "fmt"
    "makeLanguages/src"
    "makeLanguages/src/interprete/interpreteStructs"
    "makeLanguages/src/parser/parserStructs"
  )

  func printLn2(params *[]interpreteStructs.IBaseElement) interface{} {
    fmt.Println((*params)[0].GetValue())
    fmt.Println()
    return parserStructs.NullNode{}
  }

  func main() {
    makeLanguage := src.NewMakeLanguage("./conf.json", "./main.mkL")
    makeLanguage.AddFunction("printLn2", printLn2)
    makeLanguage.Run()
  }
```

```
  printLn2("hello world")
```

## Custom class

```
  package main

  import (
    "fmt"
    "makeLanguages/src"
    "makeLanguages/src/interprete/interpreteStructs"
    "makeLanguages/src/parser/parserStructs"
    "makeLanguages/src/utils"
    "os"
  )

  func makeFile(params *[]interpreteStructs.IBaseElement) interface{} {
    utils.ValidateTypes(params, "String_")
    _, err := os.Create((*params)[0].GetValue().(string))
    if err != nil {
      fmt.Println(err)
    }

    return parserStructs.NullNode{}
  }

  func main() {
    makeLanguage := src.NewMakeLanguage("./conf.json", "./main.mkL")
    methods := map[string]func(params *[]interpreteStructs.IBaseElement) interface{}{}
    methods["create"] = makeFile
    makeLanguage.AddClass("File", methods)
    makeLanguage.Run()
  }
```

```
  const file = new File()
  file.create("./hello.txt")
```


