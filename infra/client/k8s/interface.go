/*
 * Created on Mon Jul 17 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */
package k8s

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"path/filepath"

	"github.com/jacklv111/common-sdk/log"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const (
	ADMIN_SERVICE_ACCOUNT_PATH = "/admin/secret"
)

func InitK8sClient() (err error) {
	config := &rest.Config{
		Host: K8sConfig.ApiServerUrl,
		// Set other configuration options as needed, such as authentication options, timeouts, etc.
		TLSClientConfig: rest.TLSClientConfig{
			CAFile: filepath.Join(ADMIN_SERVICE_ACCOUNT_PATH, "ca.crt"),
		},
		BearerTokenFile: filepath.Join(ADMIN_SERVICE_ACCOUNT_PATH, "token"),
	}
	rest.InClusterConfig()
	// Create a Kubernetes clientset using the configuration
	Clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("error creating clientset: %w", err)
	}

	return nil
}

func getJobLogs(namespace, jobName string) (string, error) {
	podList, err := Clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: fmt.Sprintf("job-name=%s", jobName),
	})
	if err != nil {
		return "", fmt.Errorf("error listing Pods: %w", err)
	}

	var logs string
	for _, pod := range podList.Items {
		podLogs, err := Clientset.CoreV1().Pods(namespace).GetLogs(pod.Name, &corev1.PodLogOptions{}).Stream(context.TODO())
		if err != nil {
			return "", fmt.Errorf("error reading log stream for Pod %s: %w", pod.Name, err)
		}
		defer podLogs.Close()

		buf := new(bytes.Buffer)
		_, err = io.Copy(buf, podLogs)
		if err != nil {
			return "", fmt.Errorf("error reading log stream for Pod %s: %w", pod.Name, err)
		}

		logs += fmt.Sprintf("Logs for Pod %s:\n%s\n\n", pod.Name, buf.String())
	}

	return logs, nil
}

func DeleteJob(namespace, jobName string) error {
	// Delete the Job from the Kubernetes cluster
	err := Clientset.CoreV1().Pods(namespace).Delete(context.TODO(), jobName, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("error deleting Job: %w", err)
	}

	return nil
}

func IsJobCompleted(namespace, jobName string) bool {
	job, err := Clientset.CoreV1().Pods(namespace).Get(context.TODO(), jobName, metav1.GetOptions{})
	if err != nil {
		return false
	}

	if job.Status.Phase == corev1.PodSucceeded {
		return true
	} else if job.Status.Phase == corev1.PodFailed {
		logs, err := getJobLogs(namespace, jobName)
		if err != nil {
			log.Errorf("error getting Job logs: %v", err)
			return true
		}
		log.Errorf("job failed: %s, logs: %s", job.Status.Reason, logs)
		return true
	}

	return false
}

var Clientset *kubernetes.Clientset
