package main

import (
	"C"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
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

// T -- 方法本质测试
// 一个以方法所绑定类型实例作为第一个参数的普通函数
type T struct {
	a int
}

func (t *T) get() int {
	return t.a
}
func (t *T) set(a int) {
	t.a = a
}
func get(t *T) int {
	return t.a
}
func set(t *T, a int) {
	t.a = a
}
func essenceOfMethod() {
	fmt.Println("--------------")
	tt := T{1}
	a := tt.get()
	fmt.Println(a)
	a = get(&tt)
	fmt.Println(a)
	tt.set(3)
	fmt.Println(tt)
	set(&tt, 5)
	fmt.Println(tt)
	// 将 func (t *T) get() int 赋值给 f1
	// 结果 f1 变成了 func get(t *T) int
	f1 := (*T).get
	res := f1(&tt)
	fmt.Println(res)
}

// DumpMethodSet
// -- 方法集合决定接口实现
// 打印接口包含的方法集合
func DumpMethodSet(i interface{}) {
	v := reflect.TypeOf(i)
	elemType := v.Elem()
	num := elemType.NumMethod()
	if num == 0 {
		fmt.Printf("%s's method set is empty\n", elemType)
		return
	}
	fmt.Printf("%s's method set is:\n", elemType)
	for j := 0; j < num; j++ {
		fmt.Println("-", elemType.Method(j).Name)
	}
	fmt.Printf("\n")
}

type Interface interface {
	M1()
	M2()
}
type TT struct{}

func (t TT) M1()  {}
func (t *TT) M2() {}

// 对于非接口类型的自定义类型T, 起方法集合有所有receiver 为T类型的方法组成
// 而类型为 *T的方法集合则包含 所有receiver为T和*T类型的方法
func testMethodSetRules() {
	var tt TT
	var ptt *TT
	DumpMethodSet(&tt)
	DumpMethodSet(&ptt)
	DumpMethodSet((*Interface)(nil))

	fmt.Println("-------------")
	DumpMethodSet((*io.Writer)(nil))
	DumpMethodSet((*io.Reader)(nil))
	DumpMethodSet((*io.Closer)(nil))
	DumpMethodSet((*io.ReadWriter)(nil))
	DumpMethodSet((*io.ReadWriteCloser)(nil))
}

// I1 结构体类型嵌入接口类型的同时 也实现了接口类型的方法
// 若接口体类型 自己也实现了接口接口类型中的某个方法，在调用时优先使用自己实现的方法
type I1 interface {
	M4()
	M5()
	M6()
}
type T2 struct {
	I1
}

func (t *T2) M4() {
	fmt.Println("T2's Method4")
}

type S struct{}

func (s *S) M4() {
	fmt.Println("S's Method4")
}
func (s *S) M5() {
	fmt.Println("S's Method5")
}
func (s *S) M6() {
	fmt.Println("S's Method6")
}

func testEmbedMethodOrder() {
	fmt.Println("------------")
	var t2 = T2{
		I1: &S{},
	}
	// T2 结构体在定义的时候 只实现了 M4
	// 但是在使用的时候 S 实现的M5,M6 也能直接调用到
	t2.M4()
	t2.M5()
	t2.M6()
}

// 基于一个类型创建另一个类型 主要有两种方式
// 1. 使用 defined 类型 如: type myInterface I
// 2. 使用 类型别名  如 type myInterface=I
// 如果使用define类型, defined 类型能否继承到 原有(underlying) 类型的方法需要分情况来看
// -- (1) 如果defined 类型是基于接口创建的，那么接口中所有的方法, defined 类型都能继承到
// -- (2) 如果defined 类型是基于非接口类型创建的，那么defined 类型的方法集合为空, 原有类型的方法需要它重新实现一遍
// 如果是使用类型别名 定义的新类型
// -- 新类型和原有类型 具有完全相同的 方法集合
// -- 事实上 rune 类型就是 int32 的类型别名, 使用 fmt.Printf 查看类型时看到的是 int32

type T3 struct{}

func (T3) M7()  {}
func (*T3) M8() {}

type Interface1 interface {
	M7()
	M8()
}
type T4 T3
type Interface2 Interface1

func testMethodSetExtends() {
	fmt.Println("------------")
	var t3 T3
	var pt3 *T3
	var t4 T4
	var pt4 *T4

	DumpMethodSet(&t3)
	DumpMethodSet(&t4)
	DumpMethodSet(&pt3)
	DumpMethodSet(&pt4)

	DumpMethodSet((*Interface1)(nil))
	DumpMethodSet((*Interface2)(nil))

	str := "hello"
	rs := []rune(str)
	fmt.Printf("type of rune is: %T\n", rs[0])
}

func testVariadicFuncArgs() {
	fmt.Println("------------")
	//fmt.Println(concat("-", 1, 2, 3))
	fmt.Println(concat("-", "hello", 1, 2, "world", 3, []string{"test", "args"}))
}

func concat(sep string, args ...interface{}) string {
	result := ""
	for i, arg := range args {
		if i != 0 {
			result += sep
		}
		switch arg := arg.(type) {
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
			result += fmt.Sprintf("%v", arg)
		case string:
			result += arg
		case []int:
			nums := make([]string, 0, len(arg))
			for _, i := range arg {
				nums = append(nums, fmt.Sprintf("%v", i))
			}
			join := strings.Join(nums, sep)
			result += join
		case []string:
			strs := make([]string, 0, len(arg))
			for _, sss := range arg {
				strs = append(strs, sss)
			}
			result += strings.Join(strs, sep)
		default:
			fmt.Println("unknown arg type:", reflect.TypeOf(arg))
			return ""
		}
	}
	return result
}
