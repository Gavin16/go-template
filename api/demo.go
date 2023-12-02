package api

import (
	"github.com/gin-gonic/gin"
	"go-template/model"
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
// @Success 200 {object} model.Result "成功"
// @Failure 500 {object} model.HTTPError "内部错误"
// @Router /api/demo/hello [get]
func SayHello(context *gin.Context) {
	request := context.Request
	name := request.FormValue("name")
	resp := "Hello," + name
	success := model.BuildSuccess(resp)
	context.IndentedJSON(http.StatusOK, success)
}

// GetUserById
// @Summary get user by id,返参为结构体对象
// @Tags demo样例
// @Param id query int true "用户ID"
// @Description get user info by id
// @Accept multipart/form-data
// @Produce json
// @Success 200 {object} model.User "成功"
// @Failure 500 {object} model.HTTPError "内部错误"
// @Router /api/demo/getUserById [get]
func GetUserById(context *gin.Context) {
	request := context.Request
	userId := request.FormValue("id")

	user := model.User{}
	user.Id, _ = strconv.ParseInt(userId, 10, 64)
	user.Name = "Mostly"
	user.Age = 31
	user.Title = "marketing"
	user.Email = "Mostly@gmail.com"
	user.Nation = "USA"

	success := model.BuildSuccess(user)
	context.IndentedJSON(http.StatusOK, success)
}
