package main

import (
	"fmt"
	"testing"
)

func TestFunc(t *testing.T) {
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

	// **理解方法的本质&选择正确的receiver 类型
	essenceOfMethod()

	// **方法集合决定接口实现
	// 考虑接口中嵌入接口, 结构体中嵌入结构 以及 结构体中嵌入结构体 三种场景
	testMethodSetRules()
	testEmbedMethodOrder()
	testMethodSetExtends()

	// **使用变长参数函数
	testVariadicFuncArgs()
}
