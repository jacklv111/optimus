/*
 * Created on Fri Jul 14 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package constant

const (
	UPLOAD_DATASET_VERSION_ZIP = "upload-dataset-version-zip"
	UPLOAD_RAW_DATA_ZIP        = "upload-raw-data-zip"
	DATASET_ZIP_DECOMPRESSION  = "dataset-zip-decompression"
	GET_ACTION_INTERVAL_IN_SEC = 30

	HAS_ANNOTATION_FILTER_ALL               = "all"
	HAS_ANNOTATION_FILTER_HAS_ANNOTATION    = "hasAnnotation"
	HAS_ANNOTATION_FILTER_HAS_NO_ANNOTATION = "hasNoAnnotation"

	// version partition name
	VERSION_PARTITION_NAME_TRAIN = "train"
	VERSION_PARTITION_NAME_TEST  = "test"
	VERSION_PARTITION_NAME_VAL   = "val"

	// env var
	AIFS_IP        = "AIFS_IP"
	AIFS_PORT      = "AIFS_PORT"
	S3_BUCKET_NAME = "S3_BUCKET_NAME"
	S3_AK          = "S3_AK"
	S3_SK          = "S3_SK"
	S3_ENDPOINT    = "S3_ENDPOINT"
	S3_REGION      = "S3_REGION"
)
