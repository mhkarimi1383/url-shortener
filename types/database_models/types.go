package databasemodels

import "time"

type User struct {
	Id        int64
	Admin     bool      `xorm:"index"`
	Username  string    `xorm:"not null unique index"`
	Password  string    `xorm:"not null" json:"-"` // TODO: Make this to support OAuth
	Version   int64     `xorm:"version"`
	CreatedAt time.Time `xorm:"created"`
	UpdatedAt time.Time `xorm:"updated"`
}

type Url struct {
	Id        int64
	FullUrl   string    `xorm:"not null index"`
	ShortCode string    `xorm:"null unique index"`
	CreatedAt time.Time `xorm:"created"`
	UpdatedAt time.Time `xorm:"updated"`
	Version   int64     `xorm:"version"`
	Creator   User      `xorm:"bigint"`
}
