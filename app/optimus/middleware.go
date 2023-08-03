/*
 * Created on Thu Jul 06 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */
package optimus

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jacklv111/common-sdk/errors"
	"github.com/jacklv111/common-sdk/log"
	"github.com/jacklv111/optimus/app/iam/view-object/openapi"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next() // 先调用c.Next()执行后面的中间件
		// 所有中间件及router处理完毕后从这里开始执行
		// 检查c.Errors中是否有错误
		if len(c.Errors) == 0 {
			return
		}
		err := c.Errors[0].Err // 获取错误
		// 若是自定义的错误则将code、msg返回
		if appErr, ok := err.(errors.AppError); ok {
			log.Error("app error: %s", appErr.Error())
			msg := fmt.Sprintf(GetMsgTemplateByAppErrorCode(appErr.Code()), appErr.Args())
			c.JSON(GetHttpStatusByAppErrorCode(appErr.Code()),
				gin.H{
					"code": appErr.Code(),
					"msg":  msg,
				})
		} else {
			// 若非自定义错误则返回详细错误信息err.Error()
			// 比如save session出错时设置的err
			c.JSON(http.StatusOK, gin.H{
				"code": http.StatusInternalServerError,
				"msg":  err.Error(),
			})
		}

	}
}

func PanicHandler() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, err interface{}) {
		if strings.Contains(err.(string), "invalid UUID") {
			c.JSON(http.StatusBadRequest, openapi.Error{Code: INVALID_PARAMS, Message: err.(string)})
		}
	})
}
