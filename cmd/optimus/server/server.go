/*
 * Created on Mon Jul 17 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */

package server

import (
	"fmt"
	"sync"

	"github.com/jacklv111/common-sdk/cli"
	aifsclient "github.com/jacklv111/common-sdk/client/aifs-client"
	"github.com/jacklv111/common-sdk/database"
	utilerrors "github.com/jacklv111/common-sdk/errors"
	"github.com/jacklv111/common-sdk/log"
	iamapi "github.com/jacklv111/optimus/app/iam/api"
	optimusapi "github.com/jacklv111/optimus/app/optimus/api"
	"github.com/jacklv111/optimus/cmd/optimus/config"
	gormadapter "github.com/jacklv111/optimus/infra/casbin/adapter"
	"github.com/jacklv111/optimus/infra/client/k8s"
	datasetscheduler "github.com/jacklv111/optimus/pkg/dataset/scheduler"
	"github.com/spf13/cobra"
)

var waitGroup sync.WaitGroup

// NewApiServerCommand creates a *cobra.Command object with default parameters
func NewApiServerCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "aifs-apiserver",
		Long: `Run aifs apiserver to provide rest operations`,

		// stop printing usage when the command errors
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			flagSet := cmd.Flags()

			// Activate logging as soon as possible, after that
			// show flags with the final logging configuration.
			if errs := log.ValidateAndApply(log.LogConfig); len(errs) > 0 {
				return utilerrors.NewAggregate(errs)
			}
			cli.PrintFlags(flagSet)

			// validate configs
			if errs := config.ServerConfig.Validate(); len(errs) > 0 {
				return utilerrors.NewAggregate(errs)
			}

			return run()
		},
	}

	cmdFlagSet := cmd.Flags()
	serverFlagSet := config.ServerConfig.GetFlags()
	cmdFlagSet.AddFlagSet(serverFlagSet)

	return cmd
}

// run starts servers, init components with configs.
func run() error {
	if err := database.InitDb(); err != nil {
		return err
	}

	// casbin gorm adapter, init after database
	// iam component depends on casbin
	if err := gormadapter.InitCasbinGormAdapter(); err != nil {
		return err
	}

	if err := aifsclient.InitAifsClientV2(); err != nil {
		return err
	}

	if err := k8s.InitK8sClient(); err != nil {
		return err
	}

	datasetscheduler.Start()

	waitGroup.Add(2)

	go startOptimusServer()

	go startIamServer()

	waitGroup.Wait()
	return nil
}

func startOptimusServer() {
	defer waitGroup.Done()
	router := optimusapi.NewRouter()
	err := router.Run(fmt.Sprintf(":%d", config.ServerConfig.GetOptimusPort()))
	if err != nil {
		log.Errorf("start optimus server error %s", err)
	}
}

func startIamServer() {
	defer waitGroup.Done()
	router := iamapi.NewRouter()
	err := router.Run(fmt.Sprintf(":%d", config.ServerConfig.GetIamPort()))
	if err != nil {
		log.Errorf("start iam server error %s", err)
	}
}
