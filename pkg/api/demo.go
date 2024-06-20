package api

import (
	"github.com/gin-gonic/gin"
	"minsky/go-template/pkg/models"
	"net/http"
	"strconv"
)

// SayHello
// @Summary say hello,返参data为基本类型
// @Tags demo样例
// @Param name query string true "用户名"
// @Description say hello to given name
// @Accept multipart/form-data
// @Produce json
// @Success 200 {object} models.Result "成功"
// @Failure 500 {object} models.HTTPError "内部错误"
// @Router /api/demo/hello [get]
func SayHello(context *gin.Context) {
	request := context.Request
	name := request.FormValue("name")
	resp := "Hello," + name
	success := models.BuildSuccess(resp)
	context.IndentedJSON(http.StatusOK, success)
}

// GetUserById
// @Summary get user by id,返参为结构体对象
// @Tags Business-Group
// @Param id query int true "用户ID"
// @Description get user info by id
// @Accept multipart/form-data
// @Produce json
// @Success 200 {object} models.User "成功"
// @Failure 500 {object} models.HTTPError "内部错误"
// @Router /api/demo/getUserById [get]
func GetUserById(context *gin.Context) {
	request := context.Request
	userId := request.FormValue("id")

	user := models.User{}
	user.Id, _ = strconv.ParseInt(userId, 10, 64)
	user.Name = "Mostly"
	user.Age = 31
	user.Title = "marketing"
	user.Email = "Mostly@gmail.com"
	user.Nation = "USA"

	success := models.BuildSuccess(user)
	context.IndentedJSON(http.StatusOK, success)
}
