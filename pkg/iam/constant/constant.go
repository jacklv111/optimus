/*
 * Created on Thu Jul 13 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package constant

const (
	ADMIN                    = "admin"
	AUTHORIZATION            = "Authorization"
	TOKEN_EXPIRE_TIME_IN_SEC = 3600 * 24
	CASBIN_MODEL_CONF        = "conf/casbin/model"
	// todo: change this to env variable
	JWT_SECRET = "9SQzwmbp2urd0MD_IkHJ9hxj7wUCXrgw7XF9u3UONgI"
)

func GetSecret() []byte {
	return []byte(JWT_SECRET)
}

const (
	RESOURCE_TYPE_DATASET             = "dataset/dataset"
	RESOURCE_TYPE_VERSION             = "dataset/version"
	RESOURCE_TYPE_POOL                = "dataset/pool"
	RESOURCE_TYPE_ANNOTATION_TEMPLATE = "dataset/annotation-template"
)

// action
const (
	AUTHORIZE          = "authorize"
	CREATE             = "create"
	DELETE             = "delete"
	GET_DETAILS        = "getDetails"
	GET_LIST           = "getList"
	UPDATE             = "update"
	GET_DATA_ITEM_LIST = "getDataItemList"
	DELETE_DATA_ITEM   = "deleteDataItem"
)

const (
	EFFECT_ALLOW = "allow"
)
