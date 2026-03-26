package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestGinRecover(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)

	r.Use(GinRecover())
	r.GET("/panic", func(c *gin.Context) {
		panic("test panic")
	})

	c.Request = httptest.NewRequest(http.MethodGet, "/panic", nil)
	r.ServeHTTP(w, c.Request)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestGinRecover_NoPanic(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)

	r.Use(GinRecover())
	r.GET("/ok", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	c.Request = httptest.NewRequest(http.MethodGet, "/ok", nil)
	r.ServeHTTP(w, c.Request)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}
