/*
 * Created on Thu Jul 13 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package repo

import (
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/jacklv111/common-sdk/database"
)

type CasbinRepoImpl struct {
}

func (repo *CasbinRepoImpl) GetPolicyList(domain, name, resourceType, resourceId string) ([]gormadapter.CasbinRule, error) {
	var policyList []gormadapter.CasbinRule
	var gpolicyList []gormadapter.CasbinRule

	err := database.Db.Where(
		// for g
		"(ptype = ?) AND (v0 = ?) AND (v2 = ?)",
		"g",
		domain,
		name,
	).Find(&gpolicyList).Error

	if err != nil {
		return nil, err
	}

	roles := make([]string, 0)
	for _, gpolicy := range gpolicyList {
		roles = append(roles, gpolicy.V1)
	}
	roles = append(roles, "*")

	err = database.Db.Where(
		// for p
		"(ptype = ?) AND (v0 = ? or v0 = ?) AND (v1 in ?) AND (v2 = ? or v2 = ?) AND (v3 = ? or v3 = ?)",
		"p",
		domain, "*",
		roles,
		resourceType, "*",
		resourceId, "*",
	).Find(&policyList).Error

	if err != nil {
		return nil, err
	}
	return append(policyList, gpolicyList...), err
}

var CasbinRepo CasbinRepoInterface

func init() {
	CasbinRepo = &CasbinRepoImpl{}
}
