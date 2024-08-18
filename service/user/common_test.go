package user

import (
	"github.com/stretchr/testify/mock"
	"go-sample-rest-api/types"
)

type mockUserStore struct {
	mock.Mock
}

func (m *mockUserStore) UpdateUser(u types.User) error {
	args := m.Called(u)
	return args.Error(0)
}

func (m *mockUserStore) GetUserByEmail(email string) (*types.User, error) {
	args := m.Called(email)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*types.User), args.Error(1)
}

func (m *mockUserStore) CreateUser(u types.User) error {
	args := m.Called(u)
	return args.Error(0)
}

func (m *mockUserStore) GetUserByID(id int) (*types.User, error) {
	args := m.Called(id)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.User), args.Error(1)
}

type MockAuthenticator struct {
	mock.Mock
}

func (m *MockAuthenticator) CreateJWT(secret []byte, userID int) (string, error) {
	args := m.Called(secret, userID)
	return args.String(0), args.Error(1)
}

func (m *MockAuthenticator) HashPassword(password string) (string, error) {
	args := m.Called(password)
	if args.Error(1) != nil {
		return "", args.Error(1)
	}
	return args.String(0), args.Error(1)
}

func (m *MockAuthenticator) ComparePasswords(hashed string, plain []byte) bool {
	args := m.Called(hashed, plain)
	return args.Bool(0)
}
