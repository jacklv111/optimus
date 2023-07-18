/*
 * Created on Mon Jul 17 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */
package k8s

import (
	"github.com/spf13/pflag"
)

type k8sConfig struct {
	ApiServerUrl string
}

var K8sConfig *k8sConfig

func init() {
	K8sConfig = &k8sConfig{
		ApiServerUrl: "http://",
	}
}

func (cfg *k8sConfig) ReadFromFile() error {
	return nil
}

func (cfg *k8sConfig) AddFlags(flagSet *pflag.FlagSet) {
	flagSet.StringVar(&cfg.ApiServerUrl, "k8s-api-server-url", cfg.ApiServerUrl, "k8s api server url")
}

func (cfg *k8sConfig) Validate() []error {
	return []error{}
}
