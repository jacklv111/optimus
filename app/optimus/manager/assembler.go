/*
 * Created on Thu Jul 13 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package manager

import (
	aifsclientgo "github.com/jacklv111/aifs-client-go"
	"github.com/jacklv111/optimus/app/optimus/view-object/openapi"
	valueobject "github.com/jacklv111/optimus/pkg/dataset/value-object"
)

func assembleDataItemList(dataList valueobject.DataItemList) openapi.DataItemList {
	var res openapi.DataItemList
	res.RawDataType = dataList.RawDataType
	res.AnnotationTemplateId = dataList.AnnotationTemplateId
	res.AnnotationTemplateType = dataList.AnnotationTemplateType
	for _, data := range dataList.RawData {
		item := openapi.DataItemListDataListInner{
			RawDataId:  data.RawDataId,
			RawDataUrl: data.RawDataUrl,
		}
		if anno, ok := dataList.Annotations[data.RawDataId]; ok {
			item.AnnotationUrl = anno.AnnotationUrl
			item.LabelList = anno.LabelList
			item.AnnotationId = anno.AnnotationId
		}
		res.DataList = append(res.DataList, item)
	}
	return res
}

func createAnnoTempReqToBoInput(req openapi.CreateAnnotationTemplateRequest) *aifsclientgo.CreateAnnotationTemplateRequest {
	createAnnoReq := aifsclientgo.NewCreateAnnotationTemplateRequest(req.Name, req.Type)
	createAnnoReq.SetDescription(req.Description)
	createAnnoReq.SetWordList(req.WordList)
	labels := make([]aifsclientgo.Label, 0)
	for _, label := range req.Labels {
		ld := aifsclientgo.NewLabel(label.Name, label.Color)
		ld.SetSuperCategoryName(label.SuperCategoryName)
		ld.SetKeyPointDef(label.KeyPointDef)
		ld.SetKeyPointSkeleton(label.KeyPointSkeleton)
		ld.SetCoverImageUrl(label.CoverImageUrl)
		labels = append(labels, *ld)
	}
	createAnnoReq.SetLabels(labels)
	return createAnnoReq
}

func updateAnnoTempReqToBoInput(req openapi.UpdateAnnotationTemplateRequest) *aifsclientgo.UpdateAnnotationTemplateRequest {
	updateReq := aifsclientgo.NewUpdateAnnotationTemplateRequest(req.Name, req.Type)
	updateReq.SetDescription(req.Description)
	updateReq.SetWordList(req.WordList)
	labels := make([]aifsclientgo.Label, 0)
	for _, label := range req.Labels {
		ld := aifsclientgo.NewLabel(label.Name, label.Color)
		ld.SetSuperCategoryName(label.SuperCategoryName)
		ld.SetKeyPointDef(label.KeyPointDef)
		ld.SetKeyPointSkeleton(label.KeyPointSkeleton)
		ld.SetCoverImageUrl(label.CoverImageUrl)
		ld.SetId(label.Id)
		labels = append(labels, *ld)
	}
	updateReq.SetLabels(labels)
	updateReq.SetId(req.Id)
	return updateReq
}

func assembleAnnotationTemplateDetails(input *aifsclientgo.AnnotationTemplateDetails) (resp openapi.AnnotationTemplateDetails) {
	resp = openapi.AnnotationTemplateDetails{
		Id:       input.GetId(),
		Name:     input.GetName(),
		Type:     input.GetType(),
		WordList: input.GetWordList(),
		Labels:   make([]openapi.Label, 0),
		CreateAt: input.GetCreateAt(),
		UpdateAt: input.GetUpdateAt(),
	}
	for _, label := range input.GetLabels() {
		resp.Labels = append(resp.Labels, openapi.Label{
			Id:                label.GetId(),
			Name:              label.GetName(),
			SuperCategoryName: label.GetSuperCategoryName(),
			Color:             label.GetColor(),
			KeyPointDef:       label.GetKeyPointDef(),
			KeyPointSkeleton:  label.GetKeyPointSkeleton(),
			CoverImageUrl:     label.GetCoverImageUrl(),
		})
	}
	return
}
