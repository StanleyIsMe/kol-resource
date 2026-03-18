package http

import (
	"encoding/json"
	"errors"
	"flag"
	"kolresource/internal/admin"
	"kolresource/internal/kol"
	"kolresource/internal/kol/mock/usecasemock"
	"kolresource/internal/kol/usecase"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

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

func TestKolHandler_CreateKol(t *testing.T) {
	t.Parallel()

	adminID := uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6")
	tagID := uuid.MustParse("1193487b-f1a2-7a72-8ae4-197b84dc52d6")

	tests := []struct {
		name           string
		body           string
		expectedStatus int
		getMock        func(ctrl *gomock.Controller) *usecasemock.MockKolUseCase
	}{
		{
			name: "success",
			body: `{"name":"test","email":"test@example.com","description":"desc","social_media":"ig","sex":"m","tags":["1193487b-f1a2-7a72-8ae4-197b84dc52d6"]}`,
			expectedStatus: http.StatusOK,
			getMock: func(ctrl *gomock.Controller) *usecasemock.MockKolUseCase {
				mock := usecasemock.NewMockKolUseCase(ctrl)
				mock.EXPECT().CreateKol(gomock.Any(), usecase.CreateKolParam{
					Name:           "test",
					Email:          "test@example.com",
					Description:    "desc",
					SocialMedia:    "ig",
					Sex:            kol.SexMale,
					Tags:           []uuid.UUID{tagID},
					UpdatedAdminID: adminID,
				}).Return(nil)

				return mock
			},
		},
		{
			name:           "invalid_json",
			body:           `{invalid}`,
			expectedStatus: http.StatusBadRequest,
			getMock: func(ctrl *gomock.Controller) *usecasemock.MockKolUseCase {
				return usecasemock.NewMockKolUseCase(ctrl)
			},
		},
		{
			name:           "missing_required_fields",
			body:           `{"description":"only desc"}`,
			expectedStatus: http.StatusBadRequest,
			getMock: func(ctrl *gomock.Controller) *usecasemock.MockKolUseCase {
				return usecasemock.NewMockKolUseCase(ctrl)
			},
		},
		{
			name: "usecase_error",
			body: `{"name":"test","email":"test@example.com","sex":"m"}`,
			expectedStatus: http.StatusInternalServerError,
			getMock: func(ctrl *gomock.Controller) *usecasemock.MockKolUseCase {
				mock := usecasemock.NewMockKolUseCase(ctrl)
				mock.EXPECT().CreateKol(gomock.Any(), gomock.Any()).Return(errors.New("internal error"))

				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			mockUC := tt.getMock(ctrl)
			handler := NewKolHandler(mockUC)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/kols", strings.NewReader(tt.body))
			c.Request.Header.Set("Content-Type", "application/json")
			c.Set(admin.AdminIDKey, adminID)

			handler.CreateKol(c)

			if w.Code != tt.expectedStatus {
				t.Errorf("CreateKol() status = %d, want %d, body = %s", w.Code, tt.expectedStatus, w.Body.String())
			}
		})
	}
}

func TestKolHandler_GetKolByID(t *testing.T) {
	t.Parallel()

	kolID := uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6")

	tests := []struct {
		name           string
		paramID        string
		expectedStatus int
		getMock        func(ctrl *gomock.Controller) *usecasemock.MockKolUseCase
		checkBody      func(t *testing.T, body []byte)
	}{
		{
			name:           "success",
			paramID:        kolID.String(),
			expectedStatus: http.StatusOK,
			getMock: func(ctrl *gomock.Controller) *usecasemock.MockKolUseCase {
				mock := usecasemock.NewMockKolUseCase(ctrl)
				mock.EXPECT().GetKolByID(gomock.Any(), kolID).Return(&usecase.Kol{
					ID:          kolID,
					Name:        "Test Name",
					Email:       "test@example.com",
					Description: "Test Description",
					SocialMedia: "ig",
					Sex:         kol.SexMale,
					Tags:        []usecase.Tag{{ID: uuid.MustParse("1193487b-f1a2-7a72-8ae4-197b84dc52d6"), Name: "Tag1"}},
				}, nil)

				return mock
			},
			checkBody: func(t *testing.T, body []byte) {
				var result usecase.Kol
				if err := json.Unmarshal(body, &result); err != nil {
					t.Fatalf("failed to unmarshal response: %v", err)
				}

				if result.ID != kolID {
					t.Errorf("expected kol ID %v, got %v", kolID, result.ID)
				}

				if result.Name != "Test Name" {
					t.Errorf("expected name %q, got %q", "Test Name", result.Name)
				}
			},
		},
		{
			name:           "invalid_uuid",
			paramID:        "not-a-uuid",
			expectedStatus: http.StatusBadRequest,
			getMock: func(ctrl *gomock.Controller) *usecasemock.MockKolUseCase {
				return usecasemock.NewMockKolUseCase(ctrl)
			},
		},
		{
			name:           "usecase_error",
			paramID:        kolID.String(),
			expectedStatus: http.StatusInternalServerError,
			getMock: func(ctrl *gomock.Controller) *usecasemock.MockKolUseCase {
				mock := usecasemock.NewMockKolUseCase(ctrl)
				mock.EXPECT().GetKolByID(gomock.Any(), kolID).Return(nil, errors.New("internal error"))

				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			mockUC := tt.getMock(ctrl)
			handler := NewKolHandler(mockUC)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/kols/"+tt.paramID, nil)
			c.Params = gin.Params{{Key: "id", Value: tt.paramID}}

			handler.GetKolByID(c)

			if w.Code != tt.expectedStatus {
				t.Errorf("GetKolByID() status = %d, want %d, body = %s", w.Code, tt.expectedStatus, w.Body.String())
			}

			if tt.checkBody != nil {
				tt.checkBody(t, w.Body.Bytes())
			}
		})
	}
}

func TestKolHandler_DeleteKolByID(t *testing.T) {
	t.Parallel()

	kolID := uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6")

	tests := []struct {
		name           string
		paramID        string
		expectedStatus int
		getMock        func(ctrl *gomock.Controller) *usecasemock.MockKolUseCase
	}{
		{
			name:           "success",
			paramID:        kolID.String(),
			expectedStatus: http.StatusOK,
			getMock: func(ctrl *gomock.Controller) *usecasemock.MockKolUseCase {
				mock := usecasemock.NewMockKolUseCase(ctrl)
				mock.EXPECT().DeleteKolByID(gomock.Any(), kolID).Return(nil)

				return mock
			},
		},
		{
			name:           "invalid_uuid",
			paramID:        "not-a-uuid",
			expectedStatus: http.StatusBadRequest,
			getMock: func(ctrl *gomock.Controller) *usecasemock.MockKolUseCase {
				return usecasemock.NewMockKolUseCase(ctrl)
			},
		},
		{
			name:           "usecase_error",
			paramID:        kolID.String(),
			expectedStatus: http.StatusInternalServerError,
			getMock: func(ctrl *gomock.Controller) *usecasemock.MockKolUseCase {
				mock := usecasemock.NewMockKolUseCase(ctrl)
				mock.EXPECT().DeleteKolByID(gomock.Any(), kolID).Return(errors.New("internal error"))

				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			mockUC := tt.getMock(ctrl)
			handler := NewKolHandler(mockUC)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodDelete, "/api/v1/kols/"+tt.paramID, nil)
			c.Params = gin.Params{{Key: "id", Value: tt.paramID}}

			handler.DeleteKolByID(c)

			if w.Code != tt.expectedStatus {
				t.Errorf("DeleteKolByID() status = %d, want %d, body = %s", w.Code, tt.expectedStatus, w.Body.String())
			}
		})
	}
}

func TestKolHandler_ListKols(t *testing.T) {
	t.Parallel()

	kolID := uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6")

	tests := []struct {
		name           string
		queryString    string
		expectedStatus int
		getMock        func(ctrl *gomock.Controller) *usecasemock.MockKolUseCase
		checkBody      func(t *testing.T, body []byte)
	}{
		{
			name:           "success",
			queryString:    "page_index=1&page_size=10",
			expectedStatus: http.StatusOK,
			getMock: func(ctrl *gomock.Controller) *usecasemock.MockKolUseCase {
				mock := usecasemock.NewMockKolUseCase(ctrl)
				mock.EXPECT().ListKols(gomock.Any(), usecase.ListKolsParam{
					TagIDs:   []uuid.UUID{},
					Page:     1,
					PageSize: 10,
				}).Return([]*usecase.Kol{
					{
						ID:    kolID,
						Name:  "Test",
						Email: "test@example.com",
						Sex:   kol.SexMale,
						Tags:  []usecase.Tag{},
					},
				}, 1, nil)

				return mock
			},
			checkBody: func(t *testing.T, body []byte) {
				var result ListKolsResponse
				if err := json.Unmarshal(body, &result); err != nil {
					t.Fatalf("failed to unmarshal response: %v", err)
				}

				if result.Total != 1 {
					t.Errorf("expected total 1, got %d", result.Total)
				}

				if len(result.Kols) != 1 {
					t.Errorf("expected 1 kol, got %d", len(result.Kols))
				}
			},
		},
		{
			name:           "missing_page_params_uses_defaults",
			queryString:    "",
			expectedStatus: http.StatusOK,
			getMock: func(ctrl *gomock.Controller) *usecasemock.MockKolUseCase {
				mock := usecasemock.NewMockKolUseCase(ctrl)
				mock.EXPECT().ListKols(gomock.Any(), gomock.Any()).Return([]*usecase.Kol{}, 0, nil)

				return mock
			},
		},
		{
			name:           "usecase_error",
			queryString:    "page_index=1&page_size=10",
			expectedStatus: http.StatusInternalServerError,
			getMock: func(ctrl *gomock.Controller) *usecasemock.MockKolUseCase {
				mock := usecasemock.NewMockKolUseCase(ctrl)
				mock.EXPECT().ListKols(gomock.Any(), gomock.Any()).Return(nil, 0, errors.New("internal error"))

				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			mockUC := tt.getMock(ctrl)
			handler := NewKolHandler(mockUC)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			url := "/api/v1/kols"
			if tt.queryString != "" {
				url += "?" + tt.queryString
			}

			c.Request = httptest.NewRequest(http.MethodGet, url, nil)

			handler.ListKols(c)

			if w.Code != tt.expectedStatus {
				t.Errorf("ListKols() status = %d, want %d, body = %s", w.Code, tt.expectedStatus, w.Body.String())
			}

			if tt.checkBody != nil {
				tt.checkBody(t, w.Body.Bytes())
			}
		})
	}
}

func TestKolHandler_CreateTag(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		body           string
		expectedStatus int
		getMock        func(ctrl *gomock.Controller) *usecasemock.MockKolUseCase
	}{
		{
			name:           "success",
			body:           `{"name":"new-tag"}`,
			expectedStatus: http.StatusOK,
			getMock: func(ctrl *gomock.Controller) *usecasemock.MockKolUseCase {
				mock := usecasemock.NewMockKolUseCase(ctrl)
				mock.EXPECT().CreateTag(gomock.Any(), gomock.Any()).Return(nil)

				return mock
			},
		},
		{
			name:           "invalid_json",
			body:           `{invalid}`,
			expectedStatus: http.StatusBadRequest,
			getMock: func(ctrl *gomock.Controller) *usecasemock.MockKolUseCase {
				return usecasemock.NewMockKolUseCase(ctrl)
			},
		},
		{
			name:           "missing_name",
			body:           `{}`,
			expectedStatus: http.StatusBadRequest,
			getMock: func(ctrl *gomock.Controller) *usecasemock.MockKolUseCase {
				return usecasemock.NewMockKolUseCase(ctrl)
			},
		},
		{
			name:           "usecase_error",
			body:           `{"name":"new-tag"}`,
			expectedStatus: http.StatusInternalServerError,
			getMock: func(ctrl *gomock.Controller) *usecasemock.MockKolUseCase {
				mock := usecasemock.NewMockKolUseCase(ctrl)
				mock.EXPECT().CreateTag(gomock.Any(), gomock.Any()).Return(errors.New("internal error"))

				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			mockUC := tt.getMock(ctrl)
			handler := NewKolHandler(mockUC)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/tags", strings.NewReader(tt.body))
			c.Request.Header.Set("Content-Type", "application/json")

			handler.CreateTag(c)

			if w.Code != tt.expectedStatus {
				t.Errorf("CreateTag() status = %d, want %d, body = %s", w.Code, tt.expectedStatus, w.Body.String())
			}
		})
	}
}

func TestKolHandler_ListTags(t *testing.T) {
	t.Parallel()

	tagID := uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6")
	now := time.Now().Truncate(time.Second)

	tests := []struct {
		name           string
		queryName      string
		expectedStatus int
		getMock        func(ctrl *gomock.Controller) *usecasemock.MockKolUseCase
		checkBody      func(t *testing.T, body []byte)
	}{
		{
			name:           "success",
			queryName:      "test",
			expectedStatus: http.StatusOK,
			getMock: func(ctrl *gomock.Controller) *usecasemock.MockKolUseCase {
				mock := usecasemock.NewMockKolUseCase(ctrl)
				mock.EXPECT().ListTagsByName(gomock.Any(), "test").Return([]*usecase.Tag{
					{ID: tagID, Name: "test-tag", CreatedAt: now},
				}, nil)

				return mock
			},
			checkBody: func(t *testing.T, body []byte) {
				var result []*usecase.Tag
				if err := json.Unmarshal(body, &result); err != nil {
					t.Fatalf("failed to unmarshal response: %v", err)
				}

				if len(result) != 1 {
					t.Errorf("expected 1 tag, got %d", len(result))
				}

				if result[0].Name != "test-tag" {
					t.Errorf("expected tag name %q, got %q", "test-tag", result[0].Name)
				}
			},
		},
		{
			name:           "empty_query",
			queryName:      "",
			expectedStatus: http.StatusOK,
			getMock: func(ctrl *gomock.Controller) *usecasemock.MockKolUseCase {
				mock := usecasemock.NewMockKolUseCase(ctrl)
				mock.EXPECT().ListTagsByName(gomock.Any(), "").Return([]*usecase.Tag{}, nil)

				return mock
			},
		},
		{
			name:           "usecase_error",
			queryName:      "test",
			expectedStatus: http.StatusInternalServerError,
			getMock: func(ctrl *gomock.Controller) *usecasemock.MockKolUseCase {
				mock := usecasemock.NewMockKolUseCase(ctrl)
				mock.EXPECT().ListTagsByName(gomock.Any(), "test").Return(nil, errors.New("internal error"))

				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			mockUC := tt.getMock(ctrl)
			handler := NewKolHandler(mockUC)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			url := "/api/v1/tags"
			if tt.queryName != "" {
				url += "?name=" + tt.queryName
			}

			c.Request = httptest.NewRequest(http.MethodGet, url, nil)

			handler.ListTags(c)

			if w.Code != tt.expectedStatus {
				t.Errorf("ListTags() status = %d, want %d, body = %s", w.Code, tt.expectedStatus, w.Body.String())
			}

			if tt.checkBody != nil {
				tt.checkBody(t, w.Body.Bytes())
			}
		})
	}
}

func TestKolHandler_CreateProduct(t *testing.T) {
	t.Parallel()

	adminID := uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6")

	tests := []struct {
		name           string
		body           string
		expectedStatus int
		getMock        func(ctrl *gomock.Controller) *usecasemock.MockKolUseCase
	}{
		{
			name:           "success",
			body:           `{"name":"product-1","description":"a product"}`,
			expectedStatus: http.StatusOK,
			getMock: func(ctrl *gomock.Controller) *usecasemock.MockKolUseCase {
				mock := usecasemock.NewMockKolUseCase(ctrl)
				mock.EXPECT().CreateProduct(gomock.Any(), usecase.CreateProductParam{
					Name:           "product-1",
					Description:    "a product",
					UpdatedAdminID: adminID,
				}).Return(nil)

				return mock
			},
		},
		{
			name:           "invalid_json",
			body:           `{invalid}`,
			expectedStatus: http.StatusBadRequest,
			getMock: func(ctrl *gomock.Controller) *usecasemock.MockKolUseCase {
				return usecasemock.NewMockKolUseCase(ctrl)
			},
		},
		{
			name:           "missing_name",
			body:           `{"description":"only desc"}`,
			expectedStatus: http.StatusBadRequest,
			getMock: func(ctrl *gomock.Controller) *usecasemock.MockKolUseCase {
				return usecasemock.NewMockKolUseCase(ctrl)
			},
		},
		{
			name:           "usecase_error",
			body:           `{"name":"product-1","description":"a product"}`,
			expectedStatus: http.StatusInternalServerError,
			getMock: func(ctrl *gomock.Controller) *usecasemock.MockKolUseCase {
				mock := usecasemock.NewMockKolUseCase(ctrl)
				mock.EXPECT().CreateProduct(gomock.Any(), gomock.Any()).Return(errors.New("internal error"))

				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			mockUC := tt.getMock(ctrl)
			handler := NewKolHandler(mockUC)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/products", strings.NewReader(tt.body))
			c.Request.Header.Set("Content-Type", "application/json")
			c.Set(admin.AdminIDKey, adminID)

			handler.CreateProduct(c)

			if w.Code != tt.expectedStatus {
				t.Errorf("CreateProduct() status = %d, want %d, body = %s", w.Code, tt.expectedStatus, w.Body.String())
			}
		})
	}
}

func TestKolHandler_UpdateKol(t *testing.T) {
	t.Parallel()

	adminID := uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6")
	kolID := uuid.MustParse("1193487b-f1a2-7a72-8ae4-197b84dc52d6")
	tagID := uuid.MustParse("2193487b-f1a2-7a72-8ae4-197b84dc52d6")

	tests := []struct {
		name           string
		body           string
		paramID        string
		expectedStatus int
		getMock        func(ctrl *gomock.Controller) *usecasemock.MockKolUseCase
	}{
		{
			name:    "success",
			body:    `{"name":"updated","email":"up@example.com","description":"new desc","social_media":"tw","sex":"f","tags":["2193487b-f1a2-7a72-8ae4-197b84dc52d6"]}`,
			paramID: kolID.String(),
			expectedStatus: http.StatusOK,
			getMock: func(ctrl *gomock.Controller) *usecasemock.MockKolUseCase {
				mock := usecasemock.NewMockKolUseCase(ctrl)
				mock.EXPECT().UpdateKol(gomock.Any(), usecase.UpdateKolParam{
					KolID:          kolID,
					Name:           "updated",
					Email:          "up@example.com",
					Description:    "new desc",
					SocialMedia:    "tw",
					Sex:            kol.SexFemale,
					Tags:           []uuid.UUID{tagID},
					UpdatedAdminID: adminID,
				}).Return(nil)

				return mock
			},
		},
		{
			name:           "invalid_json",
			body:           `{invalid}`,
			paramID:        kolID.String(),
			expectedStatus: http.StatusBadRequest,
			getMock: func(ctrl *gomock.Controller) *usecasemock.MockKolUseCase {
				return usecasemock.NewMockKolUseCase(ctrl)
			},
		},
		{
			name:           "invalid_path_param_uuid",
			body:           `{"name":"updated","email":"up@example.com","sex":"m"}`,
			paramID:        "not-a-uuid",
			expectedStatus: http.StatusBadRequest,
			getMock: func(ctrl *gomock.Controller) *usecasemock.MockKolUseCase {
				return usecasemock.NewMockKolUseCase(ctrl)
			},
		},
		{
			name:    "usecase_error",
			body:    `{"name":"updated","email":"up@example.com","sex":"m"}`,
			paramID: kolID.String(),
			expectedStatus: http.StatusInternalServerError,
			getMock: func(ctrl *gomock.Controller) *usecasemock.MockKolUseCase {
				mock := usecasemock.NewMockKolUseCase(ctrl)
				mock.EXPECT().UpdateKol(gomock.Any(), gomock.Any()).Return(errors.New("internal error"))

				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			mockUC := tt.getMock(ctrl)
			handler := NewKolHandler(mockUC)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodPut, "/api/v1/kols/"+tt.paramID, strings.NewReader(tt.body))
			c.Request.Header.Set("Content-Type", "application/json")
			c.Params = gin.Params{{Key: "id", Value: tt.paramID}}
			c.Set(admin.AdminIDKey, adminID)

			handler.UpdateKol(c)

			if w.Code != tt.expectedStatus {
				t.Errorf("UpdateKol() status = %d, want %d, body = %s", w.Code, tt.expectedStatus, w.Body.String())
			}
		})
	}
}

func TestKolHandler_BatchCreateKolsByXlsx(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		expectedStatus int
		getMock        func(ctrl *gomock.Controller) *usecasemock.MockKolUseCase
	}{
		{
			name:           "bad_request_no_file",
			expectedStatus: http.StatusBadRequest,
			getMock: func(ctrl *gomock.Controller) *usecasemock.MockKolUseCase {
				return usecasemock.NewMockKolUseCase(ctrl)
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			mockUC := tt.getMock(ctrl)
			handler := NewKolHandler(mockUC)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/kols/upload", nil)
			c.Request.Header.Set("Content-Type", "multipart/form-data")

			handler.BatchCreateKolsByXlsx(c)

			if w.Code != tt.expectedStatus {
				t.Errorf("BatchCreateKolsByXlsx() status = %d, want %d, body = %s", w.Code, tt.expectedStatus, w.Body.String())
			}
		})
	}
}

func TestKolHandler_ListKols_WithTagIDs(t *testing.T) {
	t.Parallel()

	kolID := uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6")
	tagID := uuid.MustParse("1193487b-f1a2-7a72-8ae4-197b84dc52d6")

	tests := []struct {
		name           string
		queryString    string
		expectedStatus int
		getMock        func(ctrl *gomock.Controller) *usecasemock.MockKolUseCase
	}{
		{
			name:           "with_tag_ids",
			queryString:    "page_index=1&page_size=10&tag_ids[]=1193487b-f1a2-7a72-8ae4-197b84dc52d6",
			expectedStatus: http.StatusOK,
			getMock: func(ctrl *gomock.Controller) *usecasemock.MockKolUseCase {
				mock := usecasemock.NewMockKolUseCase(ctrl)
				mock.EXPECT().ListKols(gomock.Any(), usecase.ListKolsParam{
					TagIDs:   []uuid.UUID{tagID},
					Page:     1,
					PageSize: 10,
				}).Return([]*usecase.Kol{
					{
						ID:    kolID,
						Name:  "Test",
						Email: "test@example.com",
						Sex:   kol.SexMale,
						Tags:  []usecase.Tag{},
					},
				}, 1, nil)

				return mock
			},
		},
		{
			name:           "with_invalid_tag_id",
			queryString:    "page_index=1&page_size=10&tag_ids[]=not-a-uuid",
			expectedStatus: http.StatusBadRequest,
			getMock: func(ctrl *gomock.Controller) *usecasemock.MockKolUseCase {
				return usecasemock.NewMockKolUseCase(ctrl)
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			mockUC := tt.getMock(ctrl)
			handler := NewKolHandler(mockUC)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			url := "/api/v1/kols?" + tt.queryString
			c.Request = httptest.NewRequest(http.MethodGet, url, nil)

			handler.ListKols(c)

			if w.Code != tt.expectedStatus {
				t.Errorf("ListKols() status = %d, want %d, body = %s", w.Code, tt.expectedStatus, w.Body.String())
			}
		})
	}
}

func TestKolHandler_ListProducts(t *testing.T) {
	t.Parallel()

	productID := uuid.MustParse("0193487b-f1a2-7a72-8ae4-197b84dc52d6")
	now := time.Now().Truncate(time.Second)

	tests := []struct {
		name           string
		queryName      string
		expectedStatus int
		getMock        func(ctrl *gomock.Controller) *usecasemock.MockKolUseCase
		checkBody      func(t *testing.T, body []byte)
	}{
		{
			name:           "success",
			queryName:      "test",
			expectedStatus: http.StatusOK,
			getMock: func(ctrl *gomock.Controller) *usecasemock.MockKolUseCase {
				mock := usecasemock.NewMockKolUseCase(ctrl)
				mock.EXPECT().ListProductsByName(gomock.Any(), "test").Return([]*usecase.Product{
					{ID: productID, Name: "test-product", Description: "desc", CreatedAt: now},
				}, nil)

				return mock
			},
			checkBody: func(t *testing.T, body []byte) {
				var result []*usecase.Product
				if err := json.Unmarshal(body, &result); err != nil {
					t.Fatalf("failed to unmarshal response: %v", err)
				}

				if len(result) != 1 {
					t.Errorf("expected 1 product, got %d", len(result))
				}

				if result[0].Name != "test-product" {
					t.Errorf("expected product name %q, got %q", "test-product", result[0].Name)
				}
			},
		},
		{
			name:           "empty_query",
			queryName:      "",
			expectedStatus: http.StatusOK,
			getMock: func(ctrl *gomock.Controller) *usecasemock.MockKolUseCase {
				mock := usecasemock.NewMockKolUseCase(ctrl)
				mock.EXPECT().ListProductsByName(gomock.Any(), "").Return([]*usecase.Product{}, nil)

				return mock
			},
		},
		{
			name:           "usecase_error",
			queryName:      "test",
			expectedStatus: http.StatusInternalServerError,
			getMock: func(ctrl *gomock.Controller) *usecasemock.MockKolUseCase {
				mock := usecasemock.NewMockKolUseCase(ctrl)
				mock.EXPECT().ListProductsByName(gomock.Any(), "test").Return(nil, errors.New("internal error"))

				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			mockUC := tt.getMock(ctrl)
			handler := NewKolHandler(mockUC)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			url := "/api/v1/products"
			if tt.queryName != "" {
				url += "?name=" + tt.queryName
			}

			c.Request = httptest.NewRequest(http.MethodGet, url, nil)

			handler.ListProducts(c)

			if w.Code != tt.expectedStatus {
				t.Errorf("ListProducts() status = %d, want %d, body = %s", w.Code, tt.expectedStatus, w.Body.String())
			}

			if tt.checkBody != nil {
				tt.checkBody(t, w.Body.Bytes())
			}
		})
	}
}
