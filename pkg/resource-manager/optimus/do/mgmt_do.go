/*
 * Created on Thu Jul 13 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package do

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ResourceManagementDo struct {
	ID        string `gorm:"type:varchar(255);primaryKey;<-:create"`
	Domain    string `gorm:"type:varchar(255);uniqueIndex:domain_workspace,priority:1"`
	Workspace string `gorm:"type:varchar(255);uniqueIndex:domain_workspace,priority:2"`
	CreatedAt int64  `gorm:"type:bigint;autoCreateTime:milli;<-:create"`
}

func (ResourceManagementDo) TableName() string {
	return "resource_management"
}

func (resMgmt *ResourceManagementDo) BeforeCreate(tx *gorm.DB) (err error) {
	if resMgmt.ID == "" {
		resMgmt.ID = uuid.New().String()
	}
	return
}
