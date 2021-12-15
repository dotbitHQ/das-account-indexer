package toolib

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
)

var AllowOriginList = []string{
	`https?:\/\/127\.0\.0\.1:\d+`,
	`https?:\/\/localhost:\d+`,
}

func AllowOriginFunc(origin string) bool {
	for _, ao := range AllowOriginList {
		if ok, err := regexp.MatchString(ao, origin); err != nil {
			return false
		} else if ok {
			return true
		}
	}
	return false
}

func MiddlewareCors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.GetHeader("origin")
		ip := c.Request.RemoteAddr
		fmt.Println("MiddlewareCors:", method, origin, ip)
		if origin != "" {
			if AllowOriginFunc(origin) {
				c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			}
			c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Length,Content-Type")
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.Writer.Header().Set("content-type", "application/json")
		}
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		} else {
			c.Next()
		}
	}
}
