// Code generated by MockGen. DO NOT EDIT.
// Source: sendable.go

// Package mock_messenger is a generated GoMock package.
package mock_messenger

import (
	reflect "reflect"

	discordgo "github.com/bwmarrin/discordgo"
	gomock "github.com/golang/mock/gomock"
	messenger "github.com/oppzippy/BoostRequestBot/boost_request/messenger"
)

// MockSendable is a mock of Sendable interface.
type MockSendable struct {
	ctrl     *gomock.Controller
	recorder *MockSendableMockRecorder
}

// MockSendableMockRecorder is the mock recorder for MockSendable.
type MockSendableMockRecorder struct {
	mock *MockSendable
}

// NewMockSendable creates a new mock instance.
func NewMockSendable(ctrl *gomock.Controller) *MockSendable {
	mock := &MockSendable{ctrl: ctrl}
	mock.recorder = &MockSendableMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSendable) EXPECT() *MockSendableMockRecorder {
	return m.recorder
}

// Send mocks base method.
func (m *MockSendable) Send(discord messenger.DiscordSender) (*discordgo.Message, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Send", discord)
	ret0, _ := ret[0].(*discordgo.Message)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Send indicates an expected call of Send.
func (mr *MockSendableMockRecorder) Send(discord interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Send", reflect.TypeOf((*MockSendable)(nil).Send), discord)
}
