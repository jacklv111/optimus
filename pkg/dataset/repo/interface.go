/*
 * Created on Fri Jul 14 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package repo

import "github.com/jacklv111/optimus/pkg/dataset/do"

//go:generate mockgen -source=interface.go -destination=./mock/mock_interface.go -package=mock

type DatasetRepoInterface interface {
	CreateDataset(datasetDo do.DatasetDo, versions []do.DatasetVersionDo, pools []do.DataPoolDo) error
	CreateDatasetVersion(datasetVersionDo do.DatasetVersionDo) error
	CreateDatasetPool(datasetPoolDo do.DataPoolDo) error
	GetById(id string) (datasetDo do.DatasetDo, versions []do.DatasetVersionDo, pools []do.DataPoolDo, err error)

	// GetDatasetList
	//  @param offset
	//  @param limit
	//  @param associationId
	//  @param sortBy snake style. example: "created_at"
	//  @param sortOrder
	//  @return datasetDoList
	//  @return err
	GetDatasetList(offset, limit int, associationId, sortBy, sortOrder string) (datasetDoList []do.DatasetDo, err error)
	GetDatasetCount(associationId string) (count int64, err error)
	UpdateDataset(curData do.DatasetDo, updateFields map[string]interface{}) error
	UpdatePool(curData do.DataPoolDo, updateFields map[string]interface{}) error
	DeleteDateset(id string) error
	DeletePool(datasetId, poolName string) error
	DeleteVersion(datasetId, versionName string) error
}
