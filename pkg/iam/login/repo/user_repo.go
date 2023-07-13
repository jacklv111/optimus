/*
 * Created on Thu Jul 13 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package repo

import (
	"github.com/jacklv111/common-sdk/database"
	"github.com/jacklv111/optimus/pkg/iam/login/do"
)

type UserRepoImpl struct {
}

func (repo *UserRepoImpl) CreateUser(user do.UserDo) error {
	return database.Db.Create(&user).Error
}

func (repo *UserRepoImpl) UpdateUser() error {
	// todo
	return nil
}

func (repo *UserRepoImpl) DeleteUser() error {
	// todo
	return nil
}

func (repo *UserRepoImpl) GetUserByUK(org, name string) (do.UserDo, error) {
	userDo := do.UserDo{}
	err := database.Db.Where("domain = ? AND name = ?", org, name).First(&userDo).Error
	return userDo, err
}

var UserRepo UserRepoInterface

func init() {
	UserRepo = &UserRepoImpl{}
}
