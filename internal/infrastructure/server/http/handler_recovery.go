package http

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func (h *handler) recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if recovered := recover(); recovered != nil {
				err, ok := recovered.(error)
				if ok {
					h.responseProblem(c, err)
					return
				}
				h.responseProblem(c, fmt.Errorf("%v", recovered))
			}
		}()
		c.Next()
	}
}
