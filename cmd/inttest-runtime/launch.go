package main

import (
	"inttest-runtime/internal/config"
	mockRpcService "inttest-runtime/internal/domain/service/mockrpc"
	configRepo "inttest-runtime/internal/repository/config"
	"log"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func launchServices(cmd *cobra.Command, args []string) {
	confPath, err := cmd.PersistentFlags().GetString(string(configPathArg))
	if err != nil {
		err := errors.Wrap(err, "error getting config file path")
		log.Fatal(err)
	}
	config, err := config.FromFile(confPath)
	if err != nil {
		log.Fatal(err)
	}
	if config == nil {
		log.Fatal("runtime config is empty")
	}

	configRepo := configRepo.NewLocalRepository(*config)
	httpMocksService := mockRpcService.New(configRepo)
}
