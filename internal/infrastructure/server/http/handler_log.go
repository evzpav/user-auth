package http

import (
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"gitlab.com/evzpav/documents/pkg/log"

	"github.com/gin-gonic/gin"
)

func (h *handler) logger() gin.HandlerFunc {
	const (
		headerTraceID  string = "X-Request-ID"
		headerClientID string = "X-NW-Client"
		message        string = "Request log"
	)

	hideAuthorizationData := func(k string, v []string) []string {
		if k == "Authorization" {
			v[0] = "<sensitive data hidden>"
		}

		return v
	}

	headerToMap := func(header http.Header) map[string]string {
		headers := make(map[string]string, len(header))
		for k, v := range header {
			v = hideAuthorizationData(k, v)
			headers[strings.ToLower(k)] = strings.ToLower(strings.Join(v, ","))
		}
		return headers
	}

	createEventByStatusCode := func(log log.Logger, status int) log.LoggerEvent {
		if status >= http.StatusBadRequest && status < http.StatusInternalServerError {
			return log.Warn()
		}

		if status >= http.StatusInternalServerError {
			return log.Error()
		}

		return log.Info()
	}

	return func(c *gin.Context) {
		start := time.Now()
		defer func() {
			elapsedTime := time.Since(start)

			req := c.Request
			reqHeader := req.Header
			reqURL := req.URL

			res := c.Writer
			resHeader := res.Header()
			resStatus := res.Status()

			event := createEventByStatusCode(h.log, resStatus)

			// Org
			applicationID := ""
			clientID := reqHeader.Get(headerClientID)

			event.Org(clientID, applicationID)

			// Trace
			traceID := reqHeader.Get(headerTraceID)
			event.Trace(traceID)

			// Request
			event.Req(traceID, c.ClientIP(), reqURL.Host, reqURL.Scheme, req.Method, reqURL.String(), "", headerToMap(reqHeader))

			// Response
			event.Res(resStatus, elapsedTime, "", res.Size(), headerToMap(resHeader))

			// Errors
			if len(c.Errors.Errors()) > 0 {
				event.ErrWithStack(c.Errors.Last().Err, string(debug.Stack()))
			}

			// Message
			event.Send(message)
		}()
		c.Next()
	}
}
