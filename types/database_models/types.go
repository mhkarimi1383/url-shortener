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
	DeletedAt time.Time `xorm:"deleted" json:"-"`
}

type Entity struct {
	Id            int64
	Name          string     `xorm:"not null unique index"`
	Description   string     `xorm:"null"`
	CreatedAt     time.Time  `xorm:"created"`
	UpdatedAt     time.Time  `xorm:"updated"`
	Version       int64      `xorm:"version"`
	Creator       User       `xorm:"bigint"`
	DeletedAt     time.Time  `xorm:"deleted" json:"-"`
	LastVisitedAt *time.Time `xorm:"null"`
	VisitCount    int64      `xorm:"default 0"`
}

type Url struct {
	Id            int64
	FullUrl       string     `xorm:"not null index"`
	ShortCode     string     `xorm:"not null unique index"`
	CreatedAt     time.Time  `xorm:"created"`
	UpdatedAt     time.Time  `xorm:"updated"`
	Version       int64      `xorm:"version"`
	Creator       User       `xorm:"bigint"`
	DeletedAt     time.Time  `xorm:"deleted" json:"-"`
	Entity        Entity     `xorm:"null bigint"`
	LastVisitedAt *time.Time `xorm:"null"`
	VisitCount    int64      `xorm:"default 0"`
}
