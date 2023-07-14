/*
 * Created on Fri Jul 14 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package action

import (
	"database/sql"
)

type ActionDo struct {
	ID                   sql.NullString `gorm:"column:id"`
	Name                 string         `gorm:"column:name"`
	ResourceManagementId string         `gorm:"column:resource_management_id"`
	ResourceType         string         `gorm:"column:resource_type;primaryKey:priority:1;<-:create"`
	ResourceId           string         `gorm:"column:resource_id;primaryKey:priority:2;<-:create"`
	Params               string         `gorm:"column:params"` // json string, async task params
	Action               string         `gorm:"column:action"`
	// allow read and create
	CreateAt int64 `gorm:"autoCreateTime:milli;<-:create;column:create_at"`
	// allow read and update
	UpdateAt int64 `gorm:"autoUpdateTime:milli;<-:update,create;column:update_at"`
}

func (ActionDo) TableName() string {
	return "actions"
}
