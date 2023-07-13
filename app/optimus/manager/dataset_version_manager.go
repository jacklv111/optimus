/*
 * Created on Thu Jul 13 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package manager

import (
	"errors"

	"github.com/jacklv111/optimus/app/optimus/view-object/openapi"
	"github.com/jacklv111/optimus/pkg/dataset/bo"
	iamconst "github.com/jacklv111/optimus/pkg/iam/constant"
	loginvb "github.com/jacklv111/optimus/pkg/iam/login/value-object"
	psvc "github.com/jacklv111/optimus/pkg/iam/permission/service"
	pvb "github.com/jacklv111/optimus/pkg/iam/permission/value-object"
)

func (mgr *DatasetMgrImpl) CreateDatasetVersion(userInfo loginvb.UserInfo, datasetId string, req openapi.CreateDatasetVersionRequest) error {
	// check create dataset version permission
	hasAuth, err := psvc.PermissionSvc.Enforce(pvb.PermissionEnforce{
		Domain:       userInfo.Domain,
		Name:         userInfo.Name,
		ResourceType: iamconst.RESOURCE_TYPE_DATASET,
		ResourceId:   datasetId,
		Action:       iamconst.UPDATE,
	})
	if err != nil {
		return err
	}
	if !hasAuth {
		return errors.New("no permission")
	}
	datasetBo := bo.BuildWithId(datasetId)
	if err := datasetBo.Sync(); err != nil {
		return err
	}

	err = datasetBo.CreateNewVersionFromPool(req.PoolName, req.Name, req.Description, int(req.DataPartition.Train), int(req.DataPartition.Validation), int(req.DataPartition.Test))
	if err != nil {
		return err
	}
	return nil
}

func (mgr *DatasetMgrImpl) GetVersionDataItems(userInfo loginvb.UserInfo, datasetId, versionName, versionPartitionName string, offset, limit int, labelId, hasAnnotationFilter string) (openapi.DataItemList, error) {
	// check delete dataset permission
	hasAuth, err := psvc.PermissionSvc.Enforce(pvb.PermissionEnforce{
		Domain:       userInfo.Domain,
		Name:         userInfo.Name,
		ResourceType: iamconst.RESOURCE_TYPE_DATASET,
		ResourceId:   datasetId,
		Action:       iamconst.GET_DATA_ITEM_LIST,
	})
	if err != nil {
		return openapi.DataItemList{}, err
	}
	if !hasAuth {
		return openapi.DataItemList{}, errors.New("no permission")
	}

	datasetBo := bo.BuildWithId(datasetId)
	dataList, err := datasetBo.GetVersionDataItems(versionName, versionPartitionName, offset, limit, labelId, hasAnnotationFilter)
	if err != nil {
		return openapi.DataItemList{}, err
	}

	return assembleDataItemList(dataList), nil
}

func (mgr *DatasetMgrImpl) DeleteVersion(userInfo loginvb.UserInfo, datasetId, versionName string) error {
	hasAuth, err := psvc.PermissionSvc.Enforce(pvb.PermissionEnforce{
		Domain:       userInfo.Domain,
		Name:         userInfo.Name,
		ResourceType: iamconst.RESOURCE_TYPE_DATASET,
		ResourceId:   datasetId,
		Action:       iamconst.UPDATE,
	})
	if err != nil {
		return err
	}
	if !hasAuth {
		return errors.New("no permission")
	}
	datasetBo := bo.BuildWithId(datasetId)
	if err := datasetBo.Sync(); err != nil {
		return err
	}

	return datasetBo.DeleteVersion(versionName)
}
