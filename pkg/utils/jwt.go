package utils

import "github.com/form3tech-oss/jwt-go"

type JWTClaims struct {
	jwt.StandardClaims
	Email string
	Uid   string
}
