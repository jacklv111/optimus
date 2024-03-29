/*
 * optimus
 *
 * optimus api
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package api

import (
	"math"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jacklv111/common-sdk/errors"
	"github.com/jacklv111/common-sdk/log"
	"github.com/jacklv111/common-sdk/utils"
	"github.com/jacklv111/optimus/app/optimus/manager"
	"github.com/jacklv111/optimus/app/optimus/view-object/openapi"
	"github.com/jacklv111/optimus/pkg/dataset"
)

// CreateDataset - Create dataset
func CreateDataset(c *gin.Context) {
	userInfo, err := getUserInfoAndHandleError(c)
	if err != nil {
		return
	}

	var req openapi.CreateDatasetRequest
	err = c.BindJSON(&req)
	if err != nil {
		log.Errorf("Error occurred when binding json %s", err)
		c.Error(errors.NewAppErr(INVALID_PARAMS, err, err.Error()))
		return
	}
	workspace := c.DefaultQuery(WORKSPACE, DEFAULT_WORKSPACE)

	id, err := manager.DatasetMgr.CreateDataset(userInfo, req, workspace)
	if err != nil {
		log.Errorf("Create dataset failed, error: %s", err)
		c.Error(errors.NewAppErr(UNDEFINED_ERROR, err, err.Error()))
		return
	}

	c.JSON(http.StatusCreated, openapi.CreateDataset200Response{Id: id})
}

// DeleteDataset - Delete dataset
func DeleteDataset(c *gin.Context) {
	userInfo, err := getUserInfoAndHandleError(c)
	if err != nil {
		return
	}
	datasetId := uuid.MustParse(c.Param(DATASET_ID))
	err = manager.DatasetMgr.DeleteDataset(userInfo, datasetId.String())
	if err != nil {
		if err == dataset.ErrNotFound {
			c.Error(errors.NewAppErr(NOT_FOUND, err, err.Error()))
			return
		}
		log.Errorf("Delete dataset failed, error: %s", err)
		c.Error(errors.NewAppErr(UNDEFINED_ERROR, err, err.Error()))
		return
	}
	c.Status(http.StatusOK)
}

// GetDatasetDetails - Get dataset details
func GetDatasetDetails(c *gin.Context) {
	userInfo, err := getUserInfoAndHandleError(c)
	if err != nil {
		return
	}

	datasetId := uuid.MustParse(c.Param(DATASET_ID))
	details, err := manager.DatasetMgr.GetDetails(userInfo, datasetId.String())
	if err != nil {
		log.Errorf("Get dataset details failed, error: %s", err)
		c.Error(errors.NewAppErr(UNDEFINED_ERROR, err, err.Error()))
		return
	}

	c.JSON(http.StatusOK, details)
}

// GetDatasetList - Get dataset list
func GetDatasetList(c *gin.Context) {
	userInfo, err := getUserInfoAndHandleError(c)
	if err != nil {
		return
	}
	workspace := c.DefaultQuery(WORKSPACE, DEFAULT_WORKSPACE)

	offset, err := utils.ParseInt(c.Query(OFFSET_STR), 0, math.MaxInt, 0)
	if err != nil {
		log.Errorf("Error occurred when parsing offset type %s", err)
		c.Error(errors.NewAppErr(INVALID_PARAMS, err, err.Error()))
		return
	}
	limit, err := utils.ParseInt(c.Query(LIMIT_STR), LIMIT_MIN, LIMIT_MAX, 10)
	if err != nil {
		log.Errorf("Error occurred when parsing limit type %s", err)
		c.Error(errors.NewAppErr(INVALID_PARAMS, err, err.Error()))
		return
	}
	sortBy := c.DefaultQuery(SORT_BY, DEFAULT_SORT_BY)
	sortOrder := c.DefaultQuery(SORT_ORDER, DEFAULT_SORT_ORDER)

	datasetList, err := manager.DatasetMgr.GetList(userInfo, offset, limit, workspace, sortBy, sortOrder)

	if err != nil {
		log.Errorf("Get dataset list failed, error: %s", err)
		c.Error(errors.NewAppErr(UNDEFINED_ERROR, err, err.Error()))
		return
	}

	c.JSON(http.StatusOK, datasetList)
}

// UploadDatasetZipData - Upload dataset zip data
func UploadDatasetZipData(c *gin.Context) {
	userInfo, err := getUserInfoAndHandleError(c)
	if err != nil {
		return
	}
	zipFormat := c.GetHeader(ZIP_FORMAT)
	zipFileName := c.GetHeader(X_ZIP_FILE_NAME)
	if zipFileName == "" {
		log.Error("zip file name is empty")
		c.Error(errors.NewAppErr(INVALID_PARAMS, err, err.Error()))
		return
	}
	defer c.Request.Body.Close()
	datasetId := c.Param(DATASET_ID)
	err = manager.DatasetMgr.UploadDatasetZipData(userInfo, datasetId, zipFormat, zipFileName, c.Request.Body)
	if err != nil {
		if err == dataset.ErrNotFound {
			c.Error(errors.NewAppErr(NOT_FOUND, err, err.Error()))
			return
		}
		log.Errorf("Upload dataset raw data failed, error: %s", err)
		c.Error(errors.NewAppErr(UNDEFINED_ERROR, err, err.Error()))
		return
	}
	c.Status(http.StatusOK)
}

// GetDataPoolDataItems - Get dataset pool data items
func GetDataPoolDataItems(c *gin.Context) {
	userInfo, err := getUserInfoAndHandleError(c)
	if err != nil {
		return
	}
	datasetId := uuid.MustParse(c.Param(DATASET_ID))

	offset, err := utils.ParseInt(c.Query(OFFSET_STR), 0, math.MaxInt, 0)
	if err != nil {
		log.Errorf("Error occurred when parsing offset type %s", err)
		c.Error(errors.NewAppErr(INVALID_PARAMS, err, err.Error()))
		return
	}
	limit, err := utils.ParseInt(c.Query(LIMIT_STR), LIMIT_MIN, LIMIT_MAX, 10)
	if err != nil {
		log.Errorf("Error occurred when parsing limit type %s", err)
		c.Error(errors.NewAppErr(INVALID_PARAMS, err, err.Error()))
		return
	}
	poolName := c.Param(POOL_NAME)
	hasAnnoFilter := c.DefaultQuery(HAS_ANNOTATION_FILTER, "all")
	labelId := c.Query(LABEL_ID)

	datasetItemList, err := manager.DatasetMgr.GetDataPoolItems(userInfo, datasetId.String(), poolName, offset, limit, labelId, hasAnnoFilter)
	if err != nil {
		if err == dataset.ErrNotFound {
			c.Error(errors.NewAppErr(NOT_FOUND, err, err.Error()))
			return
		}
		log.Errorf("Get dataset data failed, error: %s", err)
		c.Error(errors.NewAppErr(UNDEFINED_ERROR, err, err.Error()))
		return
	}
	c.JSON(http.StatusOK, datasetItemList)
}

// UpdateDataset - Update dataset
func UpdateDataset(c *gin.Context) {
	userInfo, err := getUserInfoAndHandleError(c)
	if err != nil {
		return
	}
	datasetId := c.Param(DATASET_ID)
	var req openapi.UpdateDatasetRequest
	err = c.BindJSON(&req)
	if err != nil {
		log.Errorf("Bind json failed, error: %s", err)
		c.Error(errors.NewAppErr(INVALID_PARAMS, err, err.Error()))
		return
	}
	err = manager.DatasetMgr.UpdateDataset(userInfo, datasetId, req)
	if err != nil {
		if err == dataset.ErrNotFound {
			c.Error(errors.NewAppErr(NOT_FOUND, err, err.Error()))
			return
		}
		log.Errorf("Update dataset failed, error: %s", err)
		c.Error(errors.NewAppErr(UNDEFINED_ERROR, err, err.Error()))
		return
	}

	c.Status(http.StatusOK)
}
