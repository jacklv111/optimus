/*
 * Created on Thu Jul 13 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package repo

import "github.com/jacklv111/optimus/pkg/iam/login/do"

//go:generate mockgen -source=interface.go -destination=./mock/mock_interface.go -package=mock

type UserRepoInterface interface {
	// Create user
	CreateUser(user do.UserDo) error
	// Update user
	UpdateUser() error
	// Delete user
	DeleteUser() error
	// Get user
	GetUserByUK(org, name string) (do.UserDo, error)
}
