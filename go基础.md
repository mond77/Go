# Go语言基础

函数调用参数均为值传递

Go语言中，在函数调用时，引用类型（slice、map、interface、channel）都默认使用引用传递。

map == *hmap

chan == *hchan

`slice`是一种结构体+元素指针的混合类型

```go
type slice struct {
array unsafe.Pointer
len int
cap int
}
```

## 数据类型

go语言中的 int 的大小是和操作系统位数相关的，如果是32位操作系统，int 类型的大小就是4字节。如果是64位操作系统，int 类型的大小就是8个字节



## 关键字

### append

如果原来slice capacity足够大的情况下，append()函数会创建一个新的slice，它与old slice共享底层内存

创建原理：newSlice = oldSlice[:1+len(x)]

用old slice给new slice进行赋值的方式进行创建，会共享内存。并返回这个new slice。

因此为了保险，我们通常将append返回的内容赋值给原来的slice： x = appen(x,…)
如果原来的slice没有足够的容量添加内容，则创建一个新的slice，这个slice是copy的old slice。不与old slice共享内存
————————————————

### copy（与=的区别）

`=` 赋值拷贝，会将原来slice的地址拷贝，新旧slice共享内存。

`copy(new, old)` 函数copy只会将slice内容进行拷贝

### new

表达式new(T)将创建一个T类型的匿名变量，初始化为T类型的零值，然后返回变量地址，返回的指针类型为 *T 。

#### 特性

每次调用 new 函数都会返回唯一的地址变量: 

```go
package main

import "fmt"

func main() {

    p1 := new(int)
    q1 := new(int)

    fmt.Println(p1, q1, p1 == q1) // 0xc00001c0b8 0xc00001c0c0 false
}
```

但是也会有例外，当定义一个空 struct 时，通过 new 创建一个变量时，返回的地址是相同的。

```go
package main

import "fmt"


type Point struct {}

func main() {
    p2 := new(Point)
    q2 := new(Point)
    fmt.Println(p2, q2, p2 == q2) // &{} &{} true
}
```

## 函数

### defer与return



我们知道go中的函数返回值有匿名和有名两种。

对于匿名的可以理解成执行return 语句的时候，分成两步，第一步需要设置一个临时变量s用来接收返回值，第二步将临时变量s返回。

对于有名的可以理解成执行return的时候，直接将变量返回。

而我们知道，所有的defer都将在真正的return 变量之前运行，所以对于上面两种情况，defer对于返回值的影响也有两种：

对于匿名的：第一步设置临时变量保存返回值，第二步按照defer的执行步骤执行defer语句，如果其中有对变量的修改，将不会影响s变

量的值。

对于有名的：第一步先执行defer，对变量进行修改，第二步，返回被修改的返回值。



## goroutine

### goroutine和线程的区别

OS线程（操作系统线程）一本都有固定的栈内存（通常为2MB），一个 goroutine 的栈在其生命周期开始时只有很小的栈（典型情况下2KB），goroutine 的栈不是固定的，他可以按需增大和缩小，grorutine的栈大小限制可以达到1GB，但极少情况下会到1GB。所以在Go语言中一次创建十万左右的 grorutine 也是可以的。

### GOMAXPROCS

Go语言中可以通过`runtime.GOMAXPROCS()`函数设置当前程序并发时占用的CPU逻辑核心数。

Go1.5版本之前，默认使用的是单核心执行。Go1.5版本之后，默认使用全部的CPU逻辑核心数。

## channel

### 有缓冲与无缓冲的区别

无缓冲要求读写同时进行，否则一方阻塞等待；

有缓冲则只会在缓冲区空时读数据、满时写数据阻塞；

无缓存channel适用于数据要求同步的场景，而有缓存channel适用于无数据同步的场景。

### for-range遍历

channel 支持 for-range的方式进行[遍历](https://so.csdn.net/so/search?q=遍历&spm=1001.2101.3001.7020)，请注意两个细节

1）在遍历时，如果channel没有关闭，则会出现deadlock 的错误

2）在遍历时，如果channel已经关闭，则会正常遍历数据，遍历完后，就会退出遍历

### 关闭channel

```go
close(ch)
```

如果c已经关闭（c中所有值都被接收）， `x, ok := <- c`， 读取ok将会得到false。

## 错误捕获

### panic

Golang中当程序发生致命异常时（比如数组下标越界，注意这里的异常并不是error），Golang程序会panic（运行时恐慌）。当程序发生panic时，程序会执行**当前栈中**的defer 函数列表。然后打印引发panic的具体信息，最后进程退出。

```go
func RecoverError() {
	if err := recover(); err != nil {
		//输出panic信息
		fmt.Println(err)
		//输出堆栈信息
		fmt.Println(string(debug.Stack()))
	}
}
func main() {
	// 当使用defer 时，将会在程序内方法结算后，
	// 依照后进先出的方法执行defer内方法
	// 此时就能保证 捕获程序一定能捕获到错误
	defer RecoverError()
	for _, c := range []string{"1", "2"} {
		atoi, err := strconv.Atoi(c)
		Panic(err)
		fmt.Println(atoi)
	}
}
```

 defer函数要在panic之前执行（压入栈中）

### recover()

```go
package main

import (
    "fmt"
    "time"
)
func calcRem(i int) (res int) {
    defer func() {
        if err := recover(); err != nil {
            fmt.Printf("error: %s\n", err)
            res = 999 //干扰输出
        }
    }()
    res = 10 / (10 % i) //当余数取余为0时，res为0
    fmt.Printf("10 / (10 % %d) = %d\n", i, res)
    return
}
func main() {
    resCH := make(chan int, 10)
    defer close(resCH)
    for i := 1; i &lt;= 10; i++ {
        go func(i int) {
            res := calcRem(i)//recover后该goroutine遇到panic，不会跳转执行栈中的defer函数，而是继续执行下面的代码
            resCH &lt;- res
        }(i)
        time.Sleep(time.Microsecond * 100)
    }
    time.Sleep(time.Second)
    fmt.Printf("ch len is: %d\n", len(resCH))
    for i := 1; i &lt;= 10; i++ {
        res := &lt;-resCH
        fmt.Println(res)
    }
}
```



## runtime



## stdLib

### strings

## errors

```go
type error interface {
    Error() string
}
```

```go
package main
 
import (
    "fmt"
    "time"
)
 
// MyError is an error implementation that include a time and message.
type MyError struct {
    When time.Time
    What string
}
 
func (e MyError) Error() string {
    return fmt.Sprintf("%v: %v", e.When, e.What)
}
 
func oops() error {
    return MyError{
        time.Date(1989, 3, 15, 22, 30, 0, 0, time.UTC),
        "the file system has gone away",
    }
}
 
func main() {
    if err := oops(); err != nil {
        fmt.Println(err)
    }
}
```

### rpc

定义远程服务

```go
func (t *T) MethodName(argType T1, replyType *T2) error
```

## Lib（frequent）





## pprof性能优化

https://jvns.ca/blog/2017/09/24/profiling-go-with-pprof/

pprof lets you collect CPU profiles, traces, and heap profiles for your Go programs. The normal way to use pprof seems to be:

1. Set up a webserver for getting Go profiles (with `import _ "net/http/pprof"`)
2. Run `curl localhost:$PORT/debug/pprof/$PROFILE_TYPE` to save a profile
3. Use `go tool pprof` to analyze said profile

You can also generate pprof profiles in your code using the [`pprof` package](https://golang.org/pkg/runtime/pprof/) but I haven’t done that.

## reflect(反射)

s
