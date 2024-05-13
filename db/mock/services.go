// Code generated by MockGen. DO NOT EDIT.
// Source: Simple-Bank/db/services (interfaces: Services)
//
// Generated by this command:
//
//	mockgen -package mockdb -destination db/mock/services.go Simple-Bank/db/services Services
//

// Package mockdb is a generated GoMock package.
package mockdb

import (
	models "Simple-Bank/db/models"
	services "Simple-Bank/db/services"
	requests "Simple-Bank/requests"
	reflect "reflect"

	uuid "github.com/google/uuid"
	gomock "go.uber.org/mock/gomock"
)

// MockServices is a mock of Services interface.
type MockServices struct {
	ctrl     *gomock.Controller
	recorder *MockServicesMockRecorder
}

// MockServicesMockRecorder is the mock recorder for MockServices.
type MockServicesMockRecorder struct {
	mock *MockServices
}

// NewMockServices creates a new mock instance.
func NewMockServices(ctrl *gomock.Controller) *MockServices {
	mock := &MockServices{ctrl: ctrl}
	mock.recorder = &MockServicesMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockServices) EXPECT() *MockServicesMockRecorder {
	return m.recorder
}

// CreateAccount mocks base method.
func (m *MockServices) CreateAccount(arg0 string) (models.Account, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateAccount", arg0)
	ret0, _ := ret[0].(models.Account)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateAccount indicates an expected call of CreateAccount.
func (mr *MockServicesMockRecorder) CreateAccount(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateAccount", reflect.TypeOf((*MockServices)(nil).CreateAccount), arg0)
}

// CreateSession mocks base method.
func (m *MockServices) CreateSession(arg0 models.Session) (models.Session, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateSession", arg0)
	ret0, _ := ret[0].(models.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateSession indicates an expected call of CreateSession.
func (mr *MockServicesMockRecorder) CreateSession(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSession", reflect.TypeOf((*MockServices)(nil).CreateSession), arg0)
}

// CreateUser mocks base method.
func (m *MockServices) CreateUser(arg0 requests.CreateUserRequest) (models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", arg0)
	ret0, _ := ret[0].(models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockServicesMockRecorder) CreateUser(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockServices)(nil).CreateUser), arg0)
}

// DeleteAccount mocks base method.
func (m *MockServices) DeleteAccount(arg0 int64) (models.Account, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteAccount", arg0)
	ret0, _ := ret[0].(models.Account)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteAccount indicates an expected call of DeleteAccount.
func (mr *MockServicesMockRecorder) DeleteAccount(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteAccount", reflect.TypeOf((*MockServices)(nil).DeleteAccount), arg0)
}

// DepositMoney mocks base method.
func (m *MockServices) DepositMoney(arg0 requests.DepositRequest) (models.Entry, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DepositMoney", arg0)
	ret0, _ := ret[0].(models.Entry)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DepositMoney indicates an expected call of DepositMoney.
func (mr *MockServicesMockRecorder) DepositMoney(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DepositMoney", reflect.TypeOf((*MockServices)(nil).DepositMoney), arg0)
}

// GetAccount mocks base method.
func (m *MockServices) GetAccount(arg0 int64) (models.Account, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAccount", arg0)
	ret0, _ := ret[0].(models.Account)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAccount indicates an expected call of GetAccount.
func (mr *MockServicesMockRecorder) GetAccount(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAccount", reflect.TypeOf((*MockServices)(nil).GetAccount), arg0)
}

// GetEntry mocks base method.
func (m *MockServices) GetEntry(arg0 int64) (models.Entry, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEntry", arg0)
	ret0, _ := ret[0].(models.Entry)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetEntry indicates an expected call of GetEntry.
func (mr *MockServicesMockRecorder) GetEntry(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEntry", reflect.TypeOf((*MockServices)(nil).GetEntry), arg0)
}

// GetSession mocks base method.
func (m *MockServices) GetSession(arg0 uuid.UUID) (models.Session, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSession", arg0)
	ret0, _ := ret[0].(models.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSession indicates an expected call of GetSession.
func (mr *MockServicesMockRecorder) GetSession(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSession", reflect.TypeOf((*MockServices)(nil).GetSession), arg0)
}

// GetTransfer mocks base method.
func (m *MockServices) GetTransfer(arg0 int64) (models.Transfer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTransfer", arg0)
	ret0, _ := ret[0].(models.Transfer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTransfer indicates an expected call of GetTransfer.
func (mr *MockServicesMockRecorder) GetTransfer(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTransfer", reflect.TypeOf((*MockServices)(nil).GetTransfer), arg0)
}

// GetUser mocks base method.
func (m *MockServices) GetUser(arg0 string) (models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUser", arg0)
	ret0, _ := ret[0].(models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUser indicates an expected call of GetUser.
func (mr *MockServicesMockRecorder) GetUser(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUser", reflect.TypeOf((*MockServices)(nil).GetUser), arg0)
}

// ListAccounts mocks base method.
func (m *MockServices) ListAccounts(arg0 services.ListAccountsRequest) ([]models.Account, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListAccounts", arg0)
	ret0, _ := ret[0].([]models.Account)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListAccounts indicates an expected call of ListAccounts.
func (mr *MockServicesMockRecorder) ListAccounts(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListAccounts", reflect.TypeOf((*MockServices)(nil).ListAccounts), arg0)
}

// Transfer mocks base method.
func (m *MockServices) Transfer(arg0 services.TransferRequest) (models.Transfer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Transfer", arg0)
	ret0, _ := ret[0].(models.Transfer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Transfer indicates an expected call of Transfer.
func (mr *MockServicesMockRecorder) Transfer(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Transfer", reflect.TypeOf((*MockServices)(nil).Transfer), arg0)
}

// UpdateUser mocks base method.
func (m *MockServices) UpdateUser(arg0 services.UpdateUserRequest) (models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateUser", arg0)
	ret0, _ := ret[0].(models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateUser indicates an expected call of UpdateUser.
func (mr *MockServicesMockRecorder) UpdateUser(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateUser", reflect.TypeOf((*MockServices)(nil).UpdateUser), arg0)
}

// WithdrawMoney mocks base method.
func (m *MockServices) WithdrawMoney(arg0 requests.WithdrawRequest) (models.Entry, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WithdrawMoney", arg0)
	ret0, _ := ret[0].(models.Entry)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// WithdrawMoney indicates an expected call of WithdrawMoney.
func (mr *MockServicesMockRecorder) WithdrawMoney(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WithdrawMoney", reflect.TypeOf((*MockServices)(nil).WithdrawMoney), arg0)
}
