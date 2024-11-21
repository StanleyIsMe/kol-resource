// Code generated by MockGen. DO NOT EDIT.
// Source: ./repository.go

// Package repositorymock is a generated GoMock package.
package repositorymock

import (
	context "context"
	domain "kolresource/internal/kol/domain"
	entities "kolresource/internal/kol/domain/entities"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
)

// MockRepository is a mock of Repository interface.
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository.
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance.
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// CreateKol mocks base method.
func (m *MockRepository) CreateKol(ctx context.Context, param domain.CreateKolParams) (*entities.Kol, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateKol", ctx, param)
	ret0, _ := ret[0].(*entities.Kol)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateKol indicates an expected call of CreateKol.
func (mr *MockRepositoryMockRecorder) CreateKol(ctx, param interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateKol", reflect.TypeOf((*MockRepository)(nil).CreateKol), ctx, param)
}

// CreateProduct mocks base method.
func (m *MockRepository) CreateProduct(ctx context.Context, param domain.CreateProductParams) (*entities.Product, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateProduct", ctx, param)
	ret0, _ := ret[0].(*entities.Product)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateProduct indicates an expected call of CreateProduct.
func (mr *MockRepositoryMockRecorder) CreateProduct(ctx, param interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateProduct", reflect.TypeOf((*MockRepository)(nil).CreateProduct), ctx, param)
}

// CreateSendEmailLog mocks base method.
func (m *MockRepository) CreateSendEmailLog(ctx context.Context, sendEmailLog *entities.SendEmailLog) (*entities.SendEmailLog, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateSendEmailLog", ctx, sendEmailLog)
	ret0, _ := ret[0].(*entities.SendEmailLog)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateSendEmailLog indicates an expected call of CreateSendEmailLog.
func (mr *MockRepositoryMockRecorder) CreateSendEmailLog(ctx, sendEmailLog interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSendEmailLog", reflect.TypeOf((*MockRepository)(nil).CreateSendEmailLog), ctx, sendEmailLog)
}

// CreateTag mocks base method.
func (m *MockRepository) CreateTag(ctx context.Context, param domain.CreateTagParams) (*entities.Tag, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTag", ctx, param)
	ret0, _ := ret[0].(*entities.Tag)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateTag indicates an expected call of CreateTag.
func (mr *MockRepositoryMockRecorder) CreateTag(ctx, param interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTag", reflect.TypeOf((*MockRepository)(nil).CreateTag), ctx, param)
}

// DeleteKolByID mocks base method.
func (m *MockRepository) DeleteKolByID(ctx context.Context, id uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteKolByID", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteKolByID indicates an expected call of DeleteKolByID.
func (mr *MockRepositoryMockRecorder) DeleteKolByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteKolByID", reflect.TypeOf((*MockRepository)(nil).DeleteKolByID), ctx, id)
}

// DeleteProductByID mocks base method.
func (m *MockRepository) DeleteProductByID(ctx context.Context, id uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteProductByID", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteProductByID indicates an expected call of DeleteProductByID.
func (mr *MockRepositoryMockRecorder) DeleteProductByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteProductByID", reflect.TypeOf((*MockRepository)(nil).DeleteProductByID), ctx, id)
}

// DeleteTagByID mocks base method.
func (m *MockRepository) DeleteTagByID(ctx context.Context, id uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteTagByID", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteTagByID indicates an expected call of DeleteTagByID.
func (mr *MockRepositoryMockRecorder) DeleteTagByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteTagByID", reflect.TypeOf((*MockRepository)(nil).DeleteTagByID), ctx, id)
}

// GetKolByEmail mocks base method.
func (m *MockRepository) GetKolByEmail(ctx context.Context, email string) (*entities.Kol, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetKolByEmail", ctx, email)
	ret0, _ := ret[0].(*entities.Kol)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetKolByEmail indicates an expected call of GetKolByEmail.
func (mr *MockRepositoryMockRecorder) GetKolByEmail(ctx, email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetKolByEmail", reflect.TypeOf((*MockRepository)(nil).GetKolByEmail), ctx, email)
}

// GetKolByID mocks base method.
func (m *MockRepository) GetKolByID(ctx context.Context, id uuid.UUID) (*entities.Kol, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetKolByID", ctx, id)
	ret0, _ := ret[0].(*entities.Kol)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetKolByID indicates an expected call of GetKolByID.
func (mr *MockRepositoryMockRecorder) GetKolByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetKolByID", reflect.TypeOf((*MockRepository)(nil).GetKolByID), ctx, id)
}

// GetKolWithTagsByID mocks base method.
func (m *MockRepository) GetKolWithTagsByID(ctx context.Context, id uuid.UUID) (*domain.Kol, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetKolWithTagsByID", ctx, id)
	ret0, _ := ret[0].(*domain.Kol)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetKolWithTagsByID indicates an expected call of GetKolWithTagsByID.
func (mr *MockRepositoryMockRecorder) GetKolWithTagsByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetKolWithTagsByID", reflect.TypeOf((*MockRepository)(nil).GetKolWithTagsByID), ctx, id)
}

// GetProductByID mocks base method.
func (m *MockRepository) GetProductByID(ctx context.Context, id uuid.UUID) (*entities.Product, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProductByID", ctx, id)
	ret0, _ := ret[0].(*entities.Product)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProductByID indicates an expected call of GetProductByID.
func (mr *MockRepositoryMockRecorder) GetProductByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProductByID", reflect.TypeOf((*MockRepository)(nil).GetProductByID), ctx, id)
}

// GetProductByName mocks base method.
func (m *MockRepository) GetProductByName(ctx context.Context, name string) (*entities.Product, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProductByName", ctx, name)
	ret0, _ := ret[0].(*entities.Product)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProductByName indicates an expected call of GetProductByName.
func (mr *MockRepositoryMockRecorder) GetProductByName(ctx, name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProductByName", reflect.TypeOf((*MockRepository)(nil).GetProductByName), ctx, name)
}

// GetTagByName mocks base method.
func (m *MockRepository) GetTagByName(ctx context.Context, name string) (*entities.Tag, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTagByName", ctx, name)
	ret0, _ := ret[0].(*entities.Tag)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTagByName indicates an expected call of GetTagByName.
func (mr *MockRepositoryMockRecorder) GetTagByName(ctx, name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTagByName", reflect.TypeOf((*MockRepository)(nil).GetTagByName), ctx, name)
}

// ListKolWithTagsByFilters mocks base method.
func (m *MockRepository) ListKolWithTagsByFilters(ctx context.Context, param domain.ListKolWithTagsByFiltersParams) ([]*domain.Kol, int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListKolWithTagsByFilters", ctx, param)
	ret0, _ := ret[0].([]*domain.Kol)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ListKolWithTagsByFilters indicates an expected call of ListKolWithTagsByFilters.
func (mr *MockRepositoryMockRecorder) ListKolWithTagsByFilters(ctx, param interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListKolWithTagsByFilters", reflect.TypeOf((*MockRepository)(nil).ListKolWithTagsByFilters), ctx, param)
}

// ListKolsByIDs mocks base method.
func (m *MockRepository) ListKolsByIDs(ctx context.Context, ids []uuid.UUID) ([]*entities.Kol, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListKolsByIDs", ctx, ids)
	ret0, _ := ret[0].([]*entities.Kol)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListKolsByIDs indicates an expected call of ListKolsByIDs.
func (mr *MockRepositoryMockRecorder) ListKolsByIDs(ctx, ids interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListKolsByIDs", reflect.TypeOf((*MockRepository)(nil).ListKolsByIDs), ctx, ids)
}

// ListProductsByName mocks base method.
func (m *MockRepository) ListProductsByName(ctx context.Context, name string) ([]*entities.Product, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListProductsByName", ctx, name)
	ret0, _ := ret[0].([]*entities.Product)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListProductsByName indicates an expected call of ListProductsByName.
func (mr *MockRepositoryMockRecorder) ListProductsByName(ctx, name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListProductsByName", reflect.TypeOf((*MockRepository)(nil).ListProductsByName), ctx, name)
}

// ListSendEmailLogsByFilter mocks base method.
func (m *MockRepository) ListSendEmailLogsByFilter(ctx context.Context, param domain.ListSendEmailLogsByFilterParams) ([]*entities.SendEmailLog, int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListSendEmailLogsByFilter", ctx, param)
	ret0, _ := ret[0].([]*entities.SendEmailLog)
	ret1, _ := ret[1].(int)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ListSendEmailLogsByFilter indicates an expected call of ListSendEmailLogsByFilter.
func (mr *MockRepositoryMockRecorder) ListSendEmailLogsByFilter(ctx, param interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListSendEmailLogsByFilter", reflect.TypeOf((*MockRepository)(nil).ListSendEmailLogsByFilter), ctx, param)
}

// ListTagsByName mocks base method.
func (m *MockRepository) ListTagsByName(ctx context.Context, name string) ([]*entities.Tag, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListTagsByName", ctx, name)
	ret0, _ := ret[0].([]*entities.Tag)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListTagsByName indicates an expected call of ListTagsByName.
func (mr *MockRepositoryMockRecorder) ListTagsByName(ctx, name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListTagsByName", reflect.TypeOf((*MockRepository)(nil).ListTagsByName), ctx, name)
}

// UpdateKol mocks base method.
func (m *MockRepository) UpdateKol(ctx context.Context, param domain.UpdateKolParams) (*entities.Kol, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateKol", ctx, param)
	ret0, _ := ret[0].(*entities.Kol)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateKol indicates an expected call of UpdateKol.
func (mr *MockRepositoryMockRecorder) UpdateKol(ctx, param interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateKol", reflect.TypeOf((*MockRepository)(nil).UpdateKol), ctx, param)
}

// MockEmailRepository is a mock of EmailRepository interface.
type MockEmailRepository struct {
	ctrl     *gomock.Controller
	recorder *MockEmailRepositoryMockRecorder
}

// MockEmailRepositoryMockRecorder is the mock recorder for MockEmailRepository.
type MockEmailRepositoryMockRecorder struct {
	mock *MockEmailRepository
}

// NewMockEmailRepository creates a new mock instance.
func NewMockEmailRepository(ctrl *gomock.Controller) *MockEmailRepository {
	mock := &MockEmailRepository{ctrl: ctrl}
	mock.recorder = &MockEmailRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEmailRepository) EXPECT() *MockEmailRepositoryMockRecorder {
	return m.recorder
}

// SendEmail mocks base method.
func (m *MockEmailRepository) SendEmail(ctx context.Context, param domain.SendEmailParams) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendEmail", ctx, param)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendEmail indicates an expected call of SendEmail.
func (mr *MockEmailRepositoryMockRecorder) SendEmail(ctx, param interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendEmail", reflect.TypeOf((*MockEmailRepository)(nil).SendEmail), ctx, param)
}