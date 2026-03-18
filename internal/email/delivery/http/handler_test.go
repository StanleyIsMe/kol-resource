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

	"kolresource/internal/admin"
	"kolresource/internal/email/mock/usecasemock"
	"kolresource/internal/email/usecase"

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

func setupRouter(handler *EmailHandler) *gin.Engine {
	r := gin.New()

	return r
}

func setAdminContext(c *gin.Context) {
	c.Set(admin.AdminIDKey, uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6"))
	c.Set(admin.AdminNameKey, "test-admin")
}

func TestEmailHandler_CreateEmailSender(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		body           any
		setupMock      func(mock *usecasemock.MockEmailUseCase)
		expectedStatus int
	}{
		{
			name: "success",
			body: CreateEmailSenderRequest{
				Name:      "Test Sender",
				Email:     "test@example.com",
				Key:       "test-key",
				RateLimit: 100,
			},
			setupMock: func(mock *usecasemock.MockEmailUseCase) {
				mock.EXPECT().CreateEmailSender(gomock.Any(), gomock.Any()).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "bad_request_empty_body",
			body:           map[string]any{},
			setupMock:      func(mock *usecasemock.MockEmailUseCase) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "bad_request_missing_required_fields",
			body: map[string]any{
				"name": "Test Sender",
			},
			setupMock:      func(mock *usecasemock.MockEmailUseCase) {},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			mockUsecase := usecasemock.NewMockEmailUseCase(ctrl)
			tt.setupMock(mockUsecase)

			handler := NewEmailHandler(mockUsecase, nil)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			setAdminContext(c)

			bodyBytes, _ := json.Marshal(tt.body)
			c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/email_senders", bytes.NewReader(bodyBytes))
			c.Request.Header.Set("Content-Type", "application/json")

			handler.CreateEmailSender(c)

			if w.Code != tt.expectedStatus {
				t.Errorf("CreateEmailSender() status = %d, want %d, body = %s", w.Code, tt.expectedStatus, w.Body.String())
			}
		})
	}
}

func TestEmailHandler_ListEmailSenders(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		setupMock      func(mock *usecasemock.MockEmailUseCase)
		expectedStatus int
	}{
		{
			name: "success",
			setupMock: func(mock *usecasemock.MockEmailUseCase) {
				mock.EXPECT().ListEmailSenders(gomock.Any()).Return([]usecase.EmailSender{
					{
						ID:    uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6"),
						Name:  "Test Sender",
						Email: "test@example.com",
					},
				}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "success_empty_list",
			setupMock: func(mock *usecasemock.MockEmailUseCase) {
				mock.EXPECT().ListEmailSenders(gomock.Any()).Return([]usecase.EmailSender{}, nil)
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			mockUsecase := usecasemock.NewMockEmailUseCase(ctrl)
			tt.setupMock(mockUsecase)

			handler := NewEmailHandler(mockUsecase, nil)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/email_senders", nil)

			handler.ListEmailSenders(c)

			if w.Code != tt.expectedStatus {
				t.Errorf("ListEmailSenders() status = %d, want %d, body = %s", w.Code, tt.expectedStatus, w.Body.String())
			}
		})
	}
}

func TestEmailHandler_GetEmailSender(t *testing.T) {
	t.Parallel()

	senderID := uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6")

	tests := []struct {
		name           string
		pathParam      string
		setupMock      func(mock *usecasemock.MockEmailUseCase)
		expectedStatus int
	}{
		{
			name:      "success",
			pathParam: senderID.String(),
			setupMock: func(mock *usecasemock.MockEmailUseCase) {
				mock.EXPECT().GetEmailSender(gomock.Any(), senderID).Return(&usecase.EmailSender{
					ID:    senderID,
					Name:  "Test Sender",
					Email: "test@example.com",
				}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "bad_request_invalid_uuid",
			pathParam:      "not-a-valid-uuid",
			setupMock:      func(mock *usecasemock.MockEmailUseCase) {},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			mockUsecase := usecasemock.NewMockEmailUseCase(ctrl)
			tt.setupMock(mockUsecase)

			handler := NewEmailHandler(mockUsecase, nil)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Params = gin.Params{{Key: "id", Value: tt.pathParam}}
			c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/email_senders/"+tt.pathParam, nil)

			handler.GetEmailSender(c)

			if w.Code != tt.expectedStatus {
				t.Errorf("GetEmailSender() status = %d, want %d, body = %s", w.Code, tt.expectedStatus, w.Body.String())
			}
		})
	}
}

func TestEmailHandler_SendEmail(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		body           any
		setupMock      func(mock *usecasemock.MockEmailUseCase)
		expectedStatus int
	}{
		{
			name: "success",
			body: SendEmailRequest{
				Subject:      "Test Subject",
				EmailContent: "Test Content",
				KolIDs:       []uuid.UUID{uuid.MustParse("1193487b-f1a2-7a72-8ae4-197b84dc52d6")},
				ProductID:    uuid.MustParse("2193487b-f1a2-7a72-8ae4-197b84dc52d6"),
				SenderID:     uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6"),
			},
			setupMock: func(mock *usecasemock.MockEmailUseCase) {
				mock.EXPECT().SendEmail(gomock.Any(), gomock.Any()).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "bad_request_empty_body",
			body:           map[string]any{},
			setupMock:      func(mock *usecasemock.MockEmailUseCase) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "bad_request_missing_required_fields",
			body: map[string]any{
				"subject": "Test Subject",
			},
			setupMock:      func(mock *usecasemock.MockEmailUseCase) {},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			mockUsecase := usecasemock.NewMockEmailUseCase(ctrl)
			tt.setupMock(mockUsecase)

			handler := NewEmailHandler(mockUsecase, nil)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			setAdminContext(c)

			bodyBytes, _ := json.Marshal(tt.body)
			c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/email_jobs", bytes.NewReader(bodyBytes))
			c.Request.Header.Set("Content-Type", "application/json")

			handler.SendEmail(c)

			if w.Code != tt.expectedStatus {
				t.Errorf("SendEmail() status = %d, want %d, body = %s", w.Code, tt.expectedStatus, w.Body.String())
			}
		})
	}
}

func TestEmailHandler_ListEmailJobs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		queryParams    string
		setupMock      func(mock *usecasemock.MockEmailUseCase)
		expectedStatus int
	}{
		{
			name:        "success",
			queryParams: "?page_index=1&page_size=10",
			setupMock: func(mock *usecasemock.MockEmailUseCase) {
				mock.EXPECT().ListEmailJobs(gomock.Any(), gomock.Any()).Return(&usecase.ListEmailJobsResponse{
					EmailJobs: []usecase.EmailJob{},
					Total:     0,
				}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:        "success_with_default_paging",
			queryParams: "",
			setupMock: func(mock *usecasemock.MockEmailUseCase) {
				mock.EXPECT().ListEmailJobs(gomock.Any(), gomock.Any()).Return(&usecase.ListEmailJobsResponse{
					EmailJobs: []usecase.EmailJob{},
					Total:     0,
				}, nil)
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			mockUsecase := usecasemock.NewMockEmailUseCase(ctrl)
			tt.setupMock(mockUsecase)

			handler := NewEmailHandler(mockUsecase, nil)
			r := setupRouter(handler)
			r.GET("/api/v1/email_jobs", handler.ListEmailJobs)

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/api/v1/email_jobs"+tt.queryParams, nil)

			r.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("ListEmailJobs() status = %d, want %d, body = %s", w.Code, tt.expectedStatus, w.Body.String())
			}
		})
	}
}

func TestEmailHandler_GetEmailJob(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		pathParam      string
		setupMock      func(mock *usecasemock.MockEmailUseCase)
		expectedStatus int
	}{
		{
			name:      "success",
			pathParam: "1",
			setupMock: func(mock *usecasemock.MockEmailUseCase) {
				mock.EXPECT().GetEmailJob(gomock.Any(), int64(1)).Return(&usecase.EmailJob{
					ID: 1,
				}, nil)
				mock.EXPECT().ListEmailLogs(gomock.Any(), usecase.ListEmailLogsParam{
					JobID: 1,
				}).Return([]usecase.EmailLog{}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "bad_request_invalid_id",
			pathParam:      "not-a-number",
			setupMock:      func(mock *usecasemock.MockEmailUseCase) {},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			mockUsecase := usecasemock.NewMockEmailUseCase(ctrl)
			tt.setupMock(mockUsecase)

			handler := NewEmailHandler(mockUsecase, nil)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Params = gin.Params{{Key: "id", Value: tt.pathParam}}
			c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/email_jobs/"+tt.pathParam, nil)

			handler.GetEmailJob(c)

			if w.Code != tt.expectedStatus {
				t.Errorf("GetEmailJob() status = %d, want %d, body = %s", w.Code, tt.expectedStatus, w.Body.String())
			}
		})
	}
}

func TestEmailHandler_CancelEmailJob(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		pathParam      string
		setupMock      func(mock *usecasemock.MockEmailUseCase)
		expectedStatus int
	}{
		{
			name:      "success",
			pathParam: "1",
			setupMock: func(mock *usecasemock.MockEmailUseCase) {
				mock.EXPECT().CancelEmailJob(gomock.Any(), int64(1)).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "bad_request_invalid_id",
			pathParam:      "not-a-number",
			setupMock:      func(mock *usecasemock.MockEmailUseCase) {},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			mockUsecase := usecasemock.NewMockEmailUseCase(ctrl)
			tt.setupMock(mockUsecase)

			handler := NewEmailHandler(mockUsecase, nil)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Params = gin.Params{{Key: "id", Value: tt.pathParam}}
			c.Request = httptest.NewRequest(http.MethodPut, "/api/v1/email_jobs/"+tt.pathParam+"/cancel", nil)

			handler.CancelEmailJob(c)

			if w.Code != tt.expectedStatus {
				t.Errorf("CancelEmailJob() status = %d, want %d, body = %s", w.Code, tt.expectedStatus, w.Body.String())
			}
		})
	}
}

func TestEmailHandler_StartEmailJob(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		pathParam      string
		setupMock      func(mock *usecasemock.MockEmailUseCase)
		expectedStatus int
	}{
		{
			name:      "success",
			pathParam: "1",
			setupMock: func(mock *usecasemock.MockEmailUseCase) {
				mock.EXPECT().StartEmailJob(gomock.Any(), int64(1)).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "bad_request_invalid_id",
			pathParam:      "not-a-number",
			setupMock:      func(mock *usecasemock.MockEmailUseCase) {},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			mockUsecase := usecasemock.NewMockEmailUseCase(ctrl)
			tt.setupMock(mockUsecase)

			handler := NewEmailHandler(mockUsecase, nil)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Params = gin.Params{{Key: "id", Value: tt.pathParam}}
			c.Request = httptest.NewRequest(http.MethodPut, "/api/v1/email_jobs/"+tt.pathParam+"/start", nil)

			handler.StartEmailJob(c)

			if w.Code != tt.expectedStatus {
				t.Errorf("StartEmailJob() status = %d, want %d, body = %s", w.Code, tt.expectedStatus, w.Body.String())
			}
		})
	}
}

func TestEmailHandler_UpdateEmailSender(t *testing.T) {
	t.Parallel()

	senderID := uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6")
	newName := "Updated Sender"

	tests := []struct {
		name           string
		pathParam      string
		body           any
		setupMock      func(mock *usecasemock.MockEmailUseCase)
		expectedStatus int
	}{
		{
			name:      "success",
			pathParam: senderID.String(),
			body: map[string]any{
				"name": newName,
			},
			setupMock: func(mock *usecasemock.MockEmailUseCase) {
				mock.EXPECT().UpdateEmailSender(gomock.Any(), gomock.Any()).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:      "bad_request_invalid_uri_uuid",
			pathParam: "not-a-valid-uuid",
			body: map[string]any{
				"name": newName,
			},
			setupMock:      func(mock *usecasemock.MockEmailUseCase) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:      "usecase_error",
			pathParam: senderID.String(),
			body: map[string]any{
				"name": newName,
			},
			setupMock: func(mock *usecasemock.MockEmailUseCase) {
				mock.EXPECT().UpdateEmailSender(gomock.Any(), gomock.Any()).Return(errors.New("update failed"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			mockUsecase := usecasemock.NewMockEmailUseCase(ctrl)
			tt.setupMock(mockUsecase)

			handler := NewEmailHandler(mockUsecase, nil)

			r := gin.New()
			r.PUT("/api/v1/email_senders/:id", func(c *gin.Context) {
				setAdminContext(c)
				handler.UpdateEmailSender(c)
			})

			w := httptest.NewRecorder()
			bodyBytes, _ := json.Marshal(tt.body)

			path := "/api/v1/email_senders/" + tt.pathParam
			req := httptest.NewRequest(http.MethodPut, path, bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")

			r.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("UpdateEmailSender() status = %d, want %d, body = %s", w.Code, tt.expectedStatus, w.Body.String())
			}
		})
	}
}

func TestEmailHandler_ListEmailSenders_UsecaseError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockUsecase := usecasemock.NewMockEmailUseCase(ctrl)
	mockUsecase.EXPECT().ListEmailSenders(gomock.Any()).Return(nil, errors.New("internal error"))

	handler := NewEmailHandler(mockUsecase, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/email_senders", nil)

	handler.ListEmailSenders(c)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("ListEmailSenders() status = %d, want %d, body = %s", w.Code, http.StatusInternalServerError, w.Body.String())
	}
}

func TestEmailHandler_SendEmail_UsecaseError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockUsecase := usecasemock.NewMockEmailUseCase(ctrl)
	mockUsecase.EXPECT().SendEmail(gomock.Any(), gomock.Any()).Return(errors.New("send failed"))

	handler := NewEmailHandler(mockUsecase, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	setAdminContext(c)

	body := SendEmailRequest{
		Subject:      "Test Subject",
		EmailContent: "Test Content",
		KolIDs:       []uuid.UUID{uuid.MustParse("1193487b-f1a2-7a72-8ae4-197b84dc52d6")},
		ProductID:    uuid.MustParse("2193487b-f1a2-7a72-8ae4-197b84dc52d6"),
		SenderID:     uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6"),
	}
	bodyBytes, _ := json.Marshal(body)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/email_jobs", bytes.NewReader(bodyBytes))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.SendEmail(c)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("SendEmail() status = %d, want %d, body = %s", w.Code, http.StatusInternalServerError, w.Body.String())
	}
}

func TestEmailHandler_GetEmailJob_ListEmailLogsError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockUsecase := usecasemock.NewMockEmailUseCase(ctrl)
	mockUsecase.EXPECT().GetEmailJob(gomock.Any(), int64(1)).Return(&usecase.EmailJob{ID: 1}, nil)
	mockUsecase.EXPECT().ListEmailLogs(gomock.Any(), usecase.ListEmailLogsParam{
		JobID: 1,
	}).Return(nil, errors.New("log query failed"))

	handler := NewEmailHandler(mockUsecase, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "1"}}
	c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/email_jobs/1", nil)

	handler.GetEmailJob(c)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("GetEmailJob() status = %d, want %d, body = %s", w.Code, http.StatusInternalServerError, w.Body.String())
	}
}

func TestEmailHandler_GetEmailJob_GetEmailJobError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockUsecase := usecasemock.NewMockEmailUseCase(ctrl)
	mockUsecase.EXPECT().GetEmailJob(gomock.Any(), int64(1)).Return(nil, errors.New("job not found"))

	handler := NewEmailHandler(mockUsecase, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "1"}}
	c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/email_jobs/1", nil)

	handler.GetEmailJob(c)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("GetEmailJob() status = %d, want %d, body = %s", w.Code, http.StatusInternalServerError, w.Body.String())
	}
}

func TestEmailHandler_CreateEmailSender_UsecaseError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockUsecase := usecasemock.NewMockEmailUseCase(ctrl)
	mockUsecase.EXPECT().CreateEmailSender(gomock.Any(), gomock.Any()).Return(errors.New("create failed"))

	handler := NewEmailHandler(mockUsecase, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	setAdminContext(c)

	body := CreateEmailSenderRequest{
		Name:      "Test Sender",
		Email:     "test@example.com",
		Key:       "test-key",
		RateLimit: 100,
	}
	bodyBytes, _ := json.Marshal(body)
	c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/email_senders", bytes.NewReader(bodyBytes))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.CreateEmailSender(c)

	if w.Code == http.StatusOK {
		t.Errorf("CreateEmailSender() status = %d, want non-200, body = %s", w.Code, w.Body.String())
	}
}

func TestEmailHandler_GetEmailSender_UsecaseError(t *testing.T) {
	t.Parallel()

	senderID := uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6")

	ctrl := gomock.NewController(t)
	mockUsecase := usecasemock.NewMockEmailUseCase(ctrl)
	mockUsecase.EXPECT().GetEmailSender(gomock.Any(), senderID).Return(nil, errors.New("get failed"))

	handler := NewEmailHandler(mockUsecase, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: senderID.String()}}
	c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/email_senders/"+senderID.String(), nil)

	handler.GetEmailSender(c)

	if w.Code == http.StatusOK {
		t.Errorf("GetEmailSender() status = %d, want non-200, body = %s", w.Code, w.Body.String())
	}
}

func TestEmailHandler_CancelEmailJob_UsecaseError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockUsecase := usecasemock.NewMockEmailUseCase(ctrl)
	mockUsecase.EXPECT().CancelEmailJob(gomock.Any(), int64(1)).Return(errors.New("cancel failed"))

	handler := NewEmailHandler(mockUsecase, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "1"}}
	c.Request = httptest.NewRequest(http.MethodPut, "/api/v1/email_jobs/1/cancel", nil)

	handler.CancelEmailJob(c)

	if w.Code == http.StatusOK {
		t.Errorf("CancelEmailJob() status = %d, want non-200, body = %s", w.Code, w.Body.String())
	}
}

func TestEmailHandler_StartEmailJob_UsecaseError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockUsecase := usecasemock.NewMockEmailUseCase(ctrl)
	mockUsecase.EXPECT().StartEmailJob(gomock.Any(), int64(1)).Return(errors.New("start failed"))

	handler := NewEmailHandler(mockUsecase, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "1"}}
	c.Request = httptest.NewRequest(http.MethodPut, "/api/v1/email_jobs/1/start", nil)

	handler.StartEmailJob(c)

	if w.Code == http.StatusOK {
		t.Errorf("StartEmailJob() status = %d, want non-200, body = %s", w.Code, w.Body.String())
	}
}

func TestEmailHandler_ListEmailJobs_UsecaseError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockUsecase := usecasemock.NewMockEmailUseCase(ctrl)
	mockUsecase.EXPECT().ListEmailJobs(gomock.Any(), gomock.Any()).Return(nil, errors.New("list failed"))

	handler := NewEmailHandler(mockUsecase, nil)
	r := setupRouter(handler)
	r.GET("/api/v1/email_jobs", handler.ListEmailJobs)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/email_jobs?page_index=1&page_size=10", nil)

	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("ListEmailJobs() status = %d, want %d, body = %s", w.Code, http.StatusInternalServerError, w.Body.String())
	}
}
