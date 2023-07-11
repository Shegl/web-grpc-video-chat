package middleware

import (
	"github.com/gin-gonic/gin"
)

type CorsMiddleware struct{}

func (middleware *CorsMiddleware) Handle(c *gin.Context) {
	origin := c.GetHeader("Origin")
	if len(origin) == 0 {
		origin = "*"
	}

	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
	c.Writer.Header().Set("Access-Control-Max-Age", "86400")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Set-Cookie, "+
		"Access-Control-Allow-Origin, X-Requested-With, Authorization, Content-Type, Accept, Origin, User-Agent, DNT, Cache-Control, "+
		"Keep-Alive, If-Modified-Since, Accept, sentry-trace, "+
		"baggage, sec-ch-ua, sec-ch-ua-mobile, sec-ch-ua-platform")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(200)
	} else {
		c.Next()
	}
}

func NewCorsMiddleware() *CorsMiddleware {
	return &CorsMiddleware{}
}
