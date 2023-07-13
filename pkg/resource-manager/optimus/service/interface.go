/*
 * Created on Thu Jul 13 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package service

import valueobject "github.com/jacklv111/optimus/pkg/resource-manager/optimus/value-object"

//go:generate mockgen -source=interface.go -destination=./mock/mock_interface.go -package=mock

type ResourceMgmtSvcInterface interface {
	GetFirst(domain, workspace string) (resMgmtId string, err error)
	Create(domain, workspace string) (resMgmtId string, err error)
	GetById(id string) (res valueobject.ResMgmtResult, err error)
}
