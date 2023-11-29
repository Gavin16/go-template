package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// SayHello
// @BasePath /api/v1
// @Summary say hello
// @Schemes
// @Param name query true "用户名"
// @Description say hello to given name
// @Accept formData
// @Produce json
// @Success 200 {string} hello+name
// @Router /hello [get]
func SayHello(context *gin.Context) {
	request := context.Request
	name := request.FormValue("name")
	resp := "Hello," + name
	context.IndentedJSON(http.StatusOK, resp)
}
