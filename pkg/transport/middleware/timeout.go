package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

// GinTimeout is wrap Timeout for compatibility with gin.
func GinTimeout(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip timeout middleware for specific paths
		if isSkipTimeout(c.Request.URL.Path) {
			c.Next()

			return
		}

		ctx := c.Request.Context()
		newCtx, cancelCtx := context.WithTimeout(ctx, timeout)

		c.Request = c.Request.WithContext(newCtx)

		defer func() {
			cancelCtx()

			if errors.Is(newCtx.Err(), context.DeadlineExceeded) {
				c.Writer.WriteHeader(http.StatusGatewayTimeout)
				zerolog.Ctx(newCtx).Error().Fields(map[string]any{
					"method":  c.Request.Method,
					"url":     fmt.Sprintf("%s%s", c.Request.Host, c.Request.RequestURI),
					"path":    c.Request.URL.Path,
					"timeout": timeout,
				}).Msg("request timeout")
			}
		}()

		c.Next()
	}
}

// isSkipTimeout checks if the request should be skipped for timeout middleware
func isSkipTimeout(path string) bool {
	return path == "/api/v1/kols/upload" || path == "/api/v1/send_emails"
}
