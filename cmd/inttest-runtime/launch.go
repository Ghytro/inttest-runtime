package main

import (
	"context"
	"fmt"
	mockRpcApi "inttest-runtime/internal/api/mockrpc"
	"inttest-runtime/internal/config"
	mockRpcService "inttest-runtime/internal/domain/service/mockrpc"
	configRepo "inttest-runtime/internal/repository/config"
	"log"
	"os"
	"os/signal"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

func launchServices(cmd *cobra.Command, args []string) {
	confPath, err := cmd.PersistentFlags().GetString(string(configPathArg))
	if err != nil {
		err := errors.Wrap(err, "error getting config file path")
		log.Fatal(err)
	}
	cfg, err := config.FromFile(confPath)
	if err != nil {
		log.Fatal(err)
	}
	if cfg == nil {
		log.Fatal("runtime config is empty")
	}

	configRepo := configRepo.NewLocalRepository(*cfg)
	httpMocksService := mockRpcService.New(configRepo)

	errGroup, ctx := errgroup.WithContext(context.Background())
	errGroup.Go(func() error {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, os.Interrupt)
		<-sigs
		return errors.New("captured sigint, gracefully shutting down...")
	})
	for _, rpcService := range cfg.RpcServices {
		rpcService := rpcService
		var mockApi mockRpcApi.IMockApi
		switch rpcService.Type {
		case config.RpcServiceType_REST:
			mockApi = mockRpcApi.NewRestMockApi(httpMocksService)
		case config.RpcServiceType_SOAP:
			mockApi = mockRpcApi.NewSoapMockApi(httpMocksService)
		}
		if err := registerRpcRoutes(mockApi, rpcService.RpcServiceUnion.HttpService.Routes...); err != nil {
			log.Fatal(err)
		}
		errGroup.Go(func() error {
			return mockApi.Listen(ctx, fmt.Sprintf(":%d", rpcService.Port))
		})
	}

	if err := errGroup.Wait(); err != nil {
		log.Fatal(err)
	}
}

func registerRpcRoutes(api mockRpcApi.IMockApi, routes ...config.HttpRoute) error {
	for _, r := range routes {
		if err := api.Register(r.Route.String(), string(r.Method)); err != nil {
			return err
		}
	}
	return nil
}
