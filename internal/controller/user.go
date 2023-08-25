package controller

import (
	"errors"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/golang-jwt/jwt/v5"

	"github.com/mhkarimi1383/url-shortener/internal/database"
	"github.com/mhkarimi1383/url-shortener/types/configuration"
	"github.com/mhkarimi1383/url-shortener/types/database_models"
	"github.com/mhkarimi1383/url-shortener/types/request_schemas"
)

var (
	ErrInvalidUsernameOrPassword = errors.New("wrong Username or Password")
	ErrUserAlreadyExist          = errors.New("user already exist")
	ErrUserDoesNotExist          = errors.New("user does not exist")
)

func Login(info *requestschemas.Login) (*databasemodels.User, string, error) {
	user := databasemodels.User{Username: info.Username}
	has, _ := database.Engine.Get(&user)
	if !has {
		return nil, "", ErrInvalidUsernameOrPassword
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(info.Password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, "", ErrInvalidUsernameOrPassword
		}
		return nil, "", err
	}

	claims := &jwt.RegisteredClaims{
		ID:        strconv.FormatInt(user.Id, 10),
		Issuer:    "URLShortener",
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	st, err := token.SignedString([]byte(configuration.CurrentConfig.JWTSecret))
	if err != nil {
		return nil, "", err
	}

	return &user, st, nil
}

func CreateUser(info *requestschemas.Register, admin bool) error {
	user := databasemodels.User{Username: info.Username}
	has, _ := database.Engine.Get(&user)
	if has {
		return ErrUserAlreadyExist
	}
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(info.Password), 10) // 10 is also default
	if err != nil {
		return err
	}
	user.Password = string(encryptedPassword)
	user.Admin = admin
	_, err = database.Engine.Insert(&user)
	return err
}

func ChangeUserPassword(info *requestschemas.ChangeUserPassword, userId int64) error {
	user := databasemodels.User{Id: userId}
	has, _ := database.Engine.Get(&user)
	if !has {
		return ErrUserDoesNotExist
	}
	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(info.Password), 10) // 10 is also default
	if err != nil {
		return err
	}
	user.Password = string(encryptedPassword)
	_, err = database.Engine.Update(&user)
	return err
}
