/*
 * Created on Thu Jul 13 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package valueobject

type SignupInfo struct {
	Domain      string
	Name        string
	Password    string
	DisplayName string
	Email       string
	Phone       string
	Gender      string
	IsAdmin     bool
	IsForbidden bool
}
