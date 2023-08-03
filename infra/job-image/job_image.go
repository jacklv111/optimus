/*
 * Created on Thu Aug 03 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */
package jobimage

import (
	"os"
	"path/filepath"
)

const (
	JOB_PATH = "/job"

	DECOMPRESS = "decompress"
)

func Get(jobName string) (imageName string, err error) {
	filePath := filepath.Join(JOB_PATH, jobName)
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	imageName = string(content)
	return imageName, nil
}
