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

func (mgr *DatasetMgrImpl) CreateDatasetAnnotationTemplate(userInfo loginvb.UserInfo, datasetId string, req openapi.CreateAnnotationTemplateRequest) error {
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

	if req.IsCreateFromExisted {
		err = datasetBo.CreateAnnotationTemplateFromExisted(req.ExistedAnnotationTemplateId)
	} else {
		err = datasetBo.CreateAnnotationTemplate(createAnnoTempReqToBoInput(req))
	}
	return err
}
func (mgr *DatasetMgrImpl) GetDatasetAnnotationTemplate(userInfo loginvb.UserInfo, datasetId string) (resp openapi.AnnotationTemplateDetails, err error) {
	hasAuth, err := psvc.PermissionSvc.Enforce(pvb.PermissionEnforce{
		Domain:       userInfo.Domain,
		Name:         userInfo.Name,
		ResourceType: iamconst.RESOURCE_TYPE_DATASET,
		ResourceId:   datasetId,
		Action:       iamconst.GET_DETAILS,
	})
	if err != nil {
		return resp, err
	}
	if !hasAuth {
		return resp, errors.New("no permission")
	}
	datasetBo := bo.BuildWithId(datasetId)
	details, err := datasetBo.GetAnnotationTemplate()
	if err != nil {
		return resp, err
	}
	return assembleAnnotationTemplateDetails(details), nil
}

func (mgr *DatasetMgrImpl) UpdateDatasetAnnotationTemplate(userInfo loginvb.UserInfo, datasetId string, req openapi.UpdateAnnotationTemplateRequest) error {
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
	return datasetBo.UpdateAnnotationTemplate(updateAnnoTempReqToBoInput(req))
}
