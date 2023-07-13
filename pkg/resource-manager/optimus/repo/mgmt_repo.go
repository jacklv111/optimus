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
	"github.com/jacklv111/optimus/pkg/resource-manager/optimus/do"
)

type ResourceMgmtRepoImpl struct {
}

func (ResourceMgmtRepoImpl) GetFirst(domain, workspace string) (do.ResourceManagementDo, error) {
	resMgmtDo := do.ResourceManagementDo{
		Domain:    domain,
		Workspace: workspace,
	}
	err := database.Db.FirstOrCreate(&resMgmtDo, resMgmtDo).Error
	return resMgmtDo, err
}

func (ResourceMgmtRepoImpl) GetById(id string) (do.ResourceManagementDo, error) {
	resMgmtDo := do.ResourceManagementDo{}
	err := database.Db.Where("id = ?", id).First(&resMgmtDo).Error
	return resMgmtDo, err
}

func (ResourceMgmtRepoImpl) Create(domain, workspace string) (do.ResourceManagementDo, error) {
	resMgmtDo := do.ResourceManagementDo{
		Domain:    domain,
		Workspace: workspace,
	}
	err := database.Db.Create(&resMgmtDo).Error
	return resMgmtDo, err
}

var ResourceMgmtRepo ResourceMgmtRepoInterface

func init() {
	ResourceMgmtRepo = ResourceMgmtRepoImpl{}
}
