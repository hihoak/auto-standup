// Code generated by MockGen. DO NOT EDIT.
// Source: ./filters.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	jira "github.com/andygrunwald/go-jira"
	gomock "github.com/golang/mock/gomock"
	utils "github.com/hihoak/auto-standup/pkg/utils"
)

// MockFilterers is a mock of Filterers interface.
type MockFilterers struct {
	ctrl     *gomock.Controller
	recorder *MockFilterersMockRecorder
}

// MockFilterersMockRecorder is the mock recorder for MockFilterers.
type MockFilterersMockRecorder struct {
	mock *MockFilterers
}

// NewMockFilterers creates a new mock instance.
func NewMockFilterers(ctrl *gomock.Controller) *MockFilterers {
	mock := &MockFilterers{ctrl: ctrl}
	mock.recorder = &MockFilterersMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFilterers) EXPECT() *MockFilterersMockRecorder {
	return m.recorder
}

// FilterIssueByActivity mocks base method.
func (m *MockFilterers) FilterIssueByActivity(cfg *utils.Config, issue jira.Issue) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FilterIssueByActivity", cfg, issue)
	ret0, _ := ret[0].(bool)
	return ret0
}

// FilterIssueByActivity indicates an expected call of FilterIssueByActivity.
func (mr *MockFilterersMockRecorder) FilterIssueByActivity(cfg, issue interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FilterIssueByActivity", reflect.TypeOf((*MockFilterers)(nil).FilterIssueByActivity), cfg, issue)
}

// FilterIssuesByProject mocks base method.
func (m *MockFilterers) FilterIssuesByProject(cfg *utils.Config, issue jira.Issue) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FilterIssuesByProject", cfg, issue)
	ret0, _ := ret[0].(bool)
	return ret0
}

// FilterIssuesByProject indicates an expected call of FilterIssuesByProject.
func (mr *MockFilterersMockRecorder) FilterIssuesByProject(cfg, issue interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FilterIssuesByProject", reflect.TypeOf((*MockFilterers)(nil).FilterIssuesByProject), cfg, issue)
}
