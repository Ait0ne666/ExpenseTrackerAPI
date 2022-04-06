package jwt_auth

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type CustomClaims struct {
	UserID string `json:"id"`
	jwt.StandardClaims
}

type JwtAuth struct {
	key string
}

func NewJwtAuth() *JwtAuth {
	return &JwtAuth{key: "random_key"}
}

type CustomError struct {
	msg  string
	Code int
}

func NewError(msg string, code int) *CustomError {
	return &CustomError{
		msg:  msg,
		Code: code,
	}
}

func (c *CustomError) Error() string {
	return c.msg
}

const (
	authorizationHeader = "Authorization"
)

func (m *JwtAuth) GenerateJWT(userID string) (string, error) {

	claims := CustomClaims{
		userID,
		jwt.StandardClaims{
			ExpiresAt: time.Now().AddDate(1, 0, 0).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	tokenJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := tokenJWT.SignedString([]byte(m.key))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (m *JwtAuth) GetClaimsByRequest(c *gin.Context) (*CustomClaims, error) {

	header := c.GetHeader(authorizationHeader)
	if header == "" {
		return nil, NewError("you don't have enough rights", http.StatusForbidden)
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return nil, NewError("you don't have enough rights", http.StatusForbidden)
	}

	if len(headerParts[1]) == 0 {
		return nil, NewError("you don't have enough rights", http.StatusForbidden)
	}

	token, err := m.ValidateJwt(headerParts[1])
	if err != nil {
		return nil, NewError("you don't have enough rights", http.StatusForbidden)
	}

	if claims, ok := token.Claims.(*CustomClaims); ok {
		err = claims.Valid()
		if err != nil {
			return nil, err
		}
		return claims, nil
	}

	return nil, NewError("you don't have enough rights", http.StatusForbidden)
}

func (m *JwtAuth) ValidateJwt(tokenString string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unknown signature method: %v", token.Header["alg"])
		}
		return []byte(m.key), nil
	})
	if err != nil {
		return token, err
	}

	if !token.Valid {
		return token, fmt.Errorf("token is not valid, %v", token)
	}
	return token, nil
}

func (c CustomClaims) Valid() error {

	if c.UserID == "" {
		return NewError("you don't have enough rights", http.StatusForbidden)
	}

	expired := !c.VerifyExpiresAt(time.Now().Unix(), true)
	if expired {
		return NewError("you don't have enough rights", http.StatusForbidden)
	}

	return nil
}
