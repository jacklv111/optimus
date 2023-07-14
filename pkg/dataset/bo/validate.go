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

	"github.com/jacklv111/common-sdk/log"
	"github.com/jacklv111/optimus/pkg/dataset"
)

func (bo *DatasetBo) updatePoolValidate(poolName string) error {
	if !bo.ExistsPoolByName(poolName) {
		log.Errorf("dataset id %s pool %s not found", bo.datasetDo.ID, poolName)
		return dataset.ErrNotFound
	}

	hasRunningActionOnPool, err := bo.hasRunningActionOnPool(poolName)
	if err != nil {
		return err
	}
	if hasRunningActionOnPool {
		return errors.New("pool has a running action, can not update pool")
	}
	return nil
}
