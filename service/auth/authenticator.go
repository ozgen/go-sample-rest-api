package auth

type Authenticator interface {
	CreateJWT(secret []byte, userID int) (string, error)
	HashPassword(password string) (string, error)
	ComparePasswords(hashed string, plain []byte) bool
}

// Implementation that will be used in production
type RealAuthenticator struct{}

func NewAuthenticator() *RealAuthenticator {
	return &RealAuthenticator{}
}

func (*RealAuthenticator) CreateJWT(secret []byte, userID int) (string, error) {
	return CreateJWT(secret, userID)
}

func (*RealAuthenticator) HashPassword(password string) (string, error) {
	return HashPassword(password)
}

func (*RealAuthenticator) ComparePasswords(hashed string, plain []byte) bool {
	return ComparePasswords(hashed, plain)
}
