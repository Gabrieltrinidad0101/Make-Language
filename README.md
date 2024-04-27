
![Logo](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/th5xamgrr6se0x5ro4g6.png)


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

## API Reference

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

