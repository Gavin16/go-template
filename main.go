package main

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"minsky/go-template/api"
	_ "minsky/go-template/docs" // 这里需要引入,否则UI界面无法访问
	"minsky/go-template/midware"
)

func bindApi(router *gin.Engine) {
	demo := router.Group("/api/demo")
	{
		demo.GET("/hello", api.SayHello)
		demo.GET("/getUserById", api.GetUserById)
	}
	// swagger bind
	url := ginSwagger.URL("http://localhost:8000/swagger/doc.json")
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	//bind service api as follows
}

// @title go-template(Replace with your app name)
// @version 1.0
// @description 请求状态码定义
// @description code= 0, 调用成功
// @description code=-1, 系统错误
// @description code= 1, 提示接口返回的message
// @description code=10, 登录token过期
// @description code=20, 接口权限错误,没有权限访问该接口
// @termsOfService http://localhost:8000/swagger/index.html
// @license.name Apache 2.0
// @host localhost:8000
func main() {
	gin.ForceConsoleColor()
	router := gin.Default()
	addMiddleware(router)
	bindApi(router)
	err := router.Run(":8000")
	if err != nil {
		return
	}
}

func addMiddleware(router *gin.Engine) {
	router.Use(midware.Cors())
	router.Use(midware.Recover)
}

// grammar-review
//func main() {
//circle := biz.Circle{Radius: 1.0}
//area1 := circle.Area()
//fmt.Println(area1)
//rect := biz.Rectangle{Width: 1, Height: 2}
//area2 := rect.Area()
//fmt.Println(area2)

// Type Convert
// 1.number value convert
//	i := 12
//	f := float64(i)
//	d := float64(5)
//	res := f / d
//	fmt.Println(res)
//	cutRes := int(res)
//	fmt.Println("float to int:", cutRes)
//	// 2.string value convert
//	str := "123"
//	atop, _ := strconv.Atoi(str)
//	fmt.Println("string to int:", atop)
//	atop64, _ := strconv.ParseInt(str, 10, 64)
//	fmt.Println("string to int64:", atop64)
//
//	ii := 99
//	iiStr := strconv.Itoa(ii)
//	fmt.Println("int to string:", iiStr)
//
//	// 3.interface convert
//	// type assert
//	var ta interface{} = "hello"
//	tac, err := ta.(string)
//	if !err {
//		fmt.Println("assert fail:", err)
//	} else {
//		fmt.Println("convert success:", tac)
//	}
//
//	// type convert: from interface to interface
//	cc := biz.Circle{Radius: 1.0, XPos: 0, YPos: 1}
//	shape := biz.Shape(cc)
//	fmt.Println(shape.Area())
//
//	center := biz.Center(cc)
//	fmt.Println(center.Position())
//
//	// 4. slice
//	numbers := []int{0, 1, 2, 3, 4, 5, 6, 7, 8}
//	printSlice(numbers)
//
//	fmt.Println("numbers[1:4]==", numbers[1:4])
//	fmt.Println("numbers[:3]==", numbers[:3])
//	fmt.Println("numbers[4:]==", numbers[4:])
//
//	nums1 := make([]int, 0, 5)
//	nums1 = append(nums1, 1, 2, 3, 4, 5, 6, 7)
//	printSlice(nums1)
//	// append slice with slice
//	nums2 := []int{11, 22, 33}
//	nums1 = append(nums1, nums2...)
//	printSlice(nums1)
//
//	// slice len=0, size = 10
//	nums3 := make([]int, len(nums1))
//	//nums3 := make([]int, 0,10)
//	// copy the minimal len of len(src) and len(dst)
//	copy(nums3, nums1)
//	printSlice(nums3)
//
//	// goroutine and channel
//	s := []int{7, 2, 8, -5, 4, 0, 3}
//	ch := make(chan int)
//	go sum(s[:(len(s)/2)], ch)
//	go sum(s[(len(s)/2):], ch)
//	fmt.Println("front half sum is:", <-ch)
//	fmt.Println("later half sum is:", <-ch)
//
//	// channel buffer: put more data in channel buff before read
//	cBuff := make(chan int, 2)
//	cBuff <- 3
//	cBuff <- 4
//
//	fmt.Println(<-cBuff)
//	fmt.Println(<-cBuff)
//
//}

//func sum(s []int, c chan int) {
//	sum := 0
//	for _, v := range s {
//		sum += v
//	}
//	c <- sum
//}
//
//func printSlice(a []int) {
//	fmt.Printf("len=%d cap=%d slice=%v\n", len(a), cap(a), a)
//}
