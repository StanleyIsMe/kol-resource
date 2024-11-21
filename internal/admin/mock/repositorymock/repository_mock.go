// Code generated by MockGen. DO NOT EDIT.
// Source: ./repository.go

// Package repositorymock is a generated GoMock package.
package repositorymock

import (
	context "context"
	entities "kolresource/internal/admin/domain/entities"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
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

// CreateAdmin mocks base method.
func (m *MockRepository) CreateAdmin(ctx context.Context, admin *entities.Admin) (*entities.Admin, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateAdmin", ctx, admin)
	ret0, _ := ret[0].(*entities.Admin)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateAdmin indicates an expected call of CreateAdmin.
func (mr *MockRepositoryMockRecorder) CreateAdmin(ctx, admin interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateAdmin", reflect.TypeOf((*MockRepository)(nil).CreateAdmin), ctx, admin)
}

// GetAdminByUserName mocks base method.
func (m *MockRepository) GetAdminByUserName(ctx context.Context, userName string) (*entities.Admin, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAdminByUserName", ctx, userName)
	ret0, _ := ret[0].(*entities.Admin)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAdminByUserName indicates an expected call of GetAdminByUserName.
func (mr *MockRepositoryMockRecorder) GetAdminByUserName(ctx, userName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAdminByUserName", reflect.TypeOf((*MockRepository)(nil).GetAdminByUserName), ctx, userName)
}
