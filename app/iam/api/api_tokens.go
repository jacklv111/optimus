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
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jacklv111/common-sdk/errors"
	"github.com/jacklv111/common-sdk/log"
	"github.com/jacklv111/optimus/app/iam/view-object/openapi"
	"github.com/jacklv111/optimus/pkg/iam/constant"
	"github.com/jacklv111/optimus/pkg/iam/login/service"
)

// GetTokenInfo - Get token info
func GetTokenInfo(c *gin.Context) {
	token := c.GetHeader(constant.AUTHORIZATION)
	userInfo, err := service.LoginSvc.ParseUserInfoFromToken(token)
	if err != nil {
		log.Errorf("Error occurred when parsing token %s", err)
		c.Error(errors.NewAppErr(INVALID_PARAMS, err, err.Error()))
		return
	}
	var resp openapi.TokenInfo
	resp.Username = userInfo.Name
	resp.ExpiredAt = time.Unix(userInfo.ExpiresAt, 0)
	c.JSON(http.StatusOK, resp)
}
