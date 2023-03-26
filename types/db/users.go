/*
Copyright © 2023 Muhammed Hussein Karimi <info@karimi.dev>
*/
package db

import (
	"fmt"
	"net/http"
	"time"

	"github.com/mhkarimi1383/url-shortener/utils/errors"
)

const UserNotFoundError = "unable to find given user"

type User struct {
	ID        int64     `xorm:"pk autoincr" json:"id"`
	Name      string    `xorm:"varchar(60) not null unique" json:"name"`
	Password  string    `xorm:"varchar(60) not null" json:"-"`
	CreatedAt time.Time `xorm:"created not null" json:"created_at"`
	UpdatedAt time.Time `xorm:"updated null" json:"updated_at"`
	DeletedAt time.Time `xorm:"deleted null" json:"-"`
	Version   int       `xorm:"version not null default 1" json:"version"`
}

func (u *User) Insert() error {
	e := GetEngine()
	_, err := e.InsertOne(u)
	return err
}

func (u *User) Get() (has bool, err error) {
	e := GetEngine()
	has, err = e.Get(u)
	if !has {
		return has, errors.ErrorToHTTPErrorAdapter(
			http.StatusUnauthorized,
			"Invalid Credentials",
			fmt.Errorf(UserNotFoundError),
		)
	}
	return has, err
}

func CountUsers() (int64, error) {
	e := GetEngine()
	return e.Count(new(User))
}
