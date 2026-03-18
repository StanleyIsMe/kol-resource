package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func TestGinContextLogger(t *testing.T) {
	t.Parallel()

	logger := zerolog.Nop()

	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)

	var ctxHasLogger bool

	r.Use(GinContextLogger(&logger))
	r.GET("/test", func(c *gin.Context) {
		l := zerolog.Ctx(c.Request.Context())
		ctxHasLogger = l != nil
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, c.Request)

	if !ctxHasLogger {
		t.Error("expected logger to be present in request context")
	}

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}
