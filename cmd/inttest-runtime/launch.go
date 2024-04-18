package main

import (
	"context"
	"fmt"
	mockRpcApi "inttest-runtime/internal/api/mockrpc"
	"inttest-runtime/internal/config"
	mockRpcService "inttest-runtime/internal/domain/service/mockrpc"
	domainTypes "inttest-runtime/internal/domain/types"
	configRepo "inttest-runtime/internal/repository/config"
	"inttest-runtime/pkg/embedded"
	"inttest-runtime/pkg/mq"
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
	pyRuntime, err := embedded.NewPythonRuntime()
	if err != nil {
		log.Fatal(err)
	}
	pyFuncExecutor, err := compilePyFuncs(*cfg, pyRuntime)
	if err != nil {
		log.Fatal(err)
	}

	httpMockExecutor := domainTypes.NewMockLogicExecutor(pyFuncExecutor)
	httpMocksService := mockRpcService.New(configRepo, httpMockExecutor)

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
			return mockApi.Listen(ctx, fmt.Sprintf(":%d", rpcService.GetPort()))
		})
	}

	for _, broker := range cfg.Brokers {
		var client domainTypes.IMockBrokerPubSub
		switch broker.Type {
		case config.BrokerType_REDIS_PUBSUB:
			localRedis, err := mq.NewLocalRedis()
			if err != nil {
				log.Fatal(err)
			}
			redisAddr := fmt.Sprintf(":%d", broker.GetPort())
			errGroup.Go(func() error {
				return localRedis.Listen(ctx, redisAddr)
			})
			client, err = mq.ConnectRedisPubSub(redisAddr, 0, "")
			if err != nil {
				log.Fatal(err)
			}
		default:
			log.Fatalf("unknown broker type: %s", broker.Type)
		}

		brokerLogic := domainTypes.NewMockBroker(pyFuncExecutor, client)
		errGroup.Go(func() error {
			return brokerLogic.Start(ctx)
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

func compilePyFuncs(cfg config.Config, pyRuntime *embedded.PyRuntime) (*domainTypes.PyPrecompiledExecutor, error) {
	compileStore := domainTypes.NewPyPrecompiledExecutor(pyRuntime)
	for _, rpcService := range cfg.RpcServices {
		for _, route := range rpcService.RpcServiceUnion.HttpService.Routes {
			for _, behav := range route.Behavior {
				if behav.Type != config.RestHandlerBehaviorType_MOCK {
					continue
				}
				if _, err := compileStore.AddFunc(behav.HttpMockBehavior.Impl); err != nil {
					return nil, err
				}
			}
		}
	}
	for _, broker := range cfg.Brokers {
		for _, behav := range broker.BrokerBehaviorUnion.BrokerBehaviorRedis.Behavior {
			for _, generator := range behav.Generators {
				if generator.Type != config.RedisTopicGeneratorType_PROG {
					continue
				}
				if _, err := compileStore.AddFunc(generator.RedisTopicGeneratorUnion.RedisTopicGeneratorProg.Behavior); err != nil {
					return nil, err
				}
			}
		}
	}
	return compileStore, nil
}
