package test

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/hihoak/auto-standup/internal"
	"github.com/hihoak/auto-standup/internal/clients/jirer"
	"github.com/hihoak/auto-standup/internal/filters"
	"github.com/hihoak/auto-standup/pkg/utils"
	"github.com/stretchr/testify/require"
	"testing"
)

var TestError = errors.New("test error")

type (
	Case struct {
		Name string
		FuncArguments []interface{}
		Setup func() (*internal.Implementator, *utils.Config)
		ExpectedError error
		ExpectedResult interface{}
	}

	MockClients struct {
		JiraMockClient *jirer.MockJirer
		FiltersMock *filters.MockFilterers
	}
)

func InitDefaultMockClients(mc *gomock.Controller) *MockClients {
	return &MockClients{
		JiraMockClient: jirer.NewMockJirer(mc),
		FiltersMock: filters.NewMockFilterers(mc),
	}
}

func InitTestImplementator(mockClients *MockClients) *internal.Implementator {
	return internal.NewImplementator(mockClients.JiraMockClient, mockClients.FiltersMock)
}

func (c *Case) CheckCase(t *testing.T, actualRes interface{}, actualErr error) {
	if c.ExpectedResult != nil {
		require.Equal(t, c.ExpectedResult, actualRes)
	}
	require.Equal(t, c.ExpectedError, actualErr)
}
