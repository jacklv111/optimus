/*
 * Created on Thu Jul 13 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package config

import (
	"github.com/spf13/pflag"
)

const (
	CASBIN_MODEL_CONFIG_PATH = "./conf/casbin/model"
)

type casbinConfig struct {
	ModelPath string `yaml:"model_path"`
}

var CasbinConfig *casbinConfig

func init() {
	CasbinConfig = &casbinConfig{
		ModelPath: CASBIN_MODEL_CONFIG_PATH,
	}
}

func (config *casbinConfig) ReadFromFile() error {
	return nil
}

func (config *casbinConfig) AddFlags(flagSet *pflag.FlagSet) {
	flagSet.StringVar(&config.ModelPath, "casbin-model-path", CASBIN_MODEL_CONFIG_PATH, "Value to indicate the casbin model path")
}

func (config casbinConfig) Validate() []error {
	// do nothing
	var errs []error
	return errs
}
