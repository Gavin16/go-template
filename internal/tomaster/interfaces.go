package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"
)

// ** 了解接口类型变量的内部表示

func main() {
	fmt.Println("--------------")
	// error nil != nil
	testErrorNil()

	// 接口多态测试
	// 结构体只会从结构体得到接口的方法
	// 结构体如果想要获取接口方法的几种主要途径
	// 1. 内嵌实现了该接口方法的其它结构体
	// 2. 内嵌了该接口，创建结构体时指定一个该接口的实现类作为接口的实现
	fmt.Println("--------------")
	var c = People{
		animal: &Dog{},
		name:   "chang",
	}
	c.run()

	fmt.Println("--------------")
	pool := poolLocal{}
	pool.Lock()
	fmt.Println("print in lock area")
	pool.Unlock()

	// ** 以接口为连接点的水平组合
	// 1.基本形式: 函数/方法中以接口类型作为参数-> 接口类型可能有多个实现,每个实现类都能传入
	// 2.包裹函数: 函数接受接口类型作为参数,并且返回该接口类型作为返参
	//  		  这样函数可以对接口传参做一次包裹, 从而可以实现对输入数据 过滤,装饰,变换等操作
	// 			  如CapReader 操作可以实现对reader读取到的字符串转换为大写
	// 3.适配器函数类型: 讲一个满足特定函数签名的普通函数 显式转换成自身类型的实例
	// 			       这样如果自身类型实现了某各接口，转换前的函数也就能作为该传参类型为该接口的传参了
	// 4.中间件: 包裹函数内部使用了适配器函数类型 讲一个普通函数转换成实现了某个接口类型的实例
	//  		包裹函数  +  适配器函数类型
	testCapReader()
}

// CapReader -- 以接口为连接点的水平组合
// 使用包裹函数
func CapReader(r io.Reader) io.Reader {
	return &capitalizedReader{r: r}
}

type capitalizedReader struct {
	r io.Reader
}

func (r *capitalizedReader) Read(p []byte) (n int, err error) {
	n, err = r.r.Read(p)
	if err != nil {
		return 0, err
	}
	q := bytes.ToUpper(p)
	for i, v := range q {
		p[i] = v
	}
	return n, nil
}

func testCapReader() {
	r := strings.NewReader("hello gopher!\n")
	r1 := CapReader(io.LimitReader(r, 5))
	if _, err := io.Copy(os.Stdout, r1); err != nil {
		log.Fatal(err)
	}
	fmt.Println()
}

// -- 类型嵌入实现垂直组合
type poolLocal struct {
	private interface{}
	shared  []interface{}
	sync.Mutex
	pad [128]byte
}

type Move interface {
	Run
	Fly
}

type People struct {
	animal Run
	name   string
}

func (p *People) run() {
	p.animal.run()
	fmt.Println("and people run")
}

type Run interface {
	run()
}

type Dog struct {
}
type Cat struct {
}

func (d *Dog) run() {
	fmt.Println("dog run")
}

//func (Cat) run() {
//	fmt.Println("cat run")
//}

type Fly interface {
	fly()
}

type Plain struct {
}
type Bird struct {
}

func (Plain) fly() {
	fmt.Println("Plain flying...")
}
func (Bird) fly() {
	fmt.Println("Bird flying...")
}

type MyError struct {
	error
}

var ErrBad = MyError{
	error: errors.New("bad"),
}

func bad() bool {
	return false
}

func returnsError() error {
	var p *MyError = nil
	if bad() {
		p = &ErrBad
	}
	return p
}

func testErrorNil() {
	e := returnsError()
	if e != nil {
		fmt.Printf("error:%v\n\n", e)
		return
	}
	fmt.Println("ok")
}
