package models

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"strings"
	"time"
)

type UUID struct {
	ID        string     `json:"-" gorm:"primary_key"`
	CreatedAt time.Time  `json:"-" gorm:"index"`
	UpdatedAt time.Time  `json:"-" gorm:"index"`
	DeletedAt *time.Time `json:"-" sql:"index"`
}

func (u *UUID) BeforeCreate(scope *gorm.Scope) error {
	var err error
	if u.ID == "" {
		err = scope.SetColumn("ID", strings.Replace(uuid.New().String(), "-", "", -1))
	}
	return err
}
