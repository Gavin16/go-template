package main

import (
	"C"
	"fmt"
	"os"
	"sync"
)

// ** 习惯函数是"一等公民"
// Go语言中函数像变量一样 支持
// (1) 正常创建
// (2) 在函数内创建
// (3) 使用函数来定义类型
// (4) 将函数存储到变量中
// (5) 函数作为参数传入函数
// (6) 作为返回值从函数返回

// ** 使用defer 让函数更简洁,更健壮
// 无论是执行到函数尾部返回 还是在某个错误处理分支显示调用return 返回
// 又或者出现panic, 已经存储到 deferred 函数栈中的函数都会被调度执行
// 因此，deferred 函数是 一个在任何情况下都可以为函数进行收尾工作的好场合

func main() {
	// MyAdderFunc 定义成了 func(int,int) int 的别名
	// 此外 MyAdderFunc 的底层类型和 MyAdder 是一样的
	// 因此可以将 MyAdder 显示转换为 HandlerFunc 类型
	// 由于MyAdderFunc 实现了 BinaryAdder 的接口
	// 因此可以直接将结果(函数)赋值给 BinaryAdder 类型
	var adder BinaryAdder = MyAdderFunc(MyAdder)
	fmt.Println(adder.Add(5, 6))

	// 相当于将times 拆分成两个单参数的函数
	timesTwo := partialTimes(2)
	fmt.Println(timesTwo(5))

	// 函子的应用
	// 函子是一个容器类型, 它实现了一个接受函数类型的参数，并在容器中每个元素上应用那个函数，得到一个新函子
	// 原函子容器内部的元素值不受影响
	intSlice := []int{1, 2, 3, 4, 5}
	fmt.Printf("init a functor from int slice:%v\n", intSlice)
	f := NewIntSliceFunctor(intSlice)
	mapperFunc1 := func(i int) int {
		return i + 10
	}
	mapped1 := f.fMap(mapperFunc1)
	fmt.Printf("mapped functor1: %+v\n", mapped1)

	fmt.Printf("origin functor:%+v\n", f)

	// 使用defer 进行收尾处理
	// 其中有几个关键问题
	// 1. 明确哪些函数可以直接放在defer 后面使用,哪些需要自定义(匿名)函数进行包装
	//    -- 对于内置函数: close,copy,delete,print,recover 可以直接使用
	//                   append, cap, len, make, new 不可以直接作为 deferred 的函数
	// 2. 把握defer 关键字后表达式的求值时机
	//    -- defer 关键字后面的表达式是在将 deferred 函数注册到 deferred 函数栈的时候进行求值的
	//       所以以下 debugLogPrint 日志打印功能才能实现前后切面的日志打印
	_ = writeToFile()
	fmt.Println("--------------")
	modifyReturnValue()
	debugLogPrint()
	deferEvaluate()
}

// BinaryAdder --- 函数显示转换
type BinaryAdder interface {
	Add(int, int) int
}

type MyAdderFunc func(int, int) int

// Add 通过别名方式，实际上是让某个函数实现了一个方法
// 因此该函数能够赋值给定义该方法的接口
func (f MyAdderFunc) Add(x, y int) int {
	return f(x, y)
}

func MyAdder(x, y int) int {
	return x + y
}

// --- 函数式编程
func times(x, y int) int {
	return x * y
}

func partialTimes(x int) func(int) int {
	// 返回闭包
	return func(y int) int {
		return times(x, y)
	}
}

// IntSliceFunctor --- 函子做数据做映射
type IntSliceFunctor interface {
	fMap(fn func(int) int) IntSliceFunctor
}

type intSliceFunctorImpl struct {
	ints []int
}

func (isf intSliceFunctorImpl) fMap(fn func(int) int) IntSliceFunctor {
	newInts := make([]int, len(isf.ints))
	for i, v := range isf.ints {
		retInt := fn(v)
		newInts[i] = retInt
	}
	return intSliceFunctorImpl{ints: newInts}
}

func NewIntSliceFunctor(slice []int) IntSliceFunctor {
	return intSliceFunctorImpl{ints: slice}
}

// defer 函数的使用
var mu sync.Mutex

func writeToFile() error {
	mu.Lock()
	defer mu.Unlock()

	err := os.WriteFile("test_write.txt", []byte("..test defer usage .."), 0666)
	if err != nil {
		fmt.Println("write to file err:", err)
		return err
	}
	return nil
}

func modifyReturnValue() {
	foo := func(a, b int) (x int, y int) {
		defer func() {
			x *= 5
			y *= 10
		}()

		return a + 5, b + 6
	}
	x, y := foo(1, 2)
	fmt.Println(x, y)
}

func debugLogPrint() {
	leave := func(s string) {
		fmt.Println("leave:", s)
	}

	trace := func(s string) string {
		fmt.Println("entering:", s)
		return s
	}
	// 函数A 执行前后打印日志
	funcA := func() int {
		defer leave(trace("a"))
		fmt.Println("funcA")
		return 1
	}
	i := funcA()
	fmt.Println("funcA return: ", i)
}

// 测试defer 函数的求值时机和求值方式
func deferEvaluate() {
	// 测试求值时机
	fmt.Println("-------------")
	func1 := func() {
		for i := 0; i < 3; i++ {
			defer fmt.Println(i)
		}
	}
	func1()
	func2 := func() {
		for i := 0; i < 3; i++ {
			defer func() {
				fmt.Println(i)
			}()
		}
	}
	func2()

	fmt.Println("-------------")
	// 测试求值传参方式: 传参传的为值, 不是引用
	func3 := func() {
		sl := []int{1, 2, 3}
		defer func(a []int) {
			fmt.Println(a)
		}(sl)
		sl = []int{3, 2, 1}
		_ = sl
	}
	func3()

	func4 := func() {
		sl := []int{1, 2, 3}
		defer func(a *[]int) {
			fmt.Println(*a)
		}(&sl)
		sl = []int{3, 2, 1}
		_ = sl
	}
	func4()
}
