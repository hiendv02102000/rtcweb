package middleware

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	uuid "github.com/satori/go.uuid"
)

type JWTParam struct {
	UUID       uuid.UUID
	Authorized bool
	ExpriedAt  time.Time
}

func GenerateJWTToken(params JWTParam) (string, error) {
	var goClaims = jwt.MapClaims{}
	goClaims["authorized"] = params.Authorized
	goClaims["uuid"] = params.UUID
	goClaims["exp"] = params.ExpriedAt.Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, goClaims)

	tokenString, err := token.SignedString([]byte("access secret"))
	return tokenString, err
}

func GenerateTransactionCode() (string, error) {
	var goClaims = jwt.MapClaims{}
	goClaims["authorized"] = true
	goClaims["uuid"] = uuid.Must(uuid.NewV4(), nil) // generate uuid
	goClaims["time"] = time.Now().Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, goClaims)
	tokenString, err := token.SignedString([]byte("access secret"))
	return tokenString, err
}
