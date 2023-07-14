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
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	aifsclientgo "github.com/jacklv111/aifs-client-go"
	aifsclient "github.com/jacklv111/common-sdk/client/aifs-client"
	"github.com/jacklv111/common-sdk/log"
	"github.com/jacklv111/optimus/infra/action"
	"github.com/jacklv111/optimus/pkg/dataset"
	dsconst "github.com/jacklv111/optimus/pkg/dataset/constant"
	"github.com/jacklv111/optimus/pkg/dataset/do"
	"github.com/jacklv111/optimus/pkg/dataset/repo"
	valueobject "github.com/jacklv111/optimus/pkg/dataset/value-object"
	iamconst "github.com/jacklv111/optimus/pkg/iam/constant"
	"gorm.io/gorm"
)

type DatasetBo struct {
	datasetDo   do.DatasetDo
	versionMap  map[string]*do.DatasetVersionDo
	poolMap     map[string]*do.DataPoolDo
	versionList []do.DatasetVersionDo
	poolList    []do.DataPoolDo
	isSynced    bool
}

func (bo *DatasetBo) GetId() string {
	return bo.datasetDo.ID
}

func (bo *DatasetBo) GetName() string {
	return bo.datasetDo.Name
}

func (bo *DatasetBo) CreateDataset() error {
	_, err := aifsclientgo.NewRawDataTypeFromValue(bo.datasetDo.RawDataType)
	if err != nil {
		return err
	}
	return repo.DatasetRepo.CreateDataset(bo.datasetDo, nil, nil)
}

func (bo *DatasetBo) GetVersionSize() int {
	return len(bo.versionList)
}

func (bo *DatasetBo) HasAnnotationTemplate() bool {
	return bo.datasetDo.AnnotationTemplateId.Valid
}

func (bo *DatasetBo) GetAnnotationTemplateId() string {
	return bo.datasetDo.AnnotationTemplateId.String
}

func (bo *DatasetBo) ExistsVersionByName(name string) bool {
	_, ok := bo.versionMap[name]
	return ok
}

func (bo *DatasetBo) ExistsPoolByName(name string) bool {
	_, ok := bo.poolMap[name]
	return ok
}

func (bo *DatasetBo) GetDataPoolStatistics(poolName string) (res valueobject.DataPoolStatistics, err error) {
	if !bo.ExistsPoolByName(poolName) {
		log.Errorf("dataset %s pool %s not found", bo.datasetDo.ID, poolName)
		return res, dataset.ErrNotFound
	}

	rawDataViewId, annoViewId := bo.getPoolDataViews(poolName)
	if rawDataViewId != "" {
		resp, _, err := aifsclient.GetAifsClient().DataViewApi.GetDataViewStatistics(context.Background(), rawDataViewId).Execute()
		if err != nil {
			return res, err
		}
		res.RawDataCount = resp.GetItemCount()
		res.RawDataTotalSize = resp.GetTotalDataSize()
	}
	if annoViewId != "" {
		resp, _, err := aifsclient.GetAifsClient().DataViewApi.GetDataViewStatistics(context.Background(), rawDataViewId).Execute()
		if err != nil {
			return res, err
		}
		res.LabelCount = resp.GetLabelCount()
		for _, ld := range resp.GetLabelDistribution() {
			res.LabelDistribution = append(res.LabelDistribution, valueobject.LabelDistribution{
				LabelId: ld.GetLabelId(),
				Count:   ld.GetCount(),
				Ratio:   ld.GetRatio(),
			})
		}
	}
	return res, nil
}

func (bo *DatasetBo) MustGetPoolDataViews(name string) (rawDataViewId, annoViewId string, err error) {
	pool, ok := bo.poolMap[name]
	if !ok {
		return "", "", dataset.ErrNotFound
	}

	if !pool.AnnotationViewId.Valid {
		if !bo.datasetDo.AnnotationTemplateId.Valid {
			return "", "", fmt.Errorf("pool %s has no annotation view and dataset %s has no annotation template", name, bo.datasetDo.ID)
		}
		// create annotation view
		req := aifsclientgo.NewCreateDataViewRequest()
		req.SetDataViewName(name)
		req.SetDescription(pool.Description)
		req.SetViewType(aifsclientgo.ANNOTATION)
		req.SetAnnotationTemplateId(bo.datasetDo.AnnotationTemplateId.String)
		req.SetRawDataViewId(pool.RawDataViewId)
		respAnno, _, err := aifsclient.GetAifsClient().DataViewApi.CreateDataView(context.Background()).CreateDataViewRequest(*req).Execute()
		if err != nil {
			return "", "", err
		}
		annoViewId = respAnno.GetDataViewId()
		// update pool
		err = bo.UpdatePool(name, valueobject.UpdatePoolParams{
			AnnotationViewId: annoViewId,
		})
		if err != nil {
			return "", "", err
		}
	}
	return rawDataViewId, annoViewId, nil
}

func (bo *DatasetBo) UpdatePool(poolName string, params valueobject.UpdatePoolParams) error {
	pool, ok := bo.poolMap[poolName]
	if !ok {
		return dataset.ErrNotFound
	}
	updateFields := make(map[string]interface{})
	if params.Name != "" {
		updateFields["name"] = params.Name
	}
	if params.Description != "" {
		updateFields["description"] = params.Description
	}
	if params.AnnotationViewId != "" && !pool.AnnotationViewId.Valid {
		updateFields["annotation_view_id"] = params.AnnotationViewId
	}

	return repo.DatasetRepo.UpdatePool(*pool, updateFields)
}

func (bo *DatasetBo) Sync() (err error) {
	if bo.isSynced {
		return nil
	}
	bo.datasetDo, bo.versionList, bo.poolList, err = repo.DatasetRepo.GetById(bo.datasetDo.ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return dataset.ErrNotFound
		}
		return err
	}
	bo.versionMap = make(map[string]*do.DatasetVersionDo)
	bo.poolMap = make(map[string]*do.DataPoolDo)
	for _, ver := range bo.versionList {
		bo.versionMap[ver.Name] = &ver
	}
	for _, pool := range bo.poolList {
		bo.poolMap[pool.Name] = &pool
	}
	bo.isSynced = true
	return nil
}

func (bo *DatasetBo) GetDetails() (details valueobject.DatasetDetails, err error) {
	if err = bo.Sync(); err != nil {
		return details, err
	}
	details = valueobject.DatasetDetails{
		DatasetId:              bo.datasetDo.ID,
		Name:                   bo.datasetDo.Name,
		Description:            bo.datasetDo.Description,
		AssociationId:          bo.datasetDo.AssociationId,
		CreatedAt:              bo.datasetDo.CreatedAt,
		RawDataType:            bo.datasetDo.RawDataType,
		AnnotationTemplateType: bo.datasetDo.AnnotationTemplateType,
		AnnotationTemplateId:   bo.datasetDo.AnnotationTemplateId.String,
		CoverImageUrl:          bo.datasetDo.CoverImageUrl.String,
	}
	for _, ver := range bo.versionList {
		details.Versions = append(details.Versions, valueobject.VersionDetails{
			VersionName:     ver.Name,
			Description:     ver.Description,
			CreatedAt:       ver.CreatedAt,
			UpdatedAt:       ver.UpdatedAt,
			TrainRawDataNum: ver.TrainRawDataNum,
			ValRawDataNum:   ver.ValRawDataNum,
			TestRawDataNum:  ver.TestRawDataNum,
		})
	}
	for _, pool := range bo.poolList {
		details.Pools = append(details.Pools, valueobject.PoolDetails{
			PoolName:    pool.Name,
			Description: pool.Description,
			CreatedAt:   pool.CreatedAt,
			UpdatedAt:   pool.UpdatedAt,
		})
	}
	return details, nil
}

func (bo *DatasetBo) GetDataPoolItems(poolName string, offset, limit int, labelId, hasAnnotationFilter string) (dataItemList valueobject.DataItemList, err error) {
	if err = bo.Sync(); err != nil {
		return dataItemList, err
	}
	pool, ok := bo.poolMap[poolName]
	if !ok {
		return dataItemList, dataset.ErrNotFound
	}
	rawDataViewId := pool.RawDataViewId
	annotationView := pool.AnnotationViewId

	return getDataItemList(rawDataViewId, annotationView.String, labelId, hasAnnotationFilter, offset, limit)
}

func (bo *DatasetBo) GetVersionDataItems(versionName, versionPartitionName string, offset, limit int, labelId, hasAnnotationFilter string) (dataItemList valueobject.DataItemList, err error) {
	if err = bo.Sync(); err != nil {
		return dataItemList, err
	}
	version, ok := bo.versionMap[versionName]
	if !ok {
		return dataItemList, dataset.ErrNotFound
	}
	switch versionPartitionName {
	case dsconst.VERSION_PARTITION_NAME_TRAIN:
		return getDataItemList(version.TrainRawDataViewId, version.TrainAnnotationViewId.String, labelId, hasAnnotationFilter, offset, limit)
	case dsconst.VERSION_PARTITION_NAME_VAL:
		return getDataItemList(version.ValRawDataViewId, version.ValAnnotationViewId.String, labelId, hasAnnotationFilter, offset, limit)
	case dsconst.VERSION_PARTITION_NAME_TEST:
		return getDataItemList(version.TestRawDataViewId, version.TestAnnotationViewId.String, labelId, hasAnnotationFilter, offset, limit)
	default:
		return dataItemList, fmt.Errorf("invalid version partition name: %s", versionPartitionName)
	}
}

func (bo *DatasetBo) Delete() (err error) {
	if err = bo.Sync(); err != nil {
		return err
	}

	// delete all related data views
	for _, dataViewId := range bo.getAllDataViewList() {
		resp, err := aifsclient.GetAifsClient().DataViewApi.DeleteDataView(context.Background(), dataViewId).Execute()
		if err != nil && resp.StatusCode != http.StatusNotFound {
			return err
		}
	}

	if bo.HasAnnotationTemplate() {
		resp, err := aifsclient.GetAifsClient().AnnotationTemplateApi.DeleteAnnotationTemplate(context.Background(), bo.datasetDo.AnnotationTemplateId.String).Execute()
		if err != nil && resp.StatusCode != http.StatusNotFound {
			return err
		}
	}

	return repo.DatasetRepo.DeleteDateset(bo.datasetDo.ID)
}

func (bo *DatasetBo) CreateNewVersionFromPool(poolName, versionName, versionDesc string, trainDataRatio, valDataRatio, testDataRatio int) error {
	if trainDataRatio+valDataRatio+testDataRatio != 100 {
		return fmt.Errorf("trainDataRatio + valDataRatio + testDataRatio != 100, it is %d", trainDataRatio+valDataRatio+testDataRatio)
	}
	if trainDataRatio < 0 || valDataRatio < 0 || testDataRatio < 0 {
		return fmt.Errorf("trainDataRatio or valDataRatio or testDataRatio is less than 0")
	}

	if err := bo.Sync(); err != nil {
		return err
	}
	hasRunningActionOnPool, err := bo.hasRunningActionOnPool(poolName)
	if err != nil {
		return err
	}
	if hasRunningActionOnPool {
		return fmt.Errorf("there is running action on pool %s, can not create new version", poolName)
	}
	rawDataViewId, annoViewId := bo.getPoolDataViews(poolName)

	// divide raw data view
	req := make([]aifsclientgo.DivideRawDataDataViewRequestInner, 0)
	partion1 := aifsclientgo.NewDivideRawDataDataViewRequestInner()
	partion1.SetRatio(int32(trainDataRatio))
	partion1.SetName(bo.datasetDo.Name + "_train")
	partion1.SetDescription(bo.datasetDo.Description)
	req = append(req, *partion1)

	partion2 := aifsclientgo.NewDivideRawDataDataViewRequestInner()
	partion2.SetRatio(int32(valDataRatio))
	partion2.SetName(bo.datasetDo.Name + "_val")
	partion2.SetDescription(bo.datasetDo.Description)
	req = append(req, *partion2)

	partion3 := aifsclientgo.NewDivideRawDataDataViewRequestInner()
	partion3.SetRatio(int32(testDataRatio))
	partion3.SetName(bo.datasetDo.Name + "_test")
	partion3.SetDescription(bo.datasetDo.Description)
	req = append(req, *partion3)

	resp, _, err := aifsclient.GetAifsClient().DataViewApi.DivideDataView(context.Background(), rawDataViewId).DivideRawDataDataViewRequestInner(req).Execute()
	if err != nil {
		return err
	}
	resMap := make(map[string]aifsclientgo.DivideRawDataDataViewResponseInner)
	for _, part := range resp {
		resMap[*part.Name] = part
	}
	newVersion := do.DatasetVersionDo{
		DatasetId:          bo.datasetDo.ID,
		Name:               versionName,
		Description:        versionDesc,
		TrainRawDataViewId: *resMap[*partion1.Name].DataViewId,
		ValRawDataViewId:   *resMap[*partion2.Name].DataViewId,
		TestRawDataViewId:  *resMap[*partion3.Name].DataViewId,
		TrainRawDataNum:    int(*resMap[*partion1.Name].ItemCount),
		ValRawDataNum:      int(*resMap[*partion2.Name].ItemCount),
		TestRawDataNum:     int(*resMap[*partion3.Name].ItemCount),
	}

	if annoViewId != "" {
		// filter annotation view
		filter := aifsclientgo.NewFilterAnnotationsInDataViewRequest()
		filter.SetRawDataViewId(*resMap[*partion1.Name].DataViewId)

		resp2, _, err := aifsclient.GetAifsClient().DataViewApi.FilterAnnotationsInDataView(context.Background(), annoViewId).FilterAnnotationsInDataViewRequest(*filter).Execute()
		if err != nil {
			return err
		}
		newVersion.TrainAnnotationViewId = sql.NullString{String: resp2.GetAnnotationViewId(), Valid: true}

		filter.SetRawDataViewId(*resMap[*partion2.Name].DataViewId)
		resp2, _, err = aifsclient.GetAifsClient().DataViewApi.FilterAnnotationsInDataView(context.Background(), annoViewId).FilterAnnotationsInDataViewRequest(*filter).Execute()
		if err != nil {
			return err
		}
		newVersion.ValAnnotationViewId = sql.NullString{String: resp2.GetAnnotationViewId(), Valid: true}

		filter.SetRawDataViewId(*resMap[*partion3.Name].DataViewId)
		resp2, _, err = aifsclient.GetAifsClient().DataViewApi.FilterAnnotationsInDataView(context.Background(), annoViewId).FilterAnnotationsInDataViewRequest(*filter).Execute()
		if err != nil {
			return err
		}
		newVersion.TestAnnotationViewId = sql.NullString{String: resp2.GetAnnotationViewId(), Valid: true}
	}
	if err := repo.DatasetRepo.CreateDatasetVersion(newVersion); err != nil {
		return err
	}
	return nil
}

func (bo *DatasetBo) CreateNewVersionFromExistedView(params valueobject.GenVersionFromExistedViewParams) error {
	newVersion := do.DatasetVersionDo{
		DatasetId:             bo.datasetDo.ID,
		Name:                  params.Name,
		TrainRawDataViewId:    params.TrainRawDataViewId,
		TrainAnnotationViewId: sql.NullString{String: params.TrainAnnotationViewId, Valid: params.TrainAnnotationViewId != ""},
		TrainRawDataNum:       int(params.TrainRawDataNum),
		TrainTotalDataSize:    params.TrainTotalDataSize,
		TrainRawDataRatio:     params.TrainRawDataRatio,

		ValRawDataViewId:    params.ValRawDataViewId,
		ValAnnotationViewId: sql.NullString{String: params.ValAnnotationViewId, Valid: params.ValAnnotationViewId != ""},
		ValRawDataNum:       int(params.ValRawDataNum),
		ValTotalDataSize:    params.ValTotalDataSize,
		ValRawDataRatio:     params.ValRawDataRatio,

		TestRawDataViewId:    params.TestRawDataViewId,
		TestAnnotationViewId: sql.NullString{String: params.TestAnnotationViewId, Valid: params.TestAnnotationViewId != ""},
		TestRawDataNum:       int(params.TestRawDataNum),
		TestTotalDataSize:    params.TestTotalDataSize,
		TestRawDataRatio:     params.TestRawDataRatio,
	}
	if err := repo.DatasetRepo.CreateDatasetVersion(newVersion); err != nil {
		return err
	}
	return nil
}

func (bo *DatasetBo) UpdateDataset(params valueobject.DatasetUpdateParams) error {
	updateFields := make(map[string]interface{}, 0)
	if params.Name != "" {
		updateFields["name"] = params.Name
	}
	if params.Description != "" {
		updateFields["description"] = params.Description
	}
	if params.CoverImageUrl != "" {
		updateFields["cover_image_url"] = params.CoverImageUrl
	}
	if params.RawDataType != "" {
		if bo.datasetDo.RawDataType != "" {
			return errors.New("dataset already has raw data type")
		}
		updateFields["raw_data_type"] = params.RawDataType
	}
	if params.AnnotationTemplateType != "" {
		if bo.datasetDo.AnnotationTemplateType != "" {
			return errors.New("dataset already has annotation template")
		}
		updateFields["annotation_template_type"] = params.AnnotationTemplateType
	}
	if params.AnnotationTemplateId != "" {
		if bo.HasAnnotationTemplate() {
			return errors.New("dataset already has annotation template")
		}
		updateFields["annotation_template_id"] = params.AnnotationTemplateId
	}

	return repo.DatasetRepo.UpdateDataset(bo.datasetDo, updateFields)
}

func (bo *DatasetBo) DeletePoolRawData(poolName string, rawDataIdList []string) error {
	if err := bo.updatePoolValidate(poolName); err != nil {
		return err
	}

	rawDataViewId, _ := bo.getPoolDataViews(poolName)
	if rawDataViewId == "" {
		return nil
	}
	_, err := aifsclient.GetAifsClient().DataViewApi.DeleteDataItemInDataView(context.Background(), rawDataViewId).DataViewItemIdList(rawDataIdList).Execute()
	return err
}

func (bo *DatasetBo) DeletePoolAnnotation(poolName string, annotationIdList []string) error {
	if err := bo.updatePoolValidate(poolName); err != nil {
		return err
	}

	_, annoViewId := bo.getPoolDataViews(poolName)
	if annoViewId == "" {
		return nil
	}
	_, err := aifsclient.GetAifsClient().DataViewApi.DeleteDataItemInDataView(context.Background(), annoViewId).DataViewItemIdList(annotationIdList).Execute()
	return err
}

func (bo *DatasetBo) DeletePool(poolName string) error {
	if err := bo.updatePoolValidate(poolName); err != nil {
		return err
	}

	rawDataView, annoView := bo.getPoolDataViews(poolName)
	// delete raw data view
	resp, err := aifsclient.GetAifsClient().DataViewApi.DeleteDataView(context.Background(), rawDataView).Execute()
	if err != nil && resp.StatusCode != http.StatusNotFound {
		return err
	}
	// delete annotation view
	if annoView != "" {
		resp, err = aifsclient.GetAifsClient().DataViewApi.DeleteDataView(context.Background(), annoView).Execute()
		if err != nil && resp.StatusCode != http.StatusNotFound {
			return err
		}
	}
	return repo.DatasetRepo.DeletePool(bo.datasetDo.ID, poolName)
}

func (bo *DatasetBo) DeleteVersion(versionName string) error {
	if !bo.ExistsVersionByName(versionName) {
		log.Errorf("dataset %s version %s not found", bo.datasetDo.ID, versionName)
		return dataset.ErrNotFound
	}
	dataViewList := bo.getVersionDataViews(versionName)
	for _, dataView := range dataViewList {
		resp, err := aifsclient.GetAifsClient().DataViewApi.DeleteDataView(context.Background(), dataView).Execute()
		if err != nil && resp.StatusCode != http.StatusNotFound {
			return err
		}
	}
	return repo.DatasetRepo.DeleteVersion(bo.datasetDo.ID, versionName)
}

func (bo *DatasetBo) CreatePool(params valueobject.CreatePoolParams) error {
	if bo.ExistsPoolByName(params.Name) {
		return errors.New("pool name already exists")
	}
	newPool := do.DataPoolDo{
		DatasetID:   bo.datasetDo.ID,
		Name:        params.Name,
		Description: params.Desc,
	}
	if params.IsFromVersion {
		if !bo.ExistsVersionByName(params.SrcVersionName) {
			return dataset.ErrNotFound
		}
		// merge raw data view
		req := aifsclientgo.NewMergeDataViewsRequest()
		req.Name = &params.Name
		req.Description = &params.Desc
		req.DataViewIdList = bo.getVersionRawDataViews(params.SrcVersionName)
		respRawData, _, err := aifsclient.GetAifsClient().DataViewApi.MergeDataViews(context.Background()).MergeDataViewsRequest(*req).Execute()
		if err != nil {
			return err
		}
		newPool.RawDataViewId = respRawData.GetDataViewId()
		// merge annotation view
		req = aifsclientgo.NewMergeDataViewsRequest()
		req.Name = &params.Name
		req.Description = &params.Desc
		req.DataViewIdList = bo.getVersionAnnotationViews(params.SrcVersionName)
		respAnno, _, err := aifsclient.GetAifsClient().DataViewApi.MergeDataViews(context.Background()).MergeDataViewsRequest(*req).Execute()
		if err != nil {
			return err
		}
		newPool.AnnotationViewId = sql.NullString{String: respAnno.GetDataViewId(), Valid: true}
	} else {
		// must create raw data view
		req := aifsclientgo.NewCreateDataViewRequest()
		req.SetDataViewName(params.Name)
		req.SetDescription(params.Desc)
		req.SetViewType(aifsclientgo.RAW_DATA)
		req.SetRawDataType(aifsclientgo.RawDataType(bo.datasetDo.RawDataType))
		respRawData, _, err := aifsclient.GetAifsClient().DataViewApi.CreateDataView(context.Background()).CreateDataViewRequest(*req).Execute()
		if err != nil {
			return err
		}
		newPool.RawDataViewId = respRawData.GetDataViewId()
		// create annotation view if dataset have annotation template
		if bo.datasetDo.AnnotationTemplateId.Valid {
			// create annotation view
			req = aifsclientgo.NewCreateDataViewRequest()
			req.SetDataViewName(params.Name)
			req.SetDescription(params.Desc)
			req.SetViewType(aifsclientgo.ANNOTATION)
			req.SetAnnotationTemplateId(bo.datasetDo.AnnotationTemplateId.String)
			req.SetRelatedDataViewId(newPool.RawDataViewId)
			respAnno, _, err := aifsclient.GetAifsClient().DataViewApi.CreateDataView(context.Background()).CreateDataViewRequest(*req).Execute()
			if err != nil {
				return err
			}
			newPool.AnnotationViewId = sql.NullString{String: respAnno.GetDataViewId(), Valid: true}
		}
	}
	return repo.DatasetRepo.CreateDatasetPool(newPool)
}

func (bo *DatasetBo) UploadDatasetZipDataAsNewVersion(versionName, zipFormat, zipFileName string, zipData io.Reader) error {
	if err := bo.Sync(); err != nil {
		return err
	}
	if bo.ExistsVersionByName(versionName) {
		return errors.New("version name already exists")
	}
	resourceId := bo.getVersionResourceId(versionName)
	exists, err := action.ActionMgr.ExistsByResourceTypeAndResourceId(iamconst.RESOURCE_TYPE_VERSION, resourceId)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("dataset is running in a action, cannot update")
	}
	// it is not allowed to upload only raw data by this function
	if zipFormat == string(aifsclientgo.RAW_DATA_IMAGES) {
		return errors.New("it is not supported to upload only raw data by this function")
	}

	req := *aifsclientgo.NewCreateDataViewRequest()
	req.SetDataViewName(versionName)
	req.SetViewType(aifsclientgo.DATASET_ZIP)
	req.SetZipFormat(aifsclientgo.ZipFormat(zipFormat))
	req.SetDescription("from " + zipFileName)
	if bo.HasAnnotationTemplate() {
		req.SetAnnotationTemplateId(bo.datasetDo.AnnotationTemplateId.String)
	}
	resp, _, err := aifsclient.GetAifsClient().DataViewApi.CreateDataView(context.Background()).CreateDataViewRequest(req).Execute()
	if err != nil {
		log.Errorf("Create data view failed, error: %s", err)
		return err
	}
	datasetZipViewId := resp.GetDataViewId()

	_, err = aifsclient.GetAifsClient().DataViewUploadApi.UploadDatasetZipToDataView(context.Background(), datasetZipViewId).XFileName(zipFileName).Body(zipData).Execute()

	if err != nil {
		log.Errorf("Upload dataset zip to data view failed, error: %s", err)
		return err
	}

	params := valueobject.UploadDatasetZipDataParams{
		DatasetZipViewId: datasetZipViewId,
	}
	jsonParams, err := json.Marshal(params)
	if err != nil {
		return err
	}
	// create action record
	actionRecord, err := action.ActionMgr.Create(iamconst.RESOURCE_TYPE_VERSION, resourceId, dsconst.UPLOAD_DATASET_VERSION_ZIP, datasetZipViewId, string(jsonParams))
	if err != nil {
		_, _ = aifsclient.GetAifsClient().DataViewApi.HardDeleteDataView(context.Background(), datasetZipViewId).Execute()
		return err
	}

	// start to decompress dataset zip
	go bo.UploadDatasetVersionZipAsync(versionName, actionRecord, params)
	return nil
}

func (bo *DatasetBo) UploadRawDataZipToPool(poolName, zipFormat string, zipData io.Reader) error {
	if err := bo.Sync(); err != nil {
		return err
	}
	if !bo.ExistsPoolByName(poolName) {
		log.Errorf("dataset %s pool %s not found", bo.datasetDo.ID, poolName)
		return dataset.ErrNotFound
	}

	resourceId := bo.getPoolResourceId(poolName)
	exists, err := action.ActionMgr.ExistsByResourceTypeAndResourceId(iamconst.RESOURCE_TYPE_POOL, resourceId)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("dataset is running in a action, cannot update")
	}

	// 当前只支持raw data images
	if zipFormat != string(aifsclientgo.RAW_DATA_IMAGES) {
		return errors.New("only support raw data images")
	}

	req := *aifsclientgo.NewCreateDataViewRequest()
	req.SetDataViewName(poolName)
	req.SetViewType(aifsclientgo.DATASET_ZIP)
	req.SetZipFormat(aifsclientgo.ZipFormat(zipFormat))
	rawDataViewId, _ := bo.getPoolDataViews(poolName)
	req.SetRawDataViewId(rawDataViewId)
	req.SetDescription(fmt.Sprintf("upload raw data to %s", resourceId))

	resp, _, err := aifsclient.GetAifsClient().DataViewApi.CreateDataView(context.Background()).CreateDataViewRequest(req).Execute()
	if err != nil {
		log.Errorf("Create data view failed, error: %s", err)
		return err
	}
	datasetZipViewId := resp.GetDataViewId()

	_, err = aifsclient.GetAifsClient().DataViewUploadApi.UploadDatasetZipToDataView(context.Background(), datasetZipViewId).XFileName(resourceId).Body(zipData).Execute()

	if err != nil {
		log.Errorf("Upload dataset zip to data view failed, error: %s", err)
		return err
	}

	params := valueobject.UploadDatasetZipDataParams{
		DatasetZipViewId: datasetZipViewId,
	}

	jsonParams, err := json.Marshal(params)
	if err != nil {
		return err
	}
	// create action record
	actionRecord, err := action.ActionMgr.Create(iamconst.RESOURCE_TYPE_POOL, resourceId, dsconst.UPLOAD_RAW_DATA_ZIP, datasetZipViewId, string(jsonParams))

	if err != nil {
		_, _ = aifsclient.GetAifsClient().DataViewApi.HardDeleteDataView(context.Background(), datasetZipViewId).Execute()
		return err
	}

	// start to decompress dataset zip
	go bo.UploadRawDataZipToPoolAsync(actionRecord, params)
	return nil
}

func (bo *DatasetBo) CreateAnnotationTemplate(req *aifsclientgo.CreateAnnotationTemplateRequest) error {
	if bo.HasAnnotationTemplate() {
		return errors.New("dataset already has annotation template")
	}
	resp, _, err := aifsclient.GetAifsClient().AnnotationTemplateApi.CreateAnnotationTemplate(context.Background()).CreateAnnotationTemplateRequest(*req).Execute()
	if err != nil {
		return err
	}
	err = bo.UpdateDataset(valueobject.DatasetUpdateParams{
		AnnotationTemplateId:   resp.GetAnnotationTemplateId(),
		AnnotationTemplateType: req.GetType(),
	})
	return err
}

func (bo *DatasetBo) CreateAnnotationTemplateFromExisted(existedAnnoTempId string) error {
	if bo.HasAnnotationTemplate() {
		return errors.New("dataset already has annotation template")
	}
	resp, _, err := aifsclient.GetAifsClient().AnnotationTemplateApi.CopyAnnotationTemplate(context.Background(), existedAnnoTempId).Execute()
	if err != nil {
		return err
	}
	details, _, err := aifsclient.GetAifsClient().AnnotationTemplateApi.GetAnnoTemplateDetails(context.Background(), resp.GetAnnotationTemplateId()).Execute()
	if err != nil {
		return err
	}
	err = bo.UpdateDataset(valueobject.DatasetUpdateParams{
		AnnotationTemplateId:   details.GetId(),
		AnnotationTemplateType: details.GetType(),
	})
	return err
}

func (bo *DatasetBo) GetAnnotationTemplate() (resp *aifsclientgo.AnnotationTemplateDetails, err error) {
	if err := bo.Sync(); err != nil {
		return resp, err
	}
	if !bo.HasAnnotationTemplate() {
		return nil, dataset.ErrNotFound
	}
	resp, _, err = aifsclient.GetAifsClient().AnnotationTemplateApi.GetAnnoTemplateDetails(context.Background(), bo.datasetDo.AnnotationTemplateId.String).Execute()
	return
}
func (bo *DatasetBo) UpdateAnnotationTemplate(req *aifsclientgo.UpdateAnnotationTemplateRequest) error {
	if err := bo.Sync(); err != nil {
		return err
	}
	if !bo.HasAnnotationTemplate() {
		return dataset.ErrNotFound
	}
	_, err := aifsclient.GetAifsClient().AnnotationTemplateApi.UpdateAnnotationTemplate(context.Background()).UpdateAnnotationTemplateRequest(*req).Execute()
	if err != nil {
		return err
	}
	return nil
}
