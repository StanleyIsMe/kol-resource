package middleware

import (
	"flag"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/goleak"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	leak := flag.Bool("leak", false, "use leak detector")
	flag.Parse()

	if *leak {
		goleak.VerifyTestMain(m)

		return
	}

	os.Exit(m.Run())
}

func TestIsSkipTimeout(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		path string
		want bool
	}{
		{
			name: "upload path should skip",
			path: "/api/v1/kols/upload",
			want: true,
		},
		{
			name: "send_emails path should skip",
			path: "/api/v1/send_emails",
			want: true,
		},
		{
			name: "normal path should not skip",
			path: "/api/v1/kols",
			want: false,
		},
		{
			name: "root path should not skip",
			path: "/",
			want: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := isSkipTimeout(tt.path)
			if got != tt.want {
				t.Errorf("isSkipTimeout(%q) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}

func TestGinTimeout_Normal(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)

	r.Use(GinTimeout(5 * time.Second))
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	c.Request = httptest.NewRequest(http.MethodGet, "/test", nil)
	r.ServeHTTP(w, c.Request)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestGinTimeout_Skip(t *testing.T) {
	t.Parallel()

	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)

	r.Use(GinTimeout(1 * time.Millisecond))
	r.POST("/api/v1/kols/upload", func(c *gin.Context) {
		time.Sleep(10 * time.Millisecond)
		c.JSON(http.StatusOK, gin.H{"status": "uploaded"})
	})

	c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/kols/upload", nil)
	r.ServeHTTP(w, c.Request)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d (skip timeout), got %d", http.StatusOK, w.Code)
	}
}
