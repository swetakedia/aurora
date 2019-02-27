package blocksafetoml

import "github.com/stretchr/testify/mock"

// MockClient is a mockable blocksafetoml client.
type MockClient struct {
	mock.Mock
}

// GetBlocksafeToml is a mocking a method
func (m *MockClient) GetBlocksafeToml(domain string) (*Response, error) {
	a := m.Called(domain)
	return a.Get(0).(*Response), a.Error(1)
}

// GetBlocksafeTomlByAddress is a mocking a method
func (m *MockClient) GetBlocksafeTomlByAddress(address string) (*Response, error) {
	a := m.Called(address)
	return a.Get(0).(*Response), a.Error(1)
}
