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
func GinTimeout(logger *zerolog.Logger, timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
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
