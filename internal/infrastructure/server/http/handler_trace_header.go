package http

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/evzpav/user-auth/pkg/trace"
)

func (h *handler) traceHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(trace.RequestIDHeader, c.GetHeader(trace.RequestIDHeader))
		c.Next()
	}
}
