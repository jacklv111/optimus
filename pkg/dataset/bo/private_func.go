/*
 * Created on Fri Jul 14 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package bo

import (
	"context"
	"fmt"
	"path/filepath"

	aifsclientgo "github.com/jacklv111/aifs-client-go"
	aifsclient "github.com/jacklv111/common-sdk/client/aifs-client"
	"github.com/jacklv111/optimus/infra/action"
	dsconst "github.com/jacklv111/optimus/pkg/dataset/constant"
	valueobject "github.com/jacklv111/optimus/pkg/dataset/value-object"
	iamconst "github.com/jacklv111/optimus/pkg/iam/constant"
)

func (bo *DatasetBo) hasRunningActionOnPool(poolName string) (bool, error) {
	resourceId := bo.getPoolResourceId(poolName)
	exists, err := action.ActionMgr.ExistsByResourceTypeAndResourceId(iamconst.RESOURCE_TYPE_POOL, resourceId)
	if err != nil {
		return false, err
	}
	if exists {
		return true, nil
	}
	return false, nil
}

func (bo *DatasetBo) getVersionResourceId(name string) string {
	return filepath.Join(bo.datasetDo.ID, name)
}

func (bo *DatasetBo) getPoolResourceId(name string) string {
	return filepath.Join(bo.datasetDo.ID, name)
}

func (bo *DatasetBo) getAllDataViewList() (dataViewList []string) {
	for _, pool := range bo.poolList {
		if pool.RawDataViewId != "" {
			dataViewList = append(dataViewList, pool.RawDataViewId)
		}
		if pool.AnnotationViewId.Valid {
			dataViewList = append(dataViewList, pool.AnnotationViewId.String)
		}
	}
	for _, ver := range bo.versionList {
		dataViewList = append(dataViewList, bo.getVersionDataViews(ver.Name)...)
	}
	return
}

func fillDataItemListWithAnnotations(dataItemList *valueobject.DataItemList, resp *aifsclientgo.GetAnnotationsInDataView200Response) {
	dataItemList.AnnotationTemplateId = *resp.AnnotationTemplateId
	dataItemList.AnnotationTemplateType = *resp.AnnotationTemplateType

	if len(resp.AnnotationList) <= 0 {
		return
	}
	dataItemList.Annotations = make(map[string]valueobject.AnnotationItem, 0)
	for _, annotationData := range resp.AnnotationList {
		dataItemList.Annotations[*annotationData.RawDataId] = valueobject.AnnotationItem{
			RawDataId:     *annotationData.RawDataId,
			AnnotationUrl: *annotationData.Url,
			LabelList:     annotationData.Labels,
			AnnotationId:  *annotationData.DataItemId,
		}
	}
}

func fillDataItemListWithRawData(dataItemList *valueobject.DataItemList, resp *aifsclientgo.GetRawDataInDataView200Response) {
	if len(resp.RawDataList) > 0 {
		dataItemList.RawData = make([]valueobject.RawDataItem, 0)
		for _, rawData := range resp.RawDataList {
			dataItemList.RawData = append(dataItemList.RawData, valueobject.RawDataItem{
				RawDataId:  *rawData.RawDataId,
				RawDataUrl: *rawData.Url,
			})
		}
	}
	dataItemList.RawDataType = string(*resp.RawDataType)
}

func getDataItemList(rawDataViewId, annoViewId, labelId, hasAnnoFilter string, offset, limit int) (dataItemList valueobject.DataItemList, err error) {
	// 查询带某个 label 的 raw data
	if labelId != "" {
		if annoViewId == "" {
			return dataItemList, nil
		}
		respAnno, _, err := aifsclient.GetAifsClient().DataViewApi.GetAnnotationsInDataView(context.Background(), annoViewId).Offset(int32(offset)).Limit(int32(limit)).LabelId(labelId).Execute()
		if err != nil {
			return dataItemList, err
		}
		fillDataItemListWithAnnotations(&dataItemList, respAnno)
		respRawData, _, err := aifsclient.GetAifsClient().DataViewApi.GetRawDataInDataView(context.Background(), rawDataViewId).Offset(int32(offset)).Limit(int32(limit)).RawDataIdList(dataItemList.GetRawDataIdListFromAnnotations()).Execute()
		if err != nil {
			return dataItemList, err
		}
		fillDataItemListWithRawData(&dataItemList, respRawData)
		return dataItemList, nil
	}
	var httpcall aifsclientgo.ApiGetRawDataInDataViewRequest
	switch hasAnnoFilter {
	case dsconst.HAS_ANNOTATION_FILTER_HAS_ANNOTATION:
		httpcall = aifsclient.GetAifsClient().DataViewApi.GetRawDataInDataView(context.Background(), rawDataViewId).Offset(int32(offset)).Limit(int32(limit)).IncludedAnnotationViewId(annoViewId)
	case dsconst.HAS_ANNOTATION_FILTER_HAS_NO_ANNOTATION:
		httpcall = aifsclient.GetAifsClient().DataViewApi.GetRawDataInDataView(context.Background(), rawDataViewId).Offset(int32(offset)).Limit(int32(limit)).ExcludedAnnotationViewId(annoViewId)
	case dsconst.HAS_ANNOTATION_FILTER_ALL:
		httpcall = aifsclient.GetAifsClient().DataViewApi.GetRawDataInDataView(context.Background(), rawDataViewId).Offset(int32(offset)).Limit(int32(limit))
		// do nothing
	default:
		return dataItemList, fmt.Errorf("invalid hasAnnotationFilter: %s", hasAnnoFilter)
	}

	respRawData, _, err := httpcall.Execute()
	if err != nil {
		return dataItemList, err
	}
	fillDataItemListWithRawData(&dataItemList, respRawData)
	if annoViewId != "" {
		respAnno, _, err := aifsclient.GetAifsClient().DataViewApi.GetAnnotationsInDataView(context.Background(), annoViewId).Offset(int32(offset)).Limit(int32(limit)).RawDataIdList(dataItemList.GetRawDataIdListFromRawData()).Execute()
		if err != nil {
			return dataItemList, err
		}
		fillDataItemListWithAnnotations(&dataItemList, respAnno)
	}
	return
}

func (bo *DatasetBo) getPoolDataViews(name string) (rawDataViewId, annoViewId string) {
	pool, ok := bo.poolMap[name]
	if !ok {
		return "", ""
	}
	return pool.RawDataViewId, pool.AnnotationViewId.String
}

func (bo *DatasetBo) getVersionDataViews(name string) (dataViewIdList []string) {
	dataViewIdList = append(dataViewIdList, bo.getVersionRawDataViews(name)...)
	dataViewIdList = append(dataViewIdList, bo.getVersionAnnotationViews(name)...)
	return
}

func (bo *DatasetBo) getVersionRawDataViews(name string) (dataViewIdList []string) {
	version, ok := bo.versionMap[name]
	if !ok {
		return
	}
	if version.TrainRawDataViewId != "" {
		dataViewIdList = append(dataViewIdList, version.TrainRawDataViewId)
	}
	if version.TestRawDataViewId != "" {
		dataViewIdList = append(dataViewIdList, version.TestRawDataViewId)
	}
	if version.ValRawDataViewId != "" {
		dataViewIdList = append(dataViewIdList, version.ValRawDataViewId)
	}
	return
}

func (bo *DatasetBo) getVersionAnnotationViews(name string) (dataViewIdList []string) {
	version, ok := bo.versionMap[name]
	if !ok {
		return
	}
	if version.TrainAnnotationViewId.Valid {
		dataViewIdList = append(dataViewIdList, version.TrainAnnotationViewId.String)
	}
	if version.TestAnnotationViewId.Valid {
		dataViewIdList = append(dataViewIdList, version.TestAnnotationViewId.String)
	}
	if version.ValAnnotationViewId.Valid {
		dataViewIdList = append(dataViewIdList, version.ValAnnotationViewId.String)
	}
	return
}
