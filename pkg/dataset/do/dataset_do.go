/*
 * Created on Fri Jul 14 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package do

import (
	"database/sql"
)

// dataset metadata
type DatasetDo struct {
	ID          string `gorm:"primary_key;<-:create"`
	Name        string `gorm:"type:varchar(255);;uniqueIndex:association_name,priority:2"`
	Description string `gorm:"type:varchar(1024)"`
	CreatedAt   int64  `gorm:"type:bigint;autoCreateTime:milli;<-:create"`
	UpdatedAt   int64  `gorm:"type:bigint;autoUpdateTime:milli;<-:update,create"`

	RawDataType            string         `gorm:"type:varchar(255)"`
	AnnotationTemplateType string         `gorm:"type:varchar(255)"`
	AnnotationTemplateId   sql.NullString `gorm:"type:varchar(255)"`

	// 将一些 dataset 关联到一起的 id。使用 ex：用于查询某个管理空间下的所有 dataset
	AssociationId string         `gorm:"type:varchar(255);uniqueIndex:association_name,priority:1"`
	CoverImageUrl sql.NullString `gorm:"type:varchar(255)"`
}

func (DatasetDo) TableName() string {
	return "datasets"
}
