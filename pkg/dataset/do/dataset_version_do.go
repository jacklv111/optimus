/*
 * Created on Fri Jul 14 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package do

import "database/sql"

type DatasetVersionDo struct {
	DatasetId             string         `gorm:"type:varchar(255);uniqueIndex:dataset_id_name,priority:1"`
	Name                  string         `gorm:"type:varchar(255);primaryKey:dataset_id_name,priority:2"`
	Description           string         `gorm:"type:varchar(1024)"`
	CreatedAt             int64          `gorm:"type:bigint;autoCreateTime:milli;<-:create"`
	UpdatedAt             int64          `gorm:"type:bigint;autoUpdateTime:milli;<-:update,create"`
	TrainRawDataViewId    string         `gorm:"type:varchar(255)"`
	TrainAnnotationViewId sql.NullString `gorm:"type:varchar(255)"`
	TrainRawDataNum       int            `gorm:"type:int"`
	TrainTotalDataSize    int64          `gorm:"type:bigint"`
	TrainRawDataRatio     float32        `gorm:"type:float"`

	ValRawDataViewId    string         `gorm:"type:varchar(255)"`
	ValAnnotationViewId sql.NullString `gorm:"type:varchar(255)"`
	ValRawDataNum       int            `gorm:"type:int"`
	ValTotalDataSize    int64          `gorm:"type:bigint"`
	ValRawDataRatio     float32        `gorm:"type:float"`

	TestRawDataViewId    string         `gorm:"type:varchar(255)"`
	TestAnnotationViewId sql.NullString `gorm:"type:varchar(255)"`
	TestRawDataNum       int            `gorm:"type:int"`
	TestTotalDataSize    int64          `gorm:"type:bigint"`
	TestRawDataRatio     float32        `gorm:"type:float"`
}

func (DatasetVersionDo) TableName() string {
	return "dataset_versions"
}
