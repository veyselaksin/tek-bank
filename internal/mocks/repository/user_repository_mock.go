// Code generated by MockGen. DO NOT EDIT.
// Source: tek-bank/internal/db/repository (interfaces: UserRepository)
//
// Generated by this command:
//
//	mockgen -destination=../../mocks/repository/user_repository_mock.go -package=repository tek-bank/internal/db/repository UserRepository
//

// Package repository is a generated GoMock package.
package repository

import (
	context "context"
	reflect "reflect"
	models "tek-bank/internal/db/models"
	repository "tek-bank/internal/db/repository"
	time "time"

	gomock "go.uber.org/mock/gomock"
	gorm "gorm.io/gorm"
)

// MockUserRepository is a mock of UserRepository interface.
type MockUserRepository struct {
	ctrl     *gomock.Controller
	recorder *MockUserRepositoryMockRecorder
}

// MockUserRepositoryMockRecorder is the mock recorder for MockUserRepository.
type MockUserRepositoryMockRecorder struct {
	mock *MockUserRepository
}

// NewMockUserRepository creates a new mock instance.
func NewMockUserRepository(ctrl *gomock.Controller) *MockUserRepository {
	mock := &MockUserRepository{ctrl: ctrl}
	mock.recorder = &MockUserRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserRepository) EXPECT() *MockUserRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockUserRepository) Create(arg0 models.User) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockUserRepositoryMockRecorder) Create(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockUserRepository)(nil).Create), arg0)
}

// FindAll mocks base method.
func (m *MockUserRepository) FindAll() ([]models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindAll")
	ret0, _ := ret[0].([]models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindAll indicates an expected call of FindAll.
func (mr *MockUserRepositoryMockRecorder) FindAll() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindAll", reflect.TypeOf((*MockUserRepository)(nil).FindAll))
}

// FindByEmail mocks base method.
func (m *MockUserRepository) FindByEmail(arg0 string) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByEmail", arg0)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByEmail indicates an expected call of FindByEmail.
func (mr *MockUserRepositoryMockRecorder) FindByEmail(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByEmail", reflect.TypeOf((*MockUserRepository)(nil).FindByEmail), arg0)
}

// FindByID mocks base method.
func (m *MockUserRepository) FindByID(arg0 string) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByID", arg0)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByID indicates an expected call of FindByID.
func (mr *MockUserRepositoryMockRecorder) FindByID(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByID", reflect.TypeOf((*MockUserRepository)(nil).FindByID), arg0)
}

// FindByUniqueIdentifier mocks base method.
func (m *MockUserRepository) FindByUniqueIdentifier(arg0 string) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByUniqueIdentifier", arg0)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByUniqueIdentifier indicates an expected call of FindByUniqueIdentifier.
func (mr *MockUserRepositoryMockRecorder) FindByUniqueIdentifier(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByUniqueIdentifier", reflect.TypeOf((*MockUserRepository)(nil).FindByUniqueIdentifier), arg0)
}

// GetTokenBlacklist mocks base method.
func (m *MockUserRepository) GetTokenBlacklist(arg0 *context.Context, arg1 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTokenBlacklist", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTokenBlacklist indicates an expected call of GetTokenBlacklist.
func (mr *MockUserRepositoryMockRecorder) GetTokenBlacklist(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTokenBlacklist", reflect.TypeOf((*MockUserRepository)(nil).GetTokenBlacklist), arg0, arg1)
}

// SetTokenBlacklist mocks base method.
func (m *MockUserRepository) SetTokenBlacklist(arg0 *context.Context, arg1, arg2 string, arg3 time.Duration) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetTokenBlacklist", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetTokenBlacklist indicates an expected call of SetTokenBlacklist.
func (mr *MockUserRepositoryMockRecorder) SetTokenBlacklist(arg0, arg1, arg2, arg3 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetTokenBlacklist", reflect.TypeOf((*MockUserRepository)(nil).SetTokenBlacklist), arg0, arg1, arg2, arg3)
}

// SoftDelete mocks base method.
func (m *MockUserRepository) SoftDelete(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SoftDelete", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SoftDelete indicates an expected call of SoftDelete.
func (mr *MockUserRepositoryMockRecorder) SoftDelete(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SoftDelete", reflect.TypeOf((*MockUserRepository)(nil).SoftDelete), arg0)
}

// WithTx mocks base method.
func (m *MockUserRepository) WithTx(arg0 *gorm.DB) repository.UserRepository {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WithTx", arg0)
	ret0, _ := ret[0].(repository.UserRepository)
	return ret0
}

// WithTx indicates an expected call of WithTx.
func (mr *MockUserRepositoryMockRecorder) WithTx(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WithTx", reflect.TypeOf((*MockUserRepository)(nil).WithTx), arg0)
}