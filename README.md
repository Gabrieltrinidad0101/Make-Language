
![Logo](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/th5xamgrr6se0x5ro4g6.png)


# Make Language

Make Language is a power full tool and framework that allow you create your own programming language and has a great api that you can add more functionality.

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

## API Reference

#### Var

```
  var a = 1
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

#### if

```
  if(10 == 10){
    print("Ten is equal ten")
  }
```

#### elif

```
  if(10 != 10){
    print("Ten is not equal ten")
  }elif(10 == 10){
    print("Ten is equal ten")
  }
```

#### else

```
  if(10 != 10){
    print("Ten is not equal ten")
  }else{
    print("Ten is equal ten")
  }
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
  }
```

