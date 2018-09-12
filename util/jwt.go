package util

import "github.com/dgrijalva/jwt-go"

const (
	TokenSecret = "test-token-secret"
)

type TokenUser struct {
	ID uint `json:"id"`
}

//生成token
func GenerateToken(data map[string]interface{}) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	for key, value := range data {
		token.Claims.(jwt.MapClaims)[key] = value
	}
	return token.SignedString([]byte(TokenSecret))
}

//解析token
func ParseToken(token_string string) (*jwt.Token, error) {
	token, err := jwt.Parse(token_string, func(token *jwt.Token) (interface{}, error) {
		return []byte(TokenSecret), nil
	})
	if err == nil && token.Valid {
		return token, nil
	} else {
		return nil, err
	}
}
