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

	"github.com/jacklv111/common-sdk/utils"
	"github.com/jacklv111/optimus/app/optimus/view-object/openapi"
	"github.com/jacklv111/optimus/pkg/dataset/bo"
	"github.com/jacklv111/optimus/pkg/dataset/repo"
	valueobject "github.com/jacklv111/optimus/pkg/dataset/value-object"
	iamconst "github.com/jacklv111/optimus/pkg/iam/constant"
	loginvb "github.com/jacklv111/optimus/pkg/iam/login/value-object"
	psvc "github.com/jacklv111/optimus/pkg/iam/permission/service"
	pvb "github.com/jacklv111/optimus/pkg/iam/permission/value-object"
	lightresmgmtsvc "github.com/jacklv111/optimus/pkg/resource-manager/optimus/service"
)

//go:generate mockgen -source=interface.go -destination=./mock/mock_interface.go -package=mock

type DatasetMgrInterface interface {
	CreateDataset(userInfo loginvb.UserInfo, req openapi.CreateDatasetRequest, workspace string) (id string, err error)
	GetDetails(userInfo loginvb.UserInfo, datasetId string) (openapi.DatasetDetails, error)
	GetList(userInfo loginvb.UserInfo, offset, limit int, workspace, sortBy, sortOrder string) (openapi.GetDatasetList200Response, error)
	UploadDatasetZipData(userInfo loginvb.UserInfo, datasetId, zipFormat, zipFileName string, zipData io.Reader) error
	DeleteDataset(userInfo loginvb.UserInfo, datasetId string) error
	UpdateDataset(userInfo loginvb.UserInfo, datasetId string, req openapi.UpdateDatasetRequest) error

	// dataset version
	DeleteVersion(userInfo loginvb.UserInfo, datasetId, versionName string) error
	GetVersionDataItems(userInfo loginvb.UserInfo, datasetId, versionName, versionPartitionName string, offset, limit int, labelId, hasAnnotationFilter string) (openapi.DataItemList, error)
	CreateDatasetVersion(userInfo loginvb.UserInfo, datasetId string, req openapi.CreateDatasetVersionRequest) error

	// dataset pool
	GetDataPoolStatistics(userInfo loginvb.UserInfo, datasetId, poolName string) (resp openapi.DataPoolStatistics, err error)
	DeletePool(userInfo loginvb.UserInfo, datasetId, poolName string) error
	CreateDatasetPool(userInfo loginvb.UserInfo, datasetId string, req openapi.CreateDatasetPoolRequest) error
	GetDataPoolItems(userInfo loginvb.UserInfo, datasetId, poolName string, offset, limit int, labelId, hasAnnotationFilter string) (openapi.DataItemList, error)
	DeletePoolDataItems(userInfo loginvb.UserInfo, datasetId, poolName string, rawDataIdList, annotationIdList []string) error
	UploadZipToPool(userInfo loginvb.UserInfo, datasetId, poolName, zipFormat string, zipData io.Reader) error

	// dataset annotation template
	CreateDatasetAnnotationTemplate(userInfo loginvb.UserInfo, datasetId string, req openapi.CreateAnnotationTemplateRequest) error
	GetDatasetAnnotationTemplate(userInfo loginvb.UserInfo, datasetId string) (resp openapi.AnnotationTemplateDetails, err error)
	UpdateDatasetAnnotationTemplate(userInfo loginvb.UserInfo, datasetId string, req openapi.UpdateAnnotationTemplateRequest) error
}

type DatasetMgrImpl struct {
}

func (mgr *DatasetMgrImpl) CreateDataset(userInfo loginvb.UserInfo, req openapi.CreateDatasetRequest, workspace string) (id string, err error) {
	// get resource management id
	resMgmtId, err := lightresmgmtsvc.LightResMgmtSvc.GetFirst(userInfo.Domain, workspace)
	if err != nil {
		return "", err
	}
	// create dataset
	datasetBo := bo.BuildFromCreateDatasetRequest(req, resMgmtId)
	err = datasetBo.CreateDataset()
	if err != nil {
		return "", err
	}
	// create domain admin permission for this dataset
	psvc.PermissionSvc.CreatePermission(pvb.Permission{
		Domain:       userInfo.Domain,
		ResourceType: iamconst.RESOURCE_TYPE_DATASET,
		ResourceId:   datasetBo.GetId(),
		Action:       []string{"*"},
		Effect:       iamconst.EFFECT_ALLOW,
	})
	return datasetBo.GetId(), nil
}

func (mgr *DatasetMgrImpl) GetDetails(userInfo loginvb.UserInfo, datasetId string) (openapi.DatasetDetails, error) {
	// check get dataset details permission
	hasAuth, err := psvc.PermissionSvc.Enforce(pvb.PermissionEnforce{
		Domain:       userInfo.Domain,
		Name:         userInfo.Name,
		ResourceType: iamconst.RESOURCE_TYPE_DATASET,
		ResourceId:   datasetId,
		Action:       iamconst.GET_DETAILS,
	})
	if err != nil {
		return openapi.DatasetDetails{}, err
	}
	if !hasAuth {
		return openapi.DatasetDetails{}, errors.New("no permission")
	}
	// get dataset details
	datasetBo := bo.BuildWithId(datasetId)
	details, err := datasetBo.GetDetails()
	if err != nil {
		return openapi.DatasetDetails{}, err
	}
	// add workspace details if necessary
	versionList := make([]openapi.DatasetVersionDetails, 0)
	for _, ver := range details.Versions {
		versionList = append(versionList, openapi.DatasetVersionDetails{
			Name:                 ver.VersionName,
			CreatedAt:            ver.CreatedAt,
			UpdatedAt:            ver.UpdatedAt,
			Description:          ver.Description,
			TrainRawDataNum:      int32(ver.TrainRawDataNum),
			ValidationRawDataNum: int32(ver.ValRawDataNum),
			TestRawDataNum:       int32(ver.TestRawDataNum),
		})
	}
	poolList := make([]openapi.DataPoolDetails, 0)
	for _, pool := range details.Pools {
		poolList = append(poolList, openapi.DataPoolDetails{
			Name:        pool.PoolName,
			Description: pool.Description,
			CreatedAt:   pool.CreatedAt,
			UpdatedAt:   pool.UpdatedAt,
		})
	}
	return openapi.DatasetDetails{
		Id:                     details.DatasetId,
		Name:                   details.Name,
		Description:            details.Description,
		CreatedAt:              details.CreatedAt,
		Versions:               versionList,
		Pools:                  poolList,
		RawDataType:            details.RawDataType,
		AnnotationTemplateType: details.AnnotationTemplateType,
		AnnotationTemplateId:   details.AnnotationTemplateId,
		CoverImageUrl:          details.CoverImageUrl,
	}, nil
}

func (mgr *DatasetMgrImpl) GetList(userInfo loginvb.UserInfo, offset, limit int, workspace, sortBy, sortOrder string) (resp openapi.GetDatasetList200Response, err error) {
	resMgmtId, err := lightresmgmtsvc.LightResMgmtSvc.GetFirst(userInfo.Domain, workspace)
	if err != nil {
		return resp, err
	}
	// get dataset list
	doList, err := repo.DatasetRepo.GetDatasetList(offset, limit, resMgmtId, utils.CamelToSnake(sortBy), sortOrder)
	if err != nil {
		return resp, err
	}
	count, err := repo.DatasetRepo.GetDatasetCount(resMgmtId)
	if err != nil {
		return resp, err
	}

	// assemble response
	datasetList := make([]openapi.DatasetListItem, 0)
	for _, do := range doList {
		datasetList = append(datasetList, openapi.DatasetListItem{
			Id:          do.ID,
			Name:        do.Name,
			Description: do.Description,
			CreatedAt:   do.CreatedAt,
			UpdatedAt:   do.UpdatedAt,
		})
	}
	resp.DatasetList = datasetList
	resp.TotalCount = int32(count)
	return resp, nil
}

func (mgr *DatasetMgrImpl) DeleteDataset(userInfo loginvb.UserInfo, datasetId string) error {
	datasetBo := bo.BuildWithId(datasetId)
	err := datasetBo.Sync()
	if err != nil {
		return err
	}
	// check delete dataset permission
	hasAuth, err := psvc.PermissionSvc.Enforce(pvb.PermissionEnforce{
		Domain:       userInfo.Domain,
		Name:         userInfo.Name,
		ResourceType: iamconst.RESOURCE_TYPE_DATASET,
		ResourceId:   datasetId,
		Action:       iamconst.DELETE,
	})
	if err != nil {
		return err
	}
	if !hasAuth {
		return errors.New("no permission")
	}

	err = datasetBo.Delete()
	if err != nil {
		return err
	}
	// delete domain admin permission for this dataset
	err = psvc.PermissionSvc.DeletePermission(pvb.Permission{
		Domain:       userInfo.Domain,
		ResourceType: iamconst.RESOURCE_TYPE_DATASET,
		ResourceId:   datasetBo.GetId(),
		Action:       []string{"*"},
		Effect:       iamconst.EFFECT_ALLOW,
	})
	return err
}

func (mgr *DatasetMgrImpl) UploadDatasetZipData(userInfo loginvb.UserInfo, datasetId, zipFormat, zipFileName string, zipData io.Reader) error {
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

	versionName := utils.GetFileNameWithoutSuffix(zipFileName)
	err = datasetBo.UploadDatasetZipDataAsNewVersion(versionName, zipFormat, zipFileName, zipData)
	return err
}

func (mgr *DatasetMgrImpl) UpdateDataset(userInfo loginvb.UserInfo, datasetId string, req openapi.UpdateDatasetRequest) error {
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
	return datasetBo.UpdateDataset(valueobject.DatasetUpdateParams{
		Name:                 req.Name,
		Description:          req.Description,
		CoverImageUrl:        req.CoverImageUrl,
		AnnotationTemplateId: req.AnnotationTemplateId,
	})
}

var DatasetMgr DatasetMgrInterface

func init() {
	DatasetMgr = &DatasetMgrImpl{}
}
