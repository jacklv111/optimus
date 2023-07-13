/*
 * Created on Thu Jul 13 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package valueobject

import "github.com/dgrijalva/jwt-go"

type UserInfo struct {
	ID          string
	Domain      string
	Name        string
	IsAdmin     bool
	IsForbidden bool
	jwt.StandardClaims
}
