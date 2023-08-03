/*
 * iam
 *
 * iam api
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jacklv111/common-sdk/errors"
	"github.com/jacklv111/common-sdk/log"
	"github.com/jacklv111/optimus/app/iam/view-object/openapi"
	"github.com/jacklv111/optimus/pkg/iam/constant"
	"github.com/jacklv111/optimus/pkg/iam/login/service"
	psvc "github.com/jacklv111/optimus/pkg/iam/permission/service"
	pvb "github.com/jacklv111/optimus/pkg/iam/permission/value-object"
)

// CreatePermission - Create permission
func CreatePermission(c *gin.Context) {
	token := c.GetHeader(constant.AUTHORIZATION)
	userInfo, err := service.LoginSvc.ParseUserInfoFromToken(token)
	if err != nil {
		log.Errorf("Error occurred when parsing token %s", err)
		c.Error(errors.NewAppErr(INVALID_PARAMS, err, err.Error()))
		return
	}
	var req openapi.Permission
	err = c.BindJSON(&req)
	if err != nil {
		log.Errorf("Error occurred when binding json %s", err)
		c.Error(errors.NewAppErr(INVALID_PARAMS, err, err.Error()))
	}

	// 检查用户是否有对资源授权的权限
	hasAuth, err := psvc.PermissionSvc.Enforce(pvb.PermissionEnforce{
		Domain:       userInfo.Domain,
		Name:         userInfo.Name,
		ResourceType: req.ResourceType,
		ResourceId:   req.ResourceId,
		Action:       constant.AUTHORIZE,
	})
	if err != nil {
		log.Errorf("Error occurred when enforcing permission %s", err)
		c.Error(errors.NewAppErr(UNDEFINED_ERROR, err, err.Error()))
		return
	}

	if !hasAuth {
		log.Errorf("User %s has no permission to authorize resource %s", userInfo.Name, req.ResourceType)
		c.Error(errors.NewAppErr(NO_PERMISSION, err, err.Error()))
		return
	}

	err = psvc.PermissionSvc.CreatePermission(pvb.Permission{
		Domain:       userInfo.Domain,
		ResourceType: req.ResourceType,
		ResourceId:   req.ResourceId,
		Action:       req.Action,
		Effect:       req.Effect,
	})
	if err != nil {
		log.Errorf("Error occurred when creating permission %s", err)
		c.Error(errors.NewAppErr(UNDEFINED_ERROR, err, err.Error()))
		return
	}

	c.JSON(http.StatusCreated, gin.H{})
}

// DeletePermission - Delete permission
func DeletePermission(c *gin.Context) {
	token := c.GetHeader(constant.AUTHORIZATION)
	userInfo, err := service.LoginSvc.ParseUserInfoFromToken(token)
	if err != nil {
		log.Errorf("Error occurred when parsing token %s", err)
		c.Error(errors.NewAppErr(INVALID_PARAMS, err, err.Error()))
		return
	}
	var req openapi.Permission
	err = c.BindJSON(&req)
	if err != nil {
		log.Errorf("Error occurred when binding json %s", err)
		c.Error(errors.NewAppErr(INVALID_PARAMS, err, err.Error()))
	}

	// 检查用户是否有对资源授权的权限
	hasAuth, err := psvc.PermissionSvc.Enforce(pvb.PermissionEnforce{
		Domain:       userInfo.Domain,
		Name:         userInfo.Name,
		ResourceType: req.ResourceType,
		ResourceId:   req.ResourceId,
		Action:       constant.AUTHORIZE,
	})
	if err != nil {
		log.Errorf("Error occurred when enforcing permission %s", err)
		c.Error(errors.NewAppErr(UNDEFINED_ERROR, err, err.Error()))
		return
	}

	if !hasAuth {
		log.Errorf("User %s has no permission to authorize resource %s", userInfo.Name, req.ResourceType)
		c.Error(errors.NewAppErr(NO_PERMISSION, err, err.Error()))
		return
	}

	err = psvc.PermissionSvc.DeletePermission(pvb.Permission{
		Domain:       userInfo.Domain,
		ResourceType: req.ResourceType,
		ResourceId:   req.ResourceId,
		Action:       req.Action,
		Effect:       req.Effect,
	})
	if err != nil {
		log.Errorf("Error occurred when deleting permission %s", err)
		c.Error(errors.NewAppErr(NO_PERMISSION, err, err.Error()))
		return
	}

	c.JSON(http.StatusCreated, gin.H{})
}

// ValidateOperation - Validate operation
func ValidateOperation(c *gin.Context) {
	token := c.GetHeader(constant.AUTHORIZATION)
	userInfo, err := service.LoginSvc.ParseUserInfoFromToken(token)
	if err != nil {
		log.Errorf("Error occurred when parsing token %s", err)
		c.Error(errors.NewAppErr(INVALID_PARAMS, err, err.Error()))
		return
	}
	var req openapi.EnforcePermissionReq
	err = c.BindJSON(&req)
	if err != nil {
		log.Errorf("Error occurred when binding json %s", err)
		c.Error(errors.NewAppErr(INVALID_PARAMS, err, err.Error()))
	}

	hasAuth, err := psvc.PermissionSvc.Enforce(pvb.PermissionEnforce{
		Domain:       userInfo.Domain,
		Name:         userInfo.Name,
		ResourceType: req.ResourceType,
		ResourceId:   req.ResourceId,
		Action:       req.Action,
	})
	if err != nil {
		log.Errorf("Error occurred when enforcing permission %s", err)
		c.Error(errors.NewAppErr(UNDEFINED_ERROR, err, err.Error()))
		return
	}

	if !hasAuth {
		errMsg := fmt.Sprintf("User %s has no permission to %s resource %s/%s", userInfo.Name, req.Action, req.ServiceName, req.ResourceType)
		log.Error(errMsg)
		c.Error(errors.NewAppErr(NO_PERMISSION, err, err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
