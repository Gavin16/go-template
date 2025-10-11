package api

import "fmt"

// API模块初始化
var swaggerUrl = "http://localhost:8000/swagger/index.html"

func init() {
	fmt.Println("Swagger API Access:", swaggerUrl)
}
