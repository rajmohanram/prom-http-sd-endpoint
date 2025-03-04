package server

import (
    "github.com/stretchr/testify/mock"
)

type MockServer struct {
    mock.Mock
}

func (m *MockServer) Start(address string) error {
    args := m.Called(address)
    return args.Error(0)
}
