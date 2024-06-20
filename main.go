package main

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "minsky/go-template/docs" // swagger WebUI访问需要
	"minsky/go-template/pkg/api"
	_ "minsky/go-template/pkg/init" // 触发viper读取json配置
	midware2 "minsky/go-template/pkg/midware"
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
	addMiddleware(router)
	bindApi(router)
	err := router.Run(":8000")
	if err != nil {
		return
	}
}

func addMiddleware(router *gin.Engine) {
	router.Use(midware2.Cors())
	router.Use(midware2.Recover)
}
