package main

import (
	"fmt"
	"testing"
)

func TestGrammar(t *testing.T) {
	h := c
	ss := s
	fmt.Println("------------")
	fmt.Printf("%T\n", h)
	fmt.Printf("%T\n", ss)

	// 切片的实现原理及应用
	var ar = [...]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	var sl1 = ar[3:7] // ar中下标为3,4,5,6 四个元素，不含7
	var sl2 = ar[4:8]
	fmt.Printf("slice sl1=%v\n", sl1)
	fmt.Printf("sl2=%v\n", sl2)
	sl1[1] = 15
	fmt.Printf("sl2=%v\n", sl2)
	fmt.Printf("slice ar=%T\n", ar)
	fmt.Printf("tpye of sl1=%T\n", sl1)
	fmt.Printf("cap of sl1=%d\n", cap(sl1))

	// 使用make 创建slice 第一个参数 len 代表slice已使用长度,同时也意味着append 操作的开始位置
	// 因此无特殊情况，len 一般都设置为 0
	var sl3 = make([]int, 0, 2)
	var sl4 = make([]int, 1, 2)
	sl3 = append(sl3, 3)
	sl3 = append(sl3, 4)
	fmt.Printf("sl3=%v,len=%v,cap=%v\n", sl3, len(sl3), cap(sl3))
	sl4 = append(sl4, 5)
	sl4 = append(sl4, 6)
	fmt.Printf("sl4=%v,len=%v,cap=%v\n", sl4, len(sl4), cap(sl4))

	sl5 := []string{"a", "b", "c", "d"}
	for i, v := range sl5 {
		fmt.Printf("i=%v, v=%v\n", i, v)
	}

	// 表达式求值
	pkgExpEvaluation()
	expressEvaluation()
	// 理解代码块和作用域
	fmt.Println("-----------")
	quiz()
	fmt.Println("-----------")
	testStrUnicode()
	fmt.Println("-----------")
	testForRangeVar()
	fmt.Println("-----------")
	testBreak()
}
