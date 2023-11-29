package main

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go-template/api"
	_ "go-template/docs" // 这里需要引入,否则UI界面无法访问
)

func bindApi(router *gin.Engine) {
	v1 := router.Group("/api/v1")
	{
		example := v1.Group("/example")
		{
			// hello api
			example.GET("/hello", api.SayHello)
		}
	}

	// swagger bind
	url := ginSwagger.URL("http://localhost:8000/swagger/doc.json")
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	//bind service api as follows

}

// @title go-template(Replace with your app name)
// @version 1.0
// @description go web template(Replace with your description)
// @termsOfService http://localhost:8000/swagger/index.html
// @license.name Apache 2.0
// @host localhost:8000
// @BasePath /api/v1
func main() {
	router := gin.Default()
	bindApi(router)
	err := router.Run(":8000")
	if err != nil {
		return
	}
}
