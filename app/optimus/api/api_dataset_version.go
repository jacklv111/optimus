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
	"github.com/jacklv111/common-sdk/log"
	"github.com/jacklv111/common-sdk/utils"
	"github.com/jacklv111/optimus/app/optimus/manager"
	"github.com/jacklv111/optimus/app/optimus/view-object/openapi"
	"github.com/jacklv111/optimus/pkg/dataset"
)

// CreateDatasetVersion - Create dataset version
func CreateDatasetVersion(c *gin.Context) {
	userInfo, err := getUserInfoAndHandleError(c)
	if err != nil {
		return
	}
	datasetId := c.Param(DATASET_ID)
	var req openapi.CreateDatasetVersionRequest
	err = c.BindJSON(&req)
	if err != nil {
		log.Errorf("Bind json failed, error: %s", err)
		c.JSON(http.StatusBadRequest, openapi.Error{Code: "1", Message: err.Error()})
		return
	}
	err = manager.DatasetMgr.CreateDatasetVersion(userInfo, datasetId, req)
	if err != nil {
		if err == dataset.ErrNotFound {
			c.Status(http.StatusNotFound)
			return
		}
		log.Errorf("Create dataset version failed, error: %s", err)
		c.JSON(http.StatusInternalServerError, openapi.Error{Code: "1", Message: err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

// DeleteDatasetVersion - Delete dataset version
func DeleteDatasetVersion(c *gin.Context) {
	userInfo, err := getUserInfoAndHandleError(c)
	if err != nil {
		return
	}
	datasetId := c.Param(DATASET_ID)
	versionName := c.Param(VERSION_NAME)
	err = manager.DatasetMgr.DeleteVersion(userInfo, datasetId, versionName)
	if err != nil {
		if err == dataset.ErrNotFound {
			c.Status(http.StatusNotFound)
			return
		}
		log.Errorf("Delete dataset version failed, error: %s", err)
		c.JSON(http.StatusInternalServerError, openapi.Error{Code: "1", Message: err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

// GetVersionDataItems - get dataset version data items
func GetVersionDataItems(c *gin.Context) {
	userInfo, err := getUserInfoAndHandleError(c)
	if err != nil {
		return
	}
	datasetId := uuid.MustParse(c.Param(DATASET_ID))

	offset, err := utils.ParseInt(c.Query(OFFSET_STR), 0, math.MaxInt, 0)
	if err != nil {
		log.Errorf("Error occurred when parsing offset type %s", err)
		// TODO: 选择合适错误码和 msg
		c.JSON(http.StatusBadRequest, openapi.Error{Code: "1", Message: err.Error()})
		return
	}
	limit, err := utils.ParseInt(c.Query(LIMIT_STR), LIMIT_MIN, LIMIT_MAX, 10)
	if err != nil {
		log.Errorf("Error occurred when parsing limit type %s", err)
		// TODO: 选择合适错误码和 msg
		c.JSON(http.StatusBadRequest, openapi.Error{Code: "1", Message: err.Error()})
		return
	}
	versionName := c.Param(VERSION_NAME)
	hasAnnoFilter := c.DefaultQuery(HAS_ANNOTATION_FILTER, "all")
	labelId := c.Query(LABEL_ID)
	versionPartitionName := c.Query(VERSION_PARTITION_NAME)

	datasetItemList, err := manager.DatasetMgr.GetVersionDataItems(userInfo, datasetId.String(), versionName, versionPartitionName, offset, limit, labelId, hasAnnoFilter)
	if err != nil {
		if err == dataset.ErrNotFound {
			c.Status(http.StatusNotFound)
			return
		}
		log.Errorf("Get dataset data failed, error: %s", err)
		c.JSON(http.StatusBadRequest, openapi.Error{Code: "1", Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, datasetItemList)
}
