/*
 * Created on Wed Aug 02 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */
package main

import (
	"context"
	"fmt"

	"github.com/jacklv111/common-sdk/utils"
	"github.com/jacklv111/optimus/infra/client/k8s"
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func main() {
	k8s.K8sConfig.ApiServerUrl = "https://192.168.0.185:5443"

	err := k8s.InitK8sClient()
	if err != nil {
		fmt.Printf("k8s init error: %s", err)
	}
	namespace := "dev"
	jobName := "hello-world"
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
							Name:  "hello-world",
							Image: "swr.cn-south-1.myhuaweicloud.com/jacklv/helloworld:v0.1",
						},
					},
					RestartPolicy: v1.RestartPolicyNever,
				},
			},
			BackoffLimit: utils.Int32Ptr(0), // Optional: Set the backoff limit
		},
	}
	// Create the Job in the Kubernetes cluster
	_, err = k8s.Clientset.BatchV1().Jobs(namespace).Create(context.TODO(), job, metav1.CreateOptions{})
	if err != nil {
		fmt.Printf("create job error: %s", err)
	}
}
