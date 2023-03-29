/*
Copyright © 2023 Muhammed Hussein Karimi <info@karimi.dev>
*/
package request

import (
	"github.com/golang-jwt/jwt/v4"
)

type User struct {
	Name     string `json:"name" xml:"name" param:"name" query:"name" form:"name"`
	Password string `json:"password" xml:"password" param:"password" query:"password" form:"password"`
}

type JWTClaims struct {
	Name   string `json:"name"`
	UserID int64  `json:"user_id"`
	jwt.RegisteredClaims
}

type TokenResponse struct {
	Token      string `json:"token" xml:"token" param:"token" query:"token" form:"token"`
	ExpireTime int64  `json:"expire_time" xml:"expire_time" param:"expire_time" query:"expire_time" form:"expire_time"`
	TokenType  string `json:"token_type" xml:"token_type" param:"token_type" query:"token_type" form:"token_type"`
}

type URL struct {
	UpsteamURL string `json:"upsteam_url" xml:"upsteam_url" param:"upsteam_url" query:"upsteam_url" form:"upsteam_url"`
}
