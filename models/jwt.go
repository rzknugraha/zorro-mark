package models

import "github.com/dgrijalva/jwt-go"

//MyClaims jwt clain
type MyClaims struct {
	jwt.StandardClaims
	Nip string `json:"nip"`
	ID  int    `json:"id"`
}

//TokenResp jwt clain
type TokenResp struct {
	Token   string `json:"token"`
	Expired string `json:"expired"`
}
