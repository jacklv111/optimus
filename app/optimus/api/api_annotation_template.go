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
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jacklv111/common-sdk/errors"
	"github.com/jacklv111/common-sdk/log"
	"github.com/jacklv111/optimus/app/optimus"
	"github.com/jacklv111/optimus/app/optimus/manager"
	"github.com/jacklv111/optimus/app/optimus/view-object/openapi"
)

// CreateAnnotationTemplate - Create an annotation template
func CreateAnnotationTemplate(c *gin.Context) {
	userInfo, err := getUserInfoAndHandleError(c)
	if err != nil {
		return
	}
	var req openapi.CreateAnnotationTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Errorf("Bind json failed, error: %s", err)
		c.Error(errors.NewAppErr(optimus.INVALID_PARAMS, err, err.Error()))
		return
	}
	workspace := c.DefaultQuery(WORKSPACE, DEFAULT_WORKSPACE)
	id, err := manager.AnnotationTemplateMgr.CreateAnnotationTemplate(userInfo, req, workspace)
	if err != nil {
		log.Errorf("Create annotation template failed, error: %s", err)
		c.Error(errors.NewAppErr(optimus.UNDEFINED_ERROR, err, err.Error()))
		return
	}
	c.JSON(http.StatusOK, openapi.CreateAnnoTemplateSuccessResp{AnnotationTemplateId: id})
}

// GetAnnoTemplateDetails - Get annotation template details
func GetAnnoTemplateDetails(c *gin.Context) {
	userInfo, err := getUserInfoAndHandleError(c)
	if err != nil {
		return
	}
	id := c.Param(ANNOTATION_TEMPLATE_ID)
	details, err := manager.AnnotationTemplateMgr.GetDetails(userInfo, id)
	if err != nil {
		log.Errorf("Get annotation template details failed, error: %s", err)
		c.Error(errors.NewAppErr(optimus.UNDEFINED_ERROR, err, err.Error()))
		return
	}
	c.JSON(http.StatusOK, details)
}

// UpdateAnnotationTemplate - Update an annotation template
func UpdateAnnotationTemplate(c *gin.Context) {
	userInfo, err := getUserInfoAndHandleError(c)
	if err != nil {
		return
	}
	var req openapi.UpdateAnnotationTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Errorf("Bind json failed, error: %s", err)
		c.Error(errors.NewAppErr(optimus.INVALID_PARAMS, err, err.Error()))
		return
	}
	err = manager.AnnotationTemplateMgr.Update(userInfo, req)
	if err != nil {
		log.Errorf("Update annotation template failed, error: %s", err)
		c.Error(errors.NewAppErr(optimus.UNDEFINED_ERROR, err, err.Error()))
		return
	}

	c.Status(http.StatusOK)
}
