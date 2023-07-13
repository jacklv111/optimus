/*
 * Created on Thu Jul 13 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package service

import (
	"github.com/jacklv111/optimus/pkg/resource-manager/optimus/repo"
	valueobject "github.com/jacklv111/optimus/pkg/resource-manager/optimus/value-object"
)

type ResourceMgmtSvcImpl struct {
}

func (resMgmt *ResourceMgmtSvcImpl) GetFirst(domain, workspace string) (resMgmtId string, err error) {
	resMgmtDo, err := repo.ResourceMgmtRepo.GetFirst(domain, workspace)
	if err != nil {
		return "", err
	}
	return resMgmtDo.ID, nil
}

func (resMgmt *ResourceMgmtSvcImpl) Create(domain, workspace string) (resMgmtId string, err error) {
	resMgmtDo, err := repo.ResourceMgmtRepo.Create(domain, workspace)
	if err != nil {
		return "", err
	}
	return resMgmtDo.ID, nil
}

func (resMgmt *ResourceMgmtSvcImpl) GetById(id string) (res valueobject.ResMgmtResult, err error) {
	resMgmtDo, err := repo.ResourceMgmtRepo.GetById(id)
	if err != nil {
		return res, err
	}
	res = valueobject.ResMgmtResult{
		Id:        resMgmtDo.ID,
		Domain:    resMgmtDo.Domain,
		Workspace: resMgmtDo.Workspace,
		CreatedAt: resMgmtDo.CreatedAt,
	}
	return res, nil
}

var ResMgmtSvc ResourceMgmtSvcInterface

func init() {
	ResMgmtSvc = &ResourceMgmtSvcImpl{}
}
