/*
 * Created on Fri Jul 14 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package bo

import (
	"errors"
	"strings"
)

func ParseResourceId(resourceId string) (datasetId, resourceName string, err error) {
	stList := strings.Split(resourceId, "/")
	if len(stList) != 2 {
		return "", "", errors.New("invalid resource id")
	}
	return stList[0], stList[1], nil
}
