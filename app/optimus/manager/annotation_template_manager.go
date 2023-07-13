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
	annotationtemplate "github.com/jacklv111/optimus/pkg/annotation-template"
	iamconst "github.com/jacklv111/optimus/pkg/iam/constant"
	loginvb "github.com/jacklv111/optimus/pkg/iam/login/value-object"
	psvc "github.com/jacklv111/optimus/pkg/iam/permission/service"
	pvb "github.com/jacklv111/optimus/pkg/iam/permission/value-object"
	lightresmgmtsvc "github.com/jacklv111/optimus/pkg/resource-manager/optimus/service"
)

//go:generate mockgen -source=interface.go -destination=./mock/mock_interface.go -package=mock

type AnnotationTemplateMgrInterface interface {
	CreateAnnotationTemplate(userInfo loginvb.UserInfo, req openapi.CreateAnnotationTemplateRequest, workspace string) (id string, err error)
	GetDetails(userInfo loginvb.UserInfo, annoTempId string) (resp openapi.AnnotationTemplateDetails, err error)
	Update(userInfo loginvb.UserInfo, req openapi.UpdateAnnotationTemplateRequest) (err error)
}

type AnnotationTemplateImpl struct {
}

func (mgr *AnnotationTemplateImpl) CreateAnnotationTemplate(userInfo loginvb.UserInfo, req openapi.CreateAnnotationTemplateRequest, workspace string) (id string, err error) {
	// get resource management id
	resMgmtId, err := lightresmgmtsvc.LightResMgmtSvc.GetFirst(userInfo.Domain, workspace)
	if err != nil {
		return "", err
	}

	if req.IsCreateFromExisted {
		id, err = annotationtemplate.AnnotationTemplateSvc.CreateFromExisted(resMgmtId, req.ExistedAnnotationTemplateId)
	} else {
		id, err = annotationtemplate.AnnotationTemplateSvc.Create(resMgmtId, createAnnoTempReqToBoInput(req))
	}
	if err != nil {
		return "", err
	}
	err = psvc.PermissionSvc.CreatePermission(pvb.Permission{
		Domain:       userInfo.Domain,
		ResourceType: iamconst.RESOURCE_TYPE_ANNOTATION_TEMPLATE,
		ResourceId:   id,
		Action:       []string{"*"},
		Effect:       iamconst.EFFECT_ALLOW,
	})
	return id, err
}

func (mgr *AnnotationTemplateImpl) GetDetails(userInfo loginvb.UserInfo, annoTempId string) (resp openapi.AnnotationTemplateDetails, err error) {
	hasAuth, err := psvc.PermissionSvc.Enforce(pvb.PermissionEnforce{
		Domain:       userInfo.Domain,
		Name:         userInfo.Name,
		ResourceType: iamconst.RESOURCE_TYPE_ANNOTATION_TEMPLATE,
		Action:       iamconst.GET_DETAILS,
	})
	if err != nil {
		return resp, err
	}
	if !hasAuth {
		return resp, errors.New("no permission")
	}

	annoTemp, err := annotationtemplate.AnnotationTemplateSvc.GetDetails(annoTempId)
	if err != nil {
		return resp, err
	}
	return assembleAnnotationTemplateDetails(annoTemp), nil
}

func (mgr *AnnotationTemplateImpl) Update(userInfo loginvb.UserInfo, req openapi.UpdateAnnotationTemplateRequest) (err error) {
	hasAuth, err := psvc.PermissionSvc.Enforce(pvb.PermissionEnforce{
		Domain:       userInfo.Domain,
		Name:         userInfo.Name,
		ResourceType: iamconst.RESOURCE_TYPE_ANNOTATION_TEMPLATE,
		Action:       iamconst.UPDATE,
	})
	if err != nil {
		return err
	}
	if !hasAuth {
		return errors.New("no permission")
	}

	return annotationtemplate.AnnotationTemplateSvc.Update(updateAnnoTempReqToBoInput(req))
}

var AnnotationTemplateMgr AnnotationTemplateMgrInterface

func init() {
	AnnotationTemplateMgr = &AnnotationTemplateImpl{}
}
