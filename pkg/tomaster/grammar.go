package main

import (
	"fmt"
	"sync"
	"time"
)

// 包级变量: package 级别可见, 如果需要导出变量，则该包级变量也可以视为全局变量
// 包级变量需要使用 var 声明
// 推荐使用 如下声明方式
var (
	hello = "hello world"
	f     = float32(3.14)
	ss    = int32(111)
)

// 声明但是延迟初始化
// 一下写法虽然没有初始化,但是golang会让这些变量拥有初始的"零值"
// 如果是自定义类型的声明, 保证其零值可用是非常必要的
var a int32
var d float64

// 声明聚类与就近原则
// 声明时可以按照变量类型  或者 是否初始化来对待声明变量做划分
var (
	bufIoReaderPool   sync.Pool
	bufIoWriter2kPool sync.Pool
	bufIoWriter4kPool sync.Pool
)

var (
	aLongTimeAgo = time.Unix(0, 0)
	noDeadLine   = time.Time{}
	noCancel     = (chan struct{})(nil)
)

// 包级变量声明的位置
// 如果变量在包内被多出使用, 那么变量还是放在源文件头部声明比较合适
// 如果变量只在一个或者少数地方使用，那么紧挨着使用处进行声明更合适

// 对于函数或者方法内的局部变量
// 如果需要延迟使用 使用var关键字进行声明
// 如果声明且显示初始化的变量, 使用短变量声明形式
// 尽量在分支控制语句中 使用短变量声明形式

// golang 中无类型常量拥有像字面一样的特性
// 该特性是的无类型常量在参与变量赋值和计算过程时，无须显式声明类型转换
// 如果如果使用有类型的常量时  即使遇到底层类型一样的变量做计算处理,编译器也会报错
const (
	c = 5
	s = "hello,gopher"
)

// ** 使用iota 实现枚举常量
// iota 是go语言的一个预定义标识符，它表示const声明块中每个常量所处位置在块中的偏移值
const (
	mutexLocked = 1 << iota
	mutexWoken
	mutexStarving
	mutexWaiterShift      = iota
	starvationThresholdNs = 1e6
)

// ** 尽量定义零值可用类型
// golang 内置数据类型具有零值可用的特性
// 其中各类型的零值如下
// 所有整数类型: 0
// 浮点类型: 0.0
// 布尔类型: false
// 字符串类型: ""
// 指针，interface,slice, channel,map, function: nil
// 对于零值为 nil 的类型如 slice, 即使slice 初始值为nil 也可以对它使用append 方法进行值追加
func testSliceNilUsage() {
	var zeroSlice []int
	zeroSlice = append(zeroSlice, 1)
	zeroSlice = append(zeroSlice, 2)
	fmt.Println(zeroSlice) // 输出[1,2]
}

// 零值可用切片不支持下标索引
// 同时map 这类原生类型没有提供对零值可用的支持
var ai []int

// s[0] = 12 // 错误
var mm map[string]int

// mm["go"] = 1 // 错误
// 此外还需要注意尽量避免值复制
//var mu sync.Mutex

// mu1 := mu // 错误
//var mu1 *sync.Mutex = &mu // 指针复制是可以的

// ** 使用符合字面值作为初值构造器
// 对于较为复杂的符合类型
type myStruct struct {
	name string
	age  int
}

var st0 = new(myStruct)                   // 不推荐
var st1 = myStruct{"tony", 32}            // 不推荐
var st2 = myStruct{name: "tony", age: 32} // 推荐使用
var arr = [5]int{1, 2, 3, 4, 5}
var sl = []int{10, 20, 30, 40, 50, 60}
var mp = map[int]string{1: "hello", 2: "golang", 3: "gopher", 4: "!"}
var m1 = map[string]*myStruct{
	"Admin": {"zhansan", 33},
	"Root":  {"Lisi", 30},
}

// ** 切片的实现原理及使用
var ar = [9]int{1, 2, 3, 4, 5, 6, 7, 8, 9} // ar 类型为[9]int
var sl1 = ar[3:5]                          // 此时 sl1 变成了切片类型 []int
var slc = cap(sl1)                         // sl1 是从数组ar第3个元素开始切取,后面还有6个元素, 因此cap(sl1) 为9-3=6
// 假设底层数组为 ar = [10]int{11,12,13,14,15,16,17,18,19,20}
// 现在创建 sla1 = ar[1:5]; sla2 = ar[2:6]
// 现在将 sla1[1] 设置为 22 --> sla1[1] = 22, 那么响应的 sla2[0] 也将被设置为 22
// 因为 sla1 和 sla2 底层的数组都是同一个
// sl1 = append(sl1, 16) // ar={1, 2, 3, 4, 5, 16, 7, 8, 9}, sl1 = {4,5,16}
// sl1 = append(sl1, 17) // ar={1, 2, 3, 4, 5, 16, 17, 8, 9}, sl1 = {4,5,16,17}
// sl1 = append(sl1, 18) // ar={1, 2, 3, 4, 5, 16, 17, 18, 9}, sl1 = {4,5,16,17,18}
// sl1 = append(sl1, 19) // ar={1, 2, 3, 4, 5, 16, 17, 18, 19}, sl1 = {4,5,16,17,18,19}
// sl1 = append(sl1, 20) // ar={1, 2, 3, 4, 5, 16, 17, 18, 19}, sl1 = {4,5,16,17,18,19,20}
//
//	--> 至此sl1将与ar解除关联，而使用新的底层数组进行存储
//
// 由于slice底层基于数据来存储，向slice中动态添加元素时 会存在动态扩容
// 因此如果能提前知道slice 可能需要存储多少数据的前提下, 创建slice 时应该尽量指定 cap 参数
var sl2 = make([]int, 0)    //不推荐
var sl3 = make([]int, 0, 8) //推荐

// ** 理解Go语言的包导入
// go 项目中包名习惯上都会和包所在的最后一级路径名相同, 但是不严格限制
// 使用 go build 指令进行编译时会经历如下步骤
// 0. go 编译器会根据各个go 文件导入的依赖进行包搜索
// 1. 如果mod文件中引入了，说明是第三方依赖，如果没有引入(模块名和当前文件模块名相同)则视为项目内依赖
// 2. 对于第三方依赖，依赖的导入会先去 $GOPATH/pkg/mod 缓存中查看 是否缓存了当前模块，若缓存了当前模块则直接读取
// 3.   若没有缓存当前模块，则从 $GOPROXY 配置的代理地址去下载当前依赖的包(下载是以整个模块进行下载,而不是以依赖的包为单位)
// 4. 开始编译各个文件得到目标文件，然后链接各个文件得到可执行文件
// 此外还有些注意事项
// a. 若依赖的不同模块(三方库)出现了包名冲突的情况，在import 时可以对包名起别名的方式避免
// b. 若希望导入包，但是又不想使用(构建项目依赖关系时可能会用到)，可以使用 import _ "github.com/demo/placeholder" 方式进行包导入

// ** 表达式求值顺序
// 包级别变量初始化
var (
	a1 = c1 + b1
	b1 = ff()
	c1 = ff()
	d1 = 3
)

var (
	a2 = b2 + ff2()
	b2 = ff2()
	c2 = b2 + ff2()
	d2 = 1
)

func ff2() int {
	d2++
	return d2
}

func ff() int {
	d1++
	return d1
}
func pkgExpEvaluation() {
	fmt.Println(a1, b1, c1, d1)
	// 包级别变量表达式求值 按照变量声明顺序,一个一个找直到找到能求值的表达式,
	// 找到之后求完一个之后，是反过头去求第一个还是继续往下走 求完后面所有的表达式的值？
	// -- 是求完一个之后，马上返回头去求之前第一个还没求出来的表达式的值！！！
	fmt.Println(a2, b2, c2, d2)
}

// **局部变量初始化
// 表达式内的所有函数,方法以及channel 操作 按照从左至右的顺序进行求值
func expressEvaluation() {
	var n0 = 1
	var n1 = 1
	n0, n1 = n0+n1, n0
	fmt.Println(n0, n1)

	fmt.Println("-------------")
	switch Expr(2) {
	case Expr(1), Expr(2), Expr(3):
		fmt.Println("enter into case1")
		// fallthrough -- 进入下一个case 语句进行判断; 默认执行完一个case语句之后就结束不再继续执行
	case Expr(4):
		fmt.Println("enter into case2")
	}
}

func Expr(n int) int {
	fmt.Println(n)
	return n
}

// ** 理解Go语言代码块和作用域
func quiz() {
	if a := 1; false {
	} else if b := 2; false {
	} else if c := 3; false {
	} else {
		fmt.Println(a, b, c)
	}
}

// 以上quiz 函数体可以等效转化为
//{
//	a := 1
//	if false{
//	}else{
//		{
//			b:=2
//			if false{
//			}else{
//				{
//					c:=3
//					if false{
//					}else{
//						fmt.println(a,b,c)
//					}
//				}
//			}
//		}
//	}
//}

// ** 字符串增强
// 默认零值可用
var str string // str = ""
// 由于 string类型数据是不可变的, 因此长度也不可变，可以将长度直接存储,故 len(str) 时间复杂度为 O(1)
// 支持通过 '+','+=' 操作符进行操作
var str1 string = str + "hello"

func testStr() {
	str2 := "hi"
	str2 += str1 // str2 = str2 + str1
}

// 支持 ==, !=, >=, <=, > 和 < 比较关系操作符
func testCompare() {
	str3 := "中国"
	str4 := "中" + "国"
	fmt.Println(str3 == str4)

	str3 = "123"
	str4 = "12345"
	fmt.Println(str3 > str4)
}

// 对非ASCII提供原生支持
// 原生支持多行字符串
func testStrUnicode() {
	str = "中国欢迎您"
	rs := []rune(str)
	sl := []byte(str)
	for i, v := range rs {
		var utf8bytes []byte
		for j := i * 3; j < (i+1)*3; j++ {
			utf8bytes = append(utf8bytes, sl[j])
		}
		fmt.Printf("%s => %X => %X\n", string(v), v, utf8bytes)
	}

	strA := string(rs)
	fmt.Println(strA)
	strB := string(sl)
	fmt.Println(strB)

	fmt.Println("-------------")
	// 使用 `` 创建多行字符串
	var strs = `野径云俱黑
江船火烛明`
	fmt.Println(strs)
}

// ** Go控制语句惯用法和注意事项
// 1. 出现错误逻辑时，快速返回
// 2. 成功逻辑不要嵌入到if-else 语句中
// 3. "快乐路径"的执行逻辑在代码布局上始终靠左, 这样读者可以一眼看到该函数的正常逻辑流程
// 4. "快乐路径"的返回值一般在函数的最后一行, 就像上面伪代码段1中的那样

// 注意for range 中声明变量的作用域
func testForRangeVar() {
	a := [...]int{1, 2, 3, 4, 5}
	// 并发访问range 中的局部变量时 可能拿到的都是最后相同的值
	// 建议增加闭包函数的传参,然后将 i,v 以参数形式传过去
	for i, v := range a {
		go func() {
			time.Sleep(100 * time.Microsecond)
			fmt.Println(i, v)
		}()
	}
	time.Sleep(200 * time.Millisecond)

	// 尝试修改数组a 中的内容
	// 结果无法修改，因为a 是数组，数组在 range 语句中是通过值传递，range 遍历的是a的副本
	fmt.Println("--------------")
	fmt.Println("a = ", a)
	r := make([]int, 0, 5)
	for i, v := range a {
		if i == 0 {
			a[1] = 12
			a[2] = 13
		}
		r = append(r, v)
	}
	fmt.Println("r = ", r)
	fmt.Println("a = ", a)

	// 对切片进行遍历, 由于切片实际存储的是一个(*T, len, cap)的三元组， *T指向的是底层数组的指针
	// 因此即使 range 遍历的是 切片的副本 也同样能通过*T 指针修改到底层的数组
	fmt.Println("------------")
	s := []int{1, 2, 3, 4, 5}
	rs := make([]int, 0, 5)
	for i, v := range s {
		if i == 0 {
			s[1] = 12
			s[2] = 13
		}
		rs = append(rs, v)
	}
	fmt.Println("s = ", s)
	fmt.Println("r = ", rs)

	// 此range 遍历其它类型的表达式(string, map, channel) 依然使用的是这些变量的副本
	// 具体特性需要根据变量类型的存储方式确定

	// 当遍历字符串时, range 遍历出来的每个元素都是 rune
	fmt.Println("--------------")
	var sc = "我们中国人"
	var se = "we Chinese"
	for i, v := range sc {
		fmt.Printf("%d  %s 0x%x %T\n", i, string(v), v, v)
	}
	for i, v := range se {
		fmt.Printf("%d  %s 0x%x %T\n", i, string(v), v, v)
	}

	// 遍历map时， map在Go运行时表示为一个hmap的描述指针，因此该指针的副本也指向同一个hmap描述符
	// 循环过程中对map 进行修改时可以修改到底层存储的数据，但是
	// 由于map 遍历的顺序具有不确定性，修改的内容譬如 新增,删除 的操作，有可能因为顺序的原因，不会出现在后续的遍历中
	fmt.Println("-------------")
	var m = map[string]int{
		"tony": 21,
		"tom":  22,
		"Jim":  23,
	}
	counter := 0
	for k, v := range m {
		if counter == 0 {
			delete(m, "tony")
		}
		counter++
		fmt.Println(k, v)
	}
	fmt.Println("counter is:", counter)
	fmt.Println("m len is:", len(m))

	// channel 使用for range 遍历
	var c = make(chan int)
	go func() {
		time.Sleep(200 * time.Millisecond)
		c <- 1
		c <- 2
		c <- 3
		// 若不关闭channel 另一端读取时会一直阻塞
		// go检测机制会发现死锁并报错 fatal error: all goroutines are asleep - deadlock!
		close(c)
	}()
	for v := range c {
		fmt.Println(v)
	}
}

func testBreak() {
	// 默认情况下(不接label时)break 语句只会跳出同一函数内break 语句所在的最内层的for,switch 和 select语句
	// 若希望break 跳出到指定位置, 需要定义label,然后再break 后指定label
	exit := make(chan interface{})
	go func() {
	loop:
		for {
			select {
			case <-time.After(time.Second):
				fmt.Println("tick")
			case <-exit:
				fmt.Println("exiting...")
				break loop
			}
		}
		fmt.Println("exit!")
	}()
	time.Sleep(3 * time.Second)
	exit <- struct{}{}
	time.Sleep(500 * time.Millisecond)
}
