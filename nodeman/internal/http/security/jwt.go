package security

type JWT interface {
	ValidateToken(tokenString string) error
}
