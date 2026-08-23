package api

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"minsky/go-template/internal/biz"
	"minsky/go-template/internal/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
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

func TesteRateLimit(c *gin.Context) {
	time.Sleep(10 * time.Second)
	c.JSON(http.StatusOK, models.BuildSuccess("success"))
}

// ContextTimeOut 调用超时
func ContextTimeOut(c *gin.Context) {
	tmCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	req, _ := http.NewRequestWithContext(tmCtx, "GET", "http://localhost:8000/api/feature/delay", nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			c.JSON(http.StatusGatewayTimeout, models.BuildError(fmt.Sprintf("%s", tmCtx.Err())))
			return
		}
	}
	defer resp.Body.Close()

	// 根据业务需要转发下游响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			models.BuildError(fmt.Sprintf("%s", tmCtx.Err())),
		)
		return
	}
	c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), body)
}

// DelayService 延时等待5秒
func DelayService(c *gin.Context) {
	ctx := c.Request.Context()
	slog.Info("=====》进入DelayService")

	timer := time.NewTimer(10 * time.Second)
	defer timer.Stop()
	select {
	case <-timer.C:
		c.IndentedJSON(http.StatusOK, "OK")
		slog.Info("《===== DelayService 完成执行")
	case <-ctx.Done():
		// 请求取消, 不需要返回响应, 直接return
		slog.Info("DelayService 被取消:", "context err", ctx.Err())
		return
	}
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
	id, _ := strconv.ParseInt(userId, 10, 64)
	userService := biz.UserSvcImpl{}
	user := userService.GetUserById(id)

	success := models.BuildSuccess(user)
	context.IndentedJSON(http.StatusOK, success)
}
