package main

import (
	_ "minsky/go-template/docs"          // swagger WebUI访问需要
	_ "minsky/go-template/internal/init" // 触发viper读取json配置
	"minsky/go-template/internal/routers"

	"github.com/gin-gonic/gin"
)

// @title go-template(项目名)
// @version v1.2
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
	routers.BindApi(router)
	err := router.Run(":8000")
	if err != nil {
		return
	}
}
