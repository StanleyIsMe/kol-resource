package middleware

import (
	"errors"
	"flag"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"kolresource/internal/admin"
	"kolresource/internal/admin/mock/usecasemock"
	"kolresource/internal/admin/usecase"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
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

func TestJWT_ValidToken(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	adminID := uuid.New()
	adminName := "testadmin"

	mockUC := usecasemock.NewMockAdminUseCase(ctrl)
	mockUC.EXPECT().
		LoginTokenParser(gomock.Any(), "valid-token").
		Return(&usecase.JWTAdminClaims{
			AdminID:   adminID,
			AdminName: adminName,
		}, nil)

	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)

	var gotAdminID uuid.UUID
	var gotAdminName string

	r.Use(JWT(mockUC))
	r.GET("/protected", func(c *gin.Context) {
		id, _ := c.Get(admin.AdminIDKey)
		gotAdminID = id.(uuid.UUID)

		name, _ := c.Get(admin.AdminNameKey)
		gotAdminName = name.(string)

		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	if gotAdminID != adminID {
		t.Errorf("admin_id = %v, want %v", gotAdminID, adminID)
	}

	if gotAdminName != adminName {
		t.Errorf("admin_name = %q, want %q", gotAdminName, adminName)
	}
}

func TestJWT_MissingHeader(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := usecasemock.NewMockAdminUseCase(ctrl)

	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)

	r.Use(JWT(mockUC))
	r.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestJWT_InvalidToken(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUC := usecasemock.NewMockAdminUseCase(ctrl)
	mockUC.EXPECT().
		LoginTokenParser(gomock.Any(), "bad-token").
		Return(nil, errors.New("invalid token"))

	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)

	r.Use(JWT(mockUC))
	r.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer bad-token")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, w.Code)
	}
}
