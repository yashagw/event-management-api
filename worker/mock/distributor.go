// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/yashagw/event-management-api/worker (interfaces: TaskDistributor)

// Package mockwk is a generated GoMock package.
package mockwk

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	asynq "github.com/hibiken/asynq"
	worker "github.com/yashagw/event-management-api/worker"
)

// MockTaskDistributor is a mock of TaskDistributor interface.
type MockTaskDistributor struct {
	ctrl     *gomock.Controller
	recorder *MockTaskDistributorMockRecorder
}

// MockTaskDistributorMockRecorder is the mock recorder for MockTaskDistributor.
type MockTaskDistributorMockRecorder struct {
	mock *MockTaskDistributor
}

// NewMockTaskDistributor creates a new mock instance.
func NewMockTaskDistributor(ctrl *gomock.Controller) *MockTaskDistributor {
	mock := &MockTaskDistributor{ctrl: ctrl}
	mock.recorder = &MockTaskDistributorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTaskDistributor) EXPECT() *MockTaskDistributorMockRecorder {
	return m.recorder
}

// DistributeTaskSendEmailVerify mocks base method.
func (m *MockTaskDistributor) DistributeTaskSendEmailVerify(arg0 context.Context, arg1 *worker.PayloadSendEmailVerify, arg2 ...asynq.Option) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DistributeTaskSendEmailVerify", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// DistributeTaskSendEmailVerify indicates an expected call of DistributeTaskSendEmailVerify.
func (mr *MockTaskDistributorMockRecorder) DistributeTaskSendEmailVerify(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DistributeTaskSendEmailVerify", reflect.TypeOf((*MockTaskDistributor)(nil).DistributeTaskSendEmailVerify), varargs...)
}
