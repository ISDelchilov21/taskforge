package models

import "github.com/golang-jwt/jwt/v5"

type Token struct {
	TokenType   string `json:"token_type"`
	AccessToken string `json:"access_token"`
}

type Claims struct {
	Role string `json:"role"`
	jwt.RegisteredClaims
}
