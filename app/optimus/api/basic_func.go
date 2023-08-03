/*
 * Created on Thu Jul 13 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jacklv111/common-sdk/log"
	"github.com/jacklv111/optimus/app/optimus/view-object/openapi"
	"github.com/jacklv111/optimus/pkg/iam/constant"
	"github.com/jacklv111/optimus/pkg/iam/login/service"
	vb "github.com/jacklv111/optimus/pkg/iam/login/value-object"
)

func getUserInfoAndHandleError(c *gin.Context) (userInfo vb.UserInfo, err error) {
	token := c.GetHeader(constant.AUTHORIZATION)
	userInfo, err = service.LoginSvc.ParseUserInfoFromToken(token)
	if err != nil {
		log.Errorf("Error occurred when parsing token %s", err)
		// todo 修改错误码和 msg
		c.JSON(http.StatusUnauthorized, openapi.Error{Code: "1", Message: err.Error()})
		return
	}
	return
}
