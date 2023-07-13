/*
 * Created on Thu Jul 13 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package service

import loginvb "github.com/jacklv111/optimus/pkg/iam/login/value-object"

//go:generate mockgen -source=interface.go -destination=./mock/mock_interface.go -package=mock

type loginServiceInterface interface {
	Login(domain, username, password string) (tokenString string, err error)
	Signup(info loginvb.SignupInfo) (loginvb.SignupResult, error)
	ParseUserInfoFromToken(tokenString string) (loginvb.UserInfo, error)
}
