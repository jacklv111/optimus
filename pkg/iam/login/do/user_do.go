/*
 * Created on Thu Jul 13 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package do

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserDo struct {
	ID             string `gorm:"type:varchar(255);primary_key"`
	Domain         string `gorm:"type:varchar(255);unique:domain_name,priority:1"`
	Name           string `gorm:"type:varchar(255);unique:domain_name,priority:2"`
	CreatedAt      int64  `gorm:"type:bigint;autoCreateTime:milli;<-:create"`
	UpdatedAt      int64  `gorm:"type:bigint;autoUpdateTime:milli;<-:update,create"`
	HashedPassword string `gorm:"type:varchar(255)"`
	DisplayName    string `gorm:"type:varchar(255)"`
	Email          string `gorm:"type:varchar(255)"`
	Phone          string `gorm:"type:varchar(255)"`
	Gender         string `gorm:"type:varchar(255)"`
	IsAdmin        bool   `gorm:"type:boolean"`
	IsForbidden    bool   `gorm:"type:boolean"`
}

func (UserDo) TableName() string {
	return "users"
}

func (user *UserDo) BeforeCreate(tx *gorm.DB) (err error) {
	if user.ID == "" {
		user.ID = uuid.New().String()
	}
	return
}
