/*
 * Created on Fri Jul 14 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package scheduler

import (
	"time"

	"github.com/jacklv111/common-sdk/log"
	"github.com/jacklv111/common-sdk/scheduler"
	"github.com/jacklv111/common-sdk/scheduler/shedlock"
	"github.com/jacklv111/optimus/infra/action"
	dsbo "github.com/jacklv111/optimus/pkg/dataset/bo"
	valueobject "github.com/jacklv111/optimus/pkg/dataset/value-object"
	iamconst "github.com/jacklv111/optimus/pkg/iam/constant"
)

func Start() {
	log.Info("action scheduler started")

	scheduleConfig := scheduler.ScheduleConfig{
		Name:         "action-recover",
		Interval:     time.Minute * 5,
		InitialDelay: time.Minute * 3,
		Runnable:     recover,
	}
	shedlockConfig := shedlock.ShedlockConfig{
		Enabled:        true,
		Name:           "action-recover-lock",
		LockAtLeastFor: time.Minute * 5,
		LockAtMostFor:  time.Minute * 10,
	}
	scheduler.Schedule(shedlockConfig, scheduleConfig)
}

// 恢复一些由于异常结束的任务
func recover() {
	aList, err := action.ActionMgr.GetUpdateAtLessThan(time.Now().Add(-time.Minute * 5).UnixMilli())
	if err != nil {
		log.Errorf("get action failed, err: %v", err)
		return
	}
	if len(aList) == 0 {
		log.Info("no action need to recover")
	}
	for _, actionRecord := range aList {
		switch actionRecord.ResourceType {
		case iamconst.RESOURCE_TYPE_POOL:
			datasetId, poolName, err := dsbo.ParseResourceId(actionRecord.ResourceId)
			if err != nil {
				log.Errorf("Failed to parse datasetId and poolName from resource id: %s", actionRecord.ResourceId)
				continue
			}
			log.Infof("recover dataset id %s, pool name %s, action: %v", datasetId, poolName, actionRecord)

			params, err := valueobject.NewUploadDatasetZipDataParams(actionRecord.Params)
			if err != nil {
				log.Errorf("Failed to parse params: %v", err)
				continue
			}
			datasetBo := dsbo.BuildWithId(datasetId)
			go datasetBo.UploadRawDataZipToPoolAsync(actionRecord, params)
		case iamconst.RESOURCE_TYPE_VERSION:
			datasetId, versionName, err := dsbo.ParseResourceId(actionRecord.ResourceId)
			if err != nil {
				log.Errorf("Failed to parse datasetId and versionName from resource id: %s", actionRecord.ResourceId)
				continue
			}
			log.Infof("recover dataset id %s, version name %s, action: %v", datasetId, versionName, actionRecord)

			params, err := valueobject.NewUploadDatasetZipDataParams(actionRecord.Params)
			if err != nil {
				log.Errorf("Failed to parse params: %v", err)
				continue
			}
			datasetBo := dsbo.BuildWithId(datasetId)
			go datasetBo.UploadDatasetVersionZipAsync(versionName, actionRecord, params)
		}
		log.Errorf("cant handle resource type: %v, action %v", actionRecord.ResourceType, actionRecord)
	}
}
