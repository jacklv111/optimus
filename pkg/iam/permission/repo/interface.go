/*
 * Created on Thu Jul 13 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package repo

import gormadapter "github.com/casbin/gorm-adapter/v3"

//go:generate mockgen -source=interface.go -destination=./mock/mock_interface.go -package=mock

type CasbinRepoInterface interface {
	GetPolicyList(domain, name, resourceType, resourceId string) ([]gormadapter.CasbinRule, error)
}
