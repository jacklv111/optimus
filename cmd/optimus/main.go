/*
 * Created on Mon Jul 17 2023
 *
 * Copyright (c) 2023 Company-placeholder. All rights reserved.
 *
 * Author Yubinlv.
 */
package main

import (
	"fmt"
	_ "net/http/pprof"
	"os"

	"github.com/jacklv111/common-sdk/cli"
	"github.com/jacklv111/common-sdk/utils"
	"github.com/jacklv111/optimus/cmd/optimus/config"
	"github.com/jacklv111/optimus/cmd/optimus/server"
)

func main() {
	utils.StartProf()
	// pre run
	// read config from file
	if err := config.ServerConfig.ReadFromFile(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	command := server.NewApiServerCommand()
	code := cli.Run(command)
	os.Exit(code)
}
