/*
 * Created on Thu Jul 13 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package repo

import "github.com/jacklv111/optimus/pkg/resource-manager/optimus/do"

//go:generate mockgen -source=interface.go -destination=./mock/mock_interface.go -package=mock

type ResourceMgmtRepoInterface interface {
	GetFirst(domain, workspace string) (do.ResourceManagementDo, error)
	Create(domain, workspace string) (do.ResourceManagementDo, error)
	GetById(id string) (do.ResourceManagementDo, error)
}
