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

const URLNotFoundError = "unable to find given url"

type URL struct {
	ID            int64     `xorm:"pk autoincr" json:"id"`
	UpstreamURL   string    `xorm:"not null" json:"upstream_url"`
	DownStreamURI string    `xorm:"not null unique" json:"downstream_uri"`
	CreatedAt     time.Time `xorm:"created not null" json:"created_at"`
	UpdatedAt     time.Time `xorm:"updated null" json:"updated_at"`
	DeletedAt     time.Time `xorm:"deleted null" json:"-"`
	Creator       int64     `xorm:"no null" json:"creator"`
	Updater       *int64    `xorm:"null" json:"updater"`
	Version       int       `xorm:"version not null default 1" json:"version"`
}

func (u *URL) Remove() error {
	user := new(User)
	user.ID = u.Creator
	_, err := user.Get()
	if err != nil {
		return err
	}
	e := GetEngine()
	_, err = e.Delete(u)
	return err
}

func (u *URL) Insert() error {
	user := new(User)
	user.ID = u.Creator
	_, err := user.Get()
	if err != nil {
		return err
	}
	e := GetEngine()
	_, err = e.InsertOne(u)
	return err
}

func (u *URL) Get() (has bool, err error) {
	e := GetEngine()
	has, err = e.Get(u)
	if !has {
		return has, errors.ErrorToHTTPErrorAdapter(
			http.StatusNotFound,
			"Unable to find given URL",
			fmt.Errorf(URLNotFoundError),
		)
	}
	return has, err
}

func ListURLs(limit, offset int) ([]URL, error) {
	var urls []URL
	e := GetEngine()
	err := e.Limit(limit, offset).Find(&urls)
	return urls, err
}

func CountURLs(limit, offset int) (int64, error) {
	e := GetEngine()
	count, err := e.Count(new(URL))
	return count, err
}

func ListURLsByUserID(limit, offset int, userID int64) ([]URL, error) {
	var urls []URL
	e := GetEngine()
	err := e.Where("Creator = ?", userID).Or("Updater = ?", userID).Limit(limit, offset).Find(&urls)
	return urls, err
}
