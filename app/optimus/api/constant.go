/*
 * Created on Thu Jul 13 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package api

import "math"

const (
	OFFSET_STR = "offset"
	LIMIT_STR  = "limit"
)

// api parameters
const (
	DATASET_ID             = "datasetId"
	ZIP_FORMAT             = "X-Zip-Format"
	WORKSPACE              = "X-Workspace"
	DEFAULT_WORKSPACE      = "default"
	X_ZIP_FILE_NAME        = "X-Zip-File-Name"
	RAW_DATA_ID_LIST       = "rawDataIdList"
	ANNOTATION_ID_LIST     = "annotationIdList"
	OFFSET_MIN             = 0
	OFFSET_MAX             = math.MaxInt
	LIMIT_MIN              = 1
	LIMIT_MAX              = 50
	OFFSET_DEFAULT         = 0
	LIMIT_DEFAULT          = 10
	POOL_NAME              = "poolName"
	VERSION_NAME           = "versionName"
	VERSION_PARTITION_NAME = "versionPartitionName"
	HAS_ANNOTATION_FILTER  = "hasAnnotationFilter"
	LABEL_ID               = "labelId"
	ANNOTATION_TEMPLATE_ID = "annotationTemplateId"

	SORT_BY            = "sortBy"
	DEFAULT_SORT_BY    = "createdAt"
	SORT_ORDER         = "sortOrder"
	DEFAULT_SORT_ORDER = "desc"
)
