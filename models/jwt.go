package models

import "github.com/dgrijalva/jwt-go"

//MyClaims jwt clain
type MyClaims struct {
	jwt.StandardClaims
	Nip        string `json:"nip"`
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Position   string `json:"position"`
	IdentityNO string `json:"identity_no"`
}

//TokenResp jwt clain
type TokenResp struct {
	Token    string    `json:"token"`
	Expired  string    `json:"expired"`
	DataUser UserLogin `json:"data_user"`
}

//UserLogin detail user
type UserLogin struct {
	Nip        string `json:"nip"`
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Position   string `json:"position"`
	IdentityNO string `json:"identity_no"`
}
