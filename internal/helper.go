package internal

import "github.com/hihoak/auto-standup/test"

// InitTestImplementator - ...
func InitTestImplementator(mockClients *test.MockClients) *Implementator {
	return NewImplementator(mockClients.JiraMockClient, mockClients.FiltersMock)
}
