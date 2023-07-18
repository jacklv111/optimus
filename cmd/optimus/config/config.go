/*
 * Created on Mon Jul 17 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package config

import (
	aifsclient "github.com/jacklv111/common-sdk/client/aifs-client"
	basecfg "github.com/jacklv111/common-sdk/config"
	"github.com/jacklv111/common-sdk/database"
	"github.com/jacklv111/common-sdk/database/gorm"
	"github.com/jacklv111/common-sdk/env"
	"github.com/jacklv111/common-sdk/log"
	casbincfg "github.com/jacklv111/optimus/infra/casbin/config"
	"github.com/jacklv111/optimus/infra/client/k8s"
	"github.com/spf13/pflag"
)

// ServerConfig aifs api server config.
type serverConfig struct {
	basecfg.Configs
	optimusPort int
	iamPort     int
}

var ServerConfig serverConfig

func init() {
	ServerConfig = serverConfig{}

	ServerConfig.AddConfig(env.EnvConfig)
	ServerConfig.AddConfig(log.LogConfig)
	ServerConfig.AddConfig(gorm.GormConfig)
	ServerConfig.AddConfig(database.DbConfig)
	ServerConfig.AddConfig(aifsclient.AifsConfig)
	ServerConfig.AddConfig(casbincfg.CasbinConfig)
	ServerConfig.AddConfig(k8s.K8sConfig)
}

// optimusPort getter
func (serverCfg serverConfig) GetOptimusPort() int {
	return serverCfg.optimusPort
}

// iamPort getter
func (serverCfg serverConfig) GetIamPort() int {
	return serverCfg.iamPort
}

func (serverCfg serverConfig) GetFlags() (flagSet *pflag.FlagSet) {
	flagSet = serverCfg.Configs.GetFlags()
	flagSet.IntVar(&ServerConfig.optimusPort, "optimus-port", 8080, "optimus server port.")
	flagSet.IntVar(&ServerConfig.iamPort, "iam-port", 8081, "iam server port.")

	return
}
