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

	"github.com/jacklv111/optimus/infra/client/k8s"
	corev1 "k8s.io/api/core/v1"
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
	job := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: jobName,
		},

		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "hello-world",
					Image: "swr.cn-south-1.myhuaweicloud.com/jacklv/helloworld:v0.1",
				},
			},
			RestartPolicy: corev1.RestartPolicyNever,
		},
	}
	// Create the Job in the Kubernetes cluster
	_, err = k8s.Clientset.CoreV1().Pods(namespace).Create(context.TODO(), job, metav1.CreateOptions{})
	if err != nil {
		fmt.Printf("create job error: %s", err)
	}
}
