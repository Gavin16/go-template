package midware

import (
	"fmt"
	"minsky/go-template/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ConcurrentLimit(maxCurr int) gin.HandlerFunc {
	sem := make(chan struct{}, maxCurr)

	return func(c *gin.Context) {
		select {
		case sem <- struct{}{}:
			// 获得令牌
			defer func() {
				<-sem
			}()
			c.Next()
		default:
			c.AbortWithStatusJSON(
				http.StatusTooManyRequests,
				models.BuildError(fmt.Sprintf("请求繁忙, 最大支持并发为%d", maxCurr)))
		}

	}
}
