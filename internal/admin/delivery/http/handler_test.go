package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"kolresource/internal/admin/mock/usecasemock"
	"kolresource/internal/admin/usecase"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"go.uber.org/goleak"
)

func TestMain(m *testing.M) {
	leak := flag.Bool("leak", false, "use leak detector")
	flag.Parse()

	if *leak {
		goleak.VerifyTestMain(m)

		return
	}

	os.Exit(m.Run())
}

func TestAdminHandler_Register(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		body           any
		setupMock      func(mock *usecasemock.MockAdminUseCase)
		expectedStatus int
	}{
		{
			name: "success",
			body: RegisterRequest{
				Name:     "test-admin",
				UserName: "admin@example.com",
				Password: "securepassword",
			},
			setupMock: func(mock *usecasemock.MockAdminUseCase) {
				mock.EXPECT().Register(gomock.Any(), usecase.RegisterParams{
					Name:     "test-admin",
					UserName: "admin@example.com",
					Password: "securepassword",
				}).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "bad_request_empty_body",
			body:           map[string]any{},
			setupMock:      func(mock *usecasemock.MockAdminUseCase) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "bad_request_missing_name",
			body: map[string]any{
				"user_name": "admin@example.com",
				"password":  "securepassword",
			},
			setupMock:      func(mock *usecasemock.MockAdminUseCase) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "bad_request_missing_password",
			body: map[string]any{
				"name":      "test-admin",
				"user_name": "admin@example.com",
			},
			setupMock:      func(mock *usecasemock.MockAdminUseCase) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "bad_request_invalid_email",
			body: map[string]any{
				"name":      "test-admin",
				"user_name": "not-an-email",
				"password":  "securepassword",
			},
			setupMock:      func(mock *usecasemock.MockAdminUseCase) {},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			mockUsecase := usecasemock.NewMockAdminUseCase(ctrl)
			tt.setupMock(mockUsecase)

			handler := NewAdminHandler(mockUsecase)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			bodyBytes, _ := json.Marshal(tt.body)
			c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/register", bytes.NewReader(bodyBytes))
			c.Request.Header.Set("Content-Type", "application/json")

			handler.Register(c)

			if w.Code != tt.expectedStatus {
				t.Errorf("Register() status = %d, want %d, body = %s", w.Code, tt.expectedStatus, w.Body.String())
			}
		})
	}
}

func TestAdminHandler_Login(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		body           any
		setupMock      func(mock *usecasemock.MockAdminUseCase)
		expectedStatus int
		checkBody      func(t *testing.T, body []byte)
	}{
		{
			name: "success",
			body: LoginRequest{
				UserName: "admin@example.com",
				Password: "securepassword",
			},
			setupMock: func(mock *usecasemock.MockAdminUseCase) {
				mock.EXPECT().Login(gomock.Any(), "admin@example.com", "securepassword").Return(&usecase.LoginResponse{
					Token:     "jwt-token-string",
					AdminName: "test-admin",
				}, nil)
			},
			expectedStatus: http.StatusOK,
			checkBody: func(t *testing.T, body []byte) {
				t.Helper()

				var resp LoginResponse
				if err := json.Unmarshal(body, &resp); err != nil {
					t.Fatalf("failed to unmarshal response: %v", err)
				}

				if resp.Token != "jwt-token-string" {
					t.Errorf("Login() token = %s, want jwt-token-string", resp.Token)
				}

				if resp.AdminName != "test-admin" {
					t.Errorf("Login() admin_name = %s, want test-admin", resp.AdminName)
				}
			},
		},
		{
			name:           "bad_request_empty_body",
			body:           map[string]any{},
			setupMock:      func(mock *usecasemock.MockAdminUseCase) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "bad_request_missing_password",
			body: map[string]any{
				"user_name": "admin@example.com",
			},
			setupMock:      func(mock *usecasemock.MockAdminUseCase) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "bad_request_invalid_email",
			body: map[string]any{
				"user_name": "not-an-email",
				"password":  "securepassword",
			},
			setupMock:      func(mock *usecasemock.MockAdminUseCase) {},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			mockUsecase := usecasemock.NewMockAdminUseCase(ctrl)
			tt.setupMock(mockUsecase)

			handler := NewAdminHandler(mockUsecase)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			bodyBytes, _ := json.Marshal(tt.body)
			c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewReader(bodyBytes))
			c.Request.Header.Set("Content-Type", "application/json")

			handler.Login(c)

			if w.Code != tt.expectedStatus {
				t.Errorf("Login() status = %d, want %d, body = %s", w.Code, tt.expectedStatus, w.Body.String())
			}

			if tt.checkBody != nil {
				tt.checkBody(t, w.Body.Bytes())
			}
		})
	}
}

func TestAdminHandler_Register_UsecaseError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockUsecase := usecasemock.NewMockAdminUseCase(ctrl)
	mockUsecase.EXPECT().Register(gomock.Any(), usecase.RegisterParams{
		Name:     "test-admin",
		UserName: "admin@example.com",
		Password: "securepassword",
	}).Return(errors.New("register failed"))

	handler := NewAdminHandler(mockUsecase)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	body := RegisterRequest{
		Name:     "test-admin",
		UserName: "admin@example.com",
		Password: "securepassword",
	}
	bodyBytes, _ := json.Marshal(body)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/register", bytes.NewReader(bodyBytes))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Register(c)

	if w.Code == http.StatusOK {
		t.Errorf("Register() status = %d, want non-200, body = %s", w.Code, w.Body.String())
	}
}

func TestAdminHandler_Login_UsecaseError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockUsecase := usecasemock.NewMockAdminUseCase(ctrl)
	mockUsecase.EXPECT().Login(gomock.Any(), "admin@example.com", "securepassword").Return(nil, errors.New("login failed"))

	handler := NewAdminHandler(mockUsecase)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	body := LoginRequest{
		UserName: "admin@example.com",
		Password: "securepassword",
	}
	bodyBytes, _ := json.Marshal(body)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewReader(bodyBytes))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Login(c)

	if w.Code == http.StatusOK {
		t.Errorf("Login() status = %d, want non-200, body = %s", w.Code, w.Body.String())
	}
}
