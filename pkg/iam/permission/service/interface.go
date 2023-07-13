/*
 * Created on Thu Jul 13 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package service

import pvb "github.com/jacklv111/optimus/pkg/iam/permission/value-object"

//go:generate mockgen -source=interface.go -destination=./mock/mock_interface.go -package=mock

type PermissionServiceInterface interface {
	CreatePermission(pvb.Permission) error
	DeletePermission(pvb.Permission) error
	Enforce(permission pvb.PermissionEnforce) (bool, error)
	AddRoleForUserInDomain(domain, role, name string) error
	DeleteRoleForUserInDomain(domain, role, name string) error
}
