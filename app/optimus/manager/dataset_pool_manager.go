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
	"io"

	"github.com/jacklv111/optimus/app/optimus/view-object/openapi"
	"github.com/jacklv111/optimus/pkg/dataset/bo"
	valueobject "github.com/jacklv111/optimus/pkg/dataset/value-object"
	iamconst "github.com/jacklv111/optimus/pkg/iam/constant"
	loginvb "github.com/jacklv111/optimus/pkg/iam/login/value-object"
	psvc "github.com/jacklv111/optimus/pkg/iam/permission/service"
	pvb "github.com/jacklv111/optimus/pkg/iam/permission/value-object"
)

func (mgr *DatasetMgrImpl) GetDataPoolItems(userInfo loginvb.UserInfo, datasetId, poolName string, offset, limit int, labelId, hasAnnotationFilter string) (openapi.DataItemList, error) {
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
	dataList, err := datasetBo.GetDataPoolItems(poolName, offset, limit, labelId, hasAnnotationFilter)
	if err != nil {
		return openapi.DataItemList{}, err
	}

	return assembleDataItemList(dataList), nil
}

func (mgr *DatasetMgrImpl) UploadZipToPool(userInfo loginvb.UserInfo, datasetId, poolName, zipFormat string, zipData io.Reader) error {
	// check upload dataset zip data permission
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
	err = datasetBo.UploadRawDataZipToPool(poolName, zipFormat, zipData)
	return err
}

func (mgr *DatasetMgrImpl) DeletePoolDataItems(userInfo loginvb.UserInfo, datasetId, poolName string, rawDataIdList, annotationIdList []string) error {
	// check delete dataset data items permission
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

	// delete raw data
	if len(rawDataIdList) > 0 {
		if err := datasetBo.DeletePoolRawData(poolName, rawDataIdList); err != nil {
			return err
		}
	}

	// delete annotation
	if len(annotationIdList) > 0 {
		if err := datasetBo.DeletePoolAnnotation(poolName, annotationIdList); err != nil {
			return err
		}
	}

	return nil
}

func (mgr *DatasetMgrImpl) CreateDatasetPool(userInfo loginvb.UserInfo, datasetId string, req openapi.CreateDatasetPoolRequest) error {
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

	return datasetBo.CreatePool(valueobject.CreatePoolParams{
		Name:           req.Name,
		Desc:           req.Description,
		IsFromVersion:  req.SrcVersion != "",
		SrcVersionName: req.SrcVersion,
	})
}

func (mgr *DatasetMgrImpl) DeletePool(userInfo loginvb.UserInfo, datasetId, poolName string) error {
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

	return datasetBo.DeletePool(poolName)
}

func (mgr *DatasetMgrImpl) GetDataPoolStatistics(userInfo loginvb.UserInfo, datasetId, poolName string) (resp openapi.DataPoolStatistics, err error) {
	hasAuth, err := psvc.PermissionSvc.Enforce(pvb.PermissionEnforce{
		Domain:       userInfo.Domain,
		Name:         userInfo.Name,
		ResourceType: iamconst.RESOURCE_TYPE_DATASET,
		ResourceId:   datasetId,
		Action:       iamconst.GET_DETAILS,
	})
	if err != nil {
		return
	}
	if !hasAuth {
		err = errors.New("no permission")
		return
	}

	datasetBo := bo.BuildWithId(datasetId)
	if err = datasetBo.Sync(); err != nil {
		return
	}

	res, err := datasetBo.GetDataPoolStatistics(poolName)
	if err != nil {
		return resp, err
	}
	resp = openapi.DataPoolStatistics{
		RawDataCount:  res.RawDataCount,
		LabelCount:    res.LabelCount,
		TotalDataSize: res.RawDataTotalSize,
	}
	for _, ld := range res.LabelDistribution {
		resp.LabelDistribution = append(resp.LabelDistribution, openapi.LabelDistribution{
			LabelId: ld.LabelId,
			Count:   ld.Count,
			Ratio:   ld.Ratio,
		})
	}
	return
}
