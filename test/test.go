package test

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	jira_mocks "github.com/hihoak/auto-standup/internal/clients/jirer/mocks"
	"github.com/hihoak/auto-standup/internal/filters/mocks"
	"github.com/hihoak/auto-standup/pkg/utils"
	"github.com/stretchr/testify/require"
)

// ErrorTest - ...
var ErrorTest = errors.New("test error")

type (
	// Case - ...
	Case struct {
		Name           string
		FuncArguments  []interface{}
		Setup          func() (*MockClients, *utils.Config)
		ExpectedError  error
		ExpectedResult interface{}
	}

	// MockClients - ...
	MockClients struct {
		JiraMockClient *jira_mocks.MockJirer
		FiltersMock    *mocks.MockFilterers
	}
)

// InitDefaultMockClients - .,,
func InitDefaultMockClients(mc *gomock.Controller) *MockClients {
	return &MockClients{
		JiraMockClient: jira_mocks.NewMockJirer(mc),
		FiltersMock:    mocks.NewMockFilterers(mc),
	}
}

// CheckCase - ...
func (c *Case) CheckCase(t *testing.T, actualRes interface{}, actualErr error) {
	if c.ExpectedResult != nil {
		require.Equal(t, c.ExpectedResult, actualRes)
	}
	require.Equal(t, c.ExpectedError, actualErr)
}
