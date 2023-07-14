/*
 * Created on Fri Jul 14 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package annotationtemplate

import (
	"context"
	"net/http"

	aifsclientgo "github.com/jacklv111/aifs-client-go"
	aifsclient "github.com/jacklv111/common-sdk/client/aifs-client"
	"github.com/jacklv111/common-sdk/database"
	"gorm.io/gorm"
)

//go:generate mockgen -source=interface.go -destination=./mock/mock_interface.go -package=mock

type AnnotationTemplateSvcInterface interface {
	Create(resMgmtId string, req *aifsclientgo.CreateAnnotationTemplateRequest) (id string, err error)
	Delete(annoTempId string) (err error)
	GetDetails(annoTempId string) (resp *aifsclientgo.AnnotationTemplateDetails, err error)
	Update(req *aifsclientgo.UpdateAnnotationTemplateRequest) (err error)
	CreateFromExisted(resMgmtId, existedAnnotationTemplateId string) (id string, err error)
}

type AnnotationTemplateSvcImpl struct {
}

func (svc *AnnotationTemplateSvcImpl) Create(resMgmtId string, req *aifsclientgo.CreateAnnotationTemplateRequest) (id string, err error) {
	resp, _, err := aifsclient.GetAifsClient().AnnotationTemplateApi.CreateAnnotationTemplate(context.Background()).CreateAnnotationTemplateRequest(*req).Execute()
	if err != nil {
		return "", err
	}

	return resp.GetAnnotationTemplateId(), database.Db.Create(&AnnotationTemplateDo{
		ID:            resp.GetAnnotationTemplateId(),
		AssociationId: resMgmtId,
	}).Error
}

func (svc *AnnotationTemplateSvcImpl) CreateFromExisted(resMgmtId, existedAnnotationTemplateId string) (id string, err error) {
	resp, _, err := aifsclient.GetAifsClient().AnnotationTemplateApi.CopyAnnotationTemplate(context.Background(), existedAnnotationTemplateId).Execute()
	if err != nil {
		return "", err
	}
	return resp.GetAnnotationTemplateId(), database.Db.Create(&AnnotationTemplateDo{
		ID:            resp.GetAnnotationTemplateId(),
		AssociationId: resMgmtId,
	}).Error
}

func (svc *AnnotationTemplateSvcImpl) Delete(annoTempId string) (err error) {

	resp, err := aifsclient.GetAifsClient().AnnotationTemplateApi.DeleteAnnotationTemplate(context.Background(), annoTempId).Execute()
	if err != nil && resp.StatusCode != http.StatusNotFound {
		return err
	}
	err = database.Db.Delete(&AnnotationTemplateDo{}, "id = ?", annoTempId).Error
	if err == gorm.ErrRecordNotFound {
		return ErrNotFound
	}
	return err
}

func (svc *AnnotationTemplateSvcImpl) GetDetails(annoTempId string) (resp *aifsclientgo.AnnotationTemplateDetails, err error) {
	resp, _, err = aifsclient.GetAifsClient().AnnotationTemplateApi.GetAnnoTemplateDetails(context.Background(), annoTempId).Execute()
	return
}

func (svc *AnnotationTemplateSvcImpl) Update(req *aifsclientgo.UpdateAnnotationTemplateRequest) (err error) {
	_, err = aifsclient.GetAifsClient().AnnotationTemplateApi.UpdateAnnotationTemplate(context.Background()).UpdateAnnotationTemplateRequest(*req).Execute()
	return
}

var AnnotationTemplateSvc AnnotationTemplateSvcInterface

func init() {
	AnnotationTemplateSvc = &AnnotationTemplateSvcImpl{}
}
