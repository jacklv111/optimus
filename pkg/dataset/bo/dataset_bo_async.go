/*
 * Created on Fri Jul 14 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package bo

import (
	"context"
	"fmt"
	"time"

	aifsclient "github.com/jacklv111/common-sdk/client/aifs-client"
	"github.com/jacklv111/common-sdk/env"
	"github.com/jacklv111/common-sdk/log"
	"github.com/jacklv111/common-sdk/scheduler"
	"github.com/jacklv111/common-sdk/utils"
	"github.com/jacklv111/optimus/infra/action"
	"github.com/jacklv111/optimus/infra/client/k8s"
	"github.com/jacklv111/optimus/pkg/dataset/constant"
	dsvb "github.com/jacklv111/optimus/pkg/dataset/value-object"
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func startDecompressionAction(namespace, jobName, datasetZipViewId string) error {
	// Create a Job object
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name: jobName,
		},

		Spec: batchv1.JobSpec{
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name: "decompression",
							// TBD
							Image:   "your-container-image",
							Command: []string{"your-command"},
						},
					},
					RestartPolicy: v1.RestartPolicyNever,
				},
			},
			BackoffLimit: utils.Int32Ptr(0), // Optional: Set the backoff limit
		},
	}

	// Create the Job in the Kubernetes cluster
	_, err := k8s.Clientset.BatchV1().Jobs(namespace).Create(context.TODO(), job, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("error creating decompress job: %w", err)
	}

	return nil

}

func waitActionFinished(actionRecord action.ActionDo, params dsvb.UploadDatasetZipDataParams) {
	namespace := env.EnvConfig.GetEnvType()
	jobName := fmt.Sprintf("decompression-%s", params.DatasetZipViewId)
	err := startDecompressionAction(namespace, jobName, params.DatasetZipViewId)
	if err != nil {
		log.Errorf("Failed to StartDecompressionAction: %v", err)
		return
	}
	// wait for action finished
	scheduler.WaitCondition(constant.GET_ACTION_INTERVAL_IN_SEC, func() bool {
		actionRecord.UpdateAt = time.Now().Unix()
		err = action.ActionMgr.Updates(actionRecord)
		if err != nil {
			log.Errorf("Failed to Updates: %v", err)
			return false
		}

		return k8s.IsJobCompleted(namespace, jobName)
	})
	k8s.DeleteJob(namespace, jobName)
}

func (bo *DatasetBo) UploadDatasetVersionZipAsync(versionName string, actionRecord action.ActionDo, params dsvb.UploadDatasetZipDataParams) {
	log.Infof("Start to UploadDatasetVersionZip: action record: %v, input: %v", actionRecord, params)

	waitActionFinished(actionRecord, params)
	// 后处理
	details, _, err := aifsclient.GetAifsClient().DataViewApi.GetDataViewDetails(context.Background(), params.DatasetZipViewId).Execute()
	if err != nil {
		log.Errorf("Failed to GetDataViewDetails for dataset-zip view: %v", err)
		return
	}

	trainRawDataStatistics, _, err := aifsclient.GetAifsClient().DataViewApi.GetDataViewStatistics(context.Background(), *details.TrainRawDataViewId).Execute()
	if err != nil {
		log.Errorf("Failed to GetDataViewStatistics for train raw data view: %v", err)
		return
	}
	valRawDataStatistics, _, err := aifsclient.GetAifsClient().DataViewApi.GetDataViewStatistics(context.Background(), *details.ValRawDataViewId).Execute()
	if err != nil {
		log.Errorf("Failed to GetDataViewStatistics for val raw data view: %v", err)
		return
	}

	trainAnnotationDetails, _, err := aifsclient.GetAifsClient().DataViewApi.GetDataViewDetails(context.Background(), *details.TrainAnnotationViewId).Execute()
	if err != nil {
		log.Errorf("Failed to GetDataViewDetails for train annotation view: %v", err)
		return
	}
	totalCount := trainRawDataStatistics.GetItemCount() + valRawDataStatistics.GetItemCount()
	err = bo.CreateNewVersionFromExistedView(dsvb.GenVersionFromExistedViewParams{
		Name:                  versionName,
		Description:           fmt.Sprintf("decompressed from %s", params.ZipFileName),
		TrainRawDataViewId:    details.GetTrainRawDataViewId(),
		TrainAnnotationViewId: details.GetTrainAnnotationViewId(),
		TrainRawDataNum:       trainRawDataStatistics.GetItemCount(),
		TrainTotalDataSize:    trainRawDataStatistics.GetTotalDataSize(),
		TrainRawDataRatio:     float32(trainRawDataStatistics.GetItemCount()) / float32(totalCount),
		ValRawDataViewId:      details.GetValRawDataViewId(),
		ValAnnotationViewId:   details.GetValAnnotationViewId(),
		ValRawDataNum:         valRawDataStatistics.GetItemCount(),
		ValTotalDataSize:      valRawDataStatistics.GetTotalDataSize(),
		ValRawDataRatio:       float32(valRawDataStatistics.GetItemCount()) / float32(totalCount),
	})
	if !bo.HasAnnotationTemplate() {
		bo.UpdateDataset(dsvb.DatasetUpdateParams{
			AnnotationTemplateId:   trainAnnotationDetails.GetAnnotationTemplateId(),
			AnnotationTemplateType: trainAnnotationDetails.GetAnnotationTemplateType(),
		})
		if err != nil {
			log.Errorf("Failed to GenerateNewVersionFromExistedView: %v", err)
			return
		}
	}

	// hard delete zip data view
	_, err = aifsclient.GetAifsClient().DataViewApi.HardDeleteDataView(context.Background(), params.DatasetZipViewId).Execute()
	if err != nil {
		log.Errorf("Failed to HardDeleteDataView: %v", err)
		return
	}
	// delete action record
	err = action.ActionMgr.Delete(actionRecord.ResourceType, actionRecord.ResourceId)
	if err != nil {
		log.Errorf("Failed to Delete action record: %v", err)
		return
	}
}

func (bo *DatasetBo) UploadRawDataZipToPoolAsync(actionRecord action.ActionDo, params dsvb.UploadDatasetZipDataParams) {
	log.Infof("Start to UploadRawDataZip: action record: %v, input: %v", actionRecord, params)

	waitActionFinished(actionRecord, params)
	// 后处理
	// hard delete zip data view
	_, err := aifsclient.GetAifsClient().DataViewApi.HardDeleteDataView(context.Background(), params.DatasetZipViewId).Execute()
	if err != nil {
		log.Errorf("Failed to HardDeleteDataView: %v", err)
		return
	}
	// delete action record
	err = action.ActionMgr.Delete(actionRecord.ResourceType, actionRecord.ResourceId)
	if err != nil {
		log.Errorf("Failed to Delete action record: %v", err)
		return
	}
}
