package middleware

import (
	"fmt"
	"net/http"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

// GinRecover is wrap Recover for compatibility with gin.
func GinRecover() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, err any) {
		const size = 64 << 10
		buf := make([]byte, size)
		buf = buf[:runtime.Stack(buf, false)]

		httpReq := c.Request
		zerolog.Ctx(c.Request.Context()).Error().Fields(map[string]any{
			"method":  httpReq.Method,
			"url":     fmt.Sprintf("%s%s", httpReq.Host, httpReq.RequestURI),
			"path":    httpReq.URL.Path,
			"request": httpReq.Body,
			"panic":   err,
			"stack":   string(buf),
		}).Msg("middleware.recover catch panic")

		c.Writer.WriteHeader(http.StatusInternalServerError)

		c.Next()
	})
}
