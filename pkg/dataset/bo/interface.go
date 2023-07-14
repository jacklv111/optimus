/*
 * Created on Fri Jul 14 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package bo

import (
	"io"

	aifsclientgo "github.com/jacklv111/aifs-client-go"
	"github.com/jacklv111/optimus/infra/action"
	valueobject "github.com/jacklv111/optimus/pkg/dataset/value-object"
)

//go:generate mockgen -source=interface.go -destination=./mock/mock_interface.go -package=mock

type DatasetBoInterface interface {
	GetId() string
	CreateDataset() error
	GetDetails() (valueobject.DatasetDetails, error)
	GetName() string
	HasAnnotationTemplate() bool
	GetAnnotationTemplateId() string
	Sync() error
	UpdateDataset(params valueobject.DatasetUpdateParams) error
	Delete() error
	UploadDatasetZipDataAsNewVersion(versionName, zipFormat, zipFileName string, zipData io.Reader) error

	// dataset version ---------------------------------------------------------
	GetVersionSize() int
	ExistsVersionByName(name string) bool
	CreateNewVersionFromExistedView(params valueobject.GenVersionFromExistedViewParams) error
	CreateNewVersionFromPool(poolName, versionName, versionDesc string, trainDataRatio, valDataRatio, testDataRatio int) error
	GetVersionDataItems(versionName, versionPartitionName string, offset, limit int, labelId, hasAnnotationFilter string) (dataItemList valueobject.DataItemList, err error)
	DeleteVersion(versionName string) error

	// dataset pool ------------------------------------------------------------
	ExistsPoolByName(name string) bool
	DeletePoolRawData(poolName string, rawDataIdList []string) error
	DeletePoolAnnotation(poolName string, annotationIdList []string) error
	GetDataPoolItems(poolName string, offset, limit int, labelId, hasAnnotationFilter string) (valueobject.DataItemList, error)
	CreatePool(valueobject.CreatePoolParams) error
	UpdatePool(poolName string, params valueobject.UpdatePoolParams) error
	DeletePool(poolName string) error
	GetDataPoolStatistics(poolName string) (res valueobject.DataPoolStatistics, err error)
	UploadRawDataZipToPool(poolName, zipFormat string, zipData io.Reader) error
	UploadRawDataZipToPoolAsync(actionRecord action.ActionDo, params valueobject.UploadDatasetZipDataParams)
	UploadDatasetVersionZipAsync(versionName string, actionRecord action.ActionDo, params valueobject.UploadDatasetZipDataParams)

	// MustGetPoolDataViews it must return valid rawDataViewId and annoViewId, if it can't return valid views, it will return error
	//  @param name pool name
	//  @return rawDataViewId
	//  @return annoViewId
	//  @return err
	MustGetPoolDataViews(name string) (rawDataViewId, annoViewId string, err error)

	// dataset annotation template ----------------------------------------------
	CreateAnnotationTemplate(req *aifsclientgo.CreateAnnotationTemplateRequest) error
	CreateAnnotationTemplateFromExisted(existedAnnoTempId string) error
	GetAnnotationTemplate() (resp *aifsclientgo.AnnotationTemplateDetails, err error)
	UpdateAnnotationTemplate(req *aifsclientgo.UpdateAnnotationTemplateRequest) error
}
