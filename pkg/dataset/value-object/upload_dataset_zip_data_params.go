/*
 * Created on Fri Jul 14 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package valueobject

import "encoding/json"

type UploadDatasetZipDataParams struct {
	DatasetZipViewId string
	ZipFileName      string
}

func NewUploadDatasetZipDataParams(jsonStr string) (UploadDatasetZipDataParams, error) {
	params := UploadDatasetZipDataParams{}
	err := json.Unmarshal([]byte(jsonStr), &params)
	if err != nil {
		return params, err
	}
	return params, nil
}
