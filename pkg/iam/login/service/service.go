/*
 * Created on Thu Jul 13 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package service

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/jacklv111/common-sdk/log"
	"github.com/jacklv111/optimus/pkg/iam/constant"
	"github.com/jacklv111/optimus/pkg/iam/login/do"
	"github.com/jacklv111/optimus/pkg/iam/login/repo"
	loginvb "github.com/jacklv111/optimus/pkg/iam/login/value-object"
	"golang.org/x/crypto/bcrypt"
)

type loginServiceInterfaceImpl struct {
}

func (svc *loginServiceInterfaceImpl) Login(domain, username, password string) (tokenString string, err error) {
	userDo, err := repo.UserRepo.GetUserByUK(domain, username)
	if err != nil {
		return "", err
	}
	err = bcrypt.CompareHashAndPassword([]byte(userDo.HashedPassword), []byte(password))

	if err != nil {
		log.Infof("Invalid password for org: %s, username: %s", domain, username)
		return "", err
	}

	// Create the JWT token
	expirationTime := time.Now().Add(time.Second * constant.TOKEN_EXPIRE_TIME_IN_SEC).Unix()

	claims := &loginvb.UserInfo{
		ID:          userDo.ID,
		Name:        userDo.Name,
		Domain:      userDo.Domain,
		IsAdmin:     userDo.IsAdmin,
		IsForbidden: userDo.IsForbidden,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret
	return token.SignedString(constant.GetSecret())
}

func (svc *loginServiceInterfaceImpl) Signup(info loginvb.SignupInfo) (loginvb.SignupResult, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(info.Password), bcrypt.DefaultCost)
	if err != nil {
		return loginvb.SignupResult{}, err
	}
	userDo := do.UserDo{
		ID:             uuid.New().String(),
		Name:           info.Name,
		HashedPassword: string(hashedPassword),
		DisplayName:    info.DisplayName,
		Email:          info.Email,
		Phone:          info.Phone,
		Domain:         info.Domain,
		Gender:         info.Gender,
		IsAdmin:        info.IsAdmin,
		IsForbidden:    info.IsForbidden,
	}
	return loginvb.SignupResult{
		ID:     userDo.ID,
		Name:   userDo.Name,
		Domain: userDo.Domain,
	}, repo.UserRepo.CreateUser(userDo)
}

func (svc *loginServiceInterfaceImpl) ParseUserInfoFromToken(tokenString string) (loginvb.UserInfo, error) {
	userInfo := loginvb.UserInfo{}
	token, err := jwt.ParseWithClaims(tokenString, &userInfo, func(token *jwt.Token) (interface{}, error) {
		// Make sure the token method conforms to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("invalid signing method: %v", token.Method)
		}
		return constant.GetSecret(), nil // return the secret key
	})

	if err != nil {
		return userInfo, err
	}

	if claims, ok := token.Claims.(*loginvb.UserInfo); ok && token.Valid {
		return *claims, nil
	} else {
		return userInfo, err
	}
}

var LoginSvc loginServiceInterface

func init() {
	LoginSvc = &loginServiceInterfaceImpl{}
}
