/*
Copyright © 2023 Muhammed Hussein Karimi <info@karimi.dev>
*/
package types

type Config struct {
	ListenAddress      string
	DBEngine           string
	DBConnectionString string
	JWTSecret          string
	Debug              bool
	TokenExpireTime    int64
}
