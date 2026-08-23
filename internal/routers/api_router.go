package routers

import (
	"minsky/go-template/internal/api"
	midware2 "minsky/go-template/internal/midware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func BindApi(router *gin.Engine) {
	router.Use(midware2.Cors())
	router.Use(midware2.Recover)

	demo := router.Group("/api/demo")
	{
		demo.GET("/hello", api.SayHello)
		demo.GET("/getUserById", api.GetUserById)
	}
	// api group 绑定并发量
	feature := router.Group("/api/feature")
	{
		feature.GET("/cxtTimeOut", api.ContextTimeOut)
		feature.GET("/delay", api.DelayService)
		feature.GET("/rateTest", midware2.ConcurrentLimit(1), api.TesteRateLimit)
	}

	// swagger bind
	url := ginSwagger.URL("http://localhost:8000/swagger/doc.json")
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	//bind service api as follows
}
