/*
 * Created on Fri Jul 14 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package repo

import (
	"database/sql"
	"sort"

	"github.com/jacklv111/common-sdk/database"
	"github.com/jacklv111/optimus/pkg/dataset/do"
	"gorm.io/gorm"
)

type DatasetRepoImpl struct {
}

func (repo *DatasetRepoImpl) CreateDataset(datasetDo do.DatasetDo, versions []do.DatasetVersionDo, pools []do.DataPoolDo) error {
	return database.Db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&datasetDo).Error; err != nil {
			return err
		}

		// use default create batch
		if len(versions) > 0 {
			if err := tx.Create(versions).Error; err != nil {
				return err
			}
		}
		if (len(pools)) > 0 {
			if err := tx.Create(pools).Error; err != nil {
				return err
			}
		}
		return nil
	}, &sql.TxOptions{Isolation: sql.LevelReadCommitted, ReadOnly: false})
}

func (repo *DatasetRepoImpl) CreateDatasetVersion(datasetVersionDo do.DatasetVersionDo) error {
	return database.Db.Create(&datasetVersionDo).Error
}

func (repo *DatasetRepoImpl) GetDatasetCount(associationId string) (count int64, err error) {
	err = database.Db.Model(&do.DatasetDo{}).Where("association_id = ?", associationId).Count(&count).Error
	return
}

func (repo *DatasetRepoImpl) CreateDatasetPool(datasetPoolDo do.DataPoolDo) error {
	return database.Db.Create(&datasetPoolDo).Error
}

func (repo *DatasetRepoImpl) GetById(id string) (datasetDo do.DatasetDo, versions []do.DatasetVersionDo, pools []do.DataPoolDo, err error) {
	err = database.Db.Where("id = ?", id).First(&datasetDo).Error
	if err != nil {
		return datasetDo, versions, pools, err
	}
	err = database.Db.Where("dataset_id = ?", id).Find(&versions).Error
	if err != nil {
		return datasetDo, versions, pools, err
	}
	err = database.Db.Where("dataset_id = ?", id).Find(&pools).Error
	if err != nil {
		return datasetDo, versions, pools, err
	}
	sort.Slice(versions, func(i, j int) bool {
		return versions[i].CreatedAt < versions[j].CreatedAt
	})
	sort.Slice(pools, func(i, j int) bool {
		return pools[i].CreatedAt < pools[j].CreatedAt
	})
	return
}

func (repo *DatasetRepoImpl) GetDatasetList(offset, limit int, associationId, sortBy, sortOrder string) (datasetDoList []do.DatasetDo, err error) {
	err = database.Db.Where("association_id = ?", associationId).Offset(offset).Limit(limit).Order(sortBy + " " + sortOrder).Find(&datasetDoList).Error
	return datasetDoList, err
}

func (repo *DatasetRepoImpl) UpdateDataset(curData do.DatasetDo, updateFields map[string]interface{}) error {
	return database.Db.Model(&curData).Updates(updateFields).Error
}

func (repo *DatasetRepoImpl) UpdatePool(curData do.DataPoolDo, updateFields map[string]interface{}) error {
	return database.Db.Model(&curData).Updates(updateFields).Error
}

func (repo *DatasetRepoImpl) DeleteDateset(datasetId string) error {
	return database.Db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ?", datasetId).Delete(&do.DatasetDo{}).Error; err != nil {
			return err
		}

		if err := tx.Where("dataset_id = ?", datasetId).Delete(&do.DatasetVersionDo{}).Error; err != nil {
			return err
		}
		if err := tx.Where("dataset_id = ?", datasetId).Delete(&do.DataPoolDo{}).Error; err != nil {
			return err
		}
		return nil
	}, &sql.TxOptions{Isolation: sql.LevelReadCommitted, ReadOnly: false})
}

func (repo *DatasetRepoImpl) DeletePool(datasetId, poolName string) error {
	return database.Db.Where("dataset_id = ? and name = ?", datasetId, poolName).Delete(&do.DataPoolDo{}).Error
}

func (repo *DatasetRepoImpl) DeleteVersion(datasetId, versionName string) error {
	return database.Db.Where("dataset_id = ? and name = ?", datasetId, versionName).Delete(&do.DatasetVersionDo{}).Error
}

var DatasetRepo DatasetRepoInterface

func init() {
	DatasetRepo = &DatasetRepoImpl{}
}
