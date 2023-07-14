/*
 * Created on Fri Jul 14 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package do

import "database/sql"

type DataPoolDo struct {
	DatasetID        string         `gorm:"type:varchar(255);primaryKey:dataset_id_name,priority:1"`
	Name             string         `gorm:"type:varchar(255);primaryKey:dataset_id_name,priority:2"`
	Description      string         `gorm:"type:varchar(1024)"`
	CreatedAt        int64          `gorm:"type:bigint;autoCreateTime:milli;<-:create"`
	UpdatedAt        int64          `gorm:"type:bigint;autoUpdateTime:milli;<-:update,create"`
	RawDataViewId    string         `gorm:"type:varchar(255)"`
	AnnotationViewId sql.NullString `gorm:"type:varchar(255)"`
}

func (DataPoolDo) TableName() string {
	return "data_pools"
}
