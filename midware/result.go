package midware

import (
	"github.com/gin-gonic/gin"
	"log"
	"minsky/go-template/model"
	"net/http"
)

// Recover error request recovery
func Recover(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("panic: %v\n", r)
			result := model.BuildError(errorToString(r))
			c.JSON(http.StatusInternalServerError, result)
			// stop subsequent executing
			c.Abort()
		}
	}()
	c.Next()
}

// convert error to string
func errorToString(r interface{}) string {
	switch v := r.(type) {
	case error:
		return v.Error()
	default:
		return r.(string)
	}
}
