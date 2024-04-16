package rest

import (
	"context"
	"errors"

	"inttest-runtime/internal/config"
	domainTypes "inttest-runtime/internal/domain/types"
)

type Service struct {
	repo Repository
}

func New(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s Service) HandleRestRequest(
	ctx context.Context,

	urlPattern string,
	method string,
	reqParams domainTypes.HttpClientRequestParams,
) (response *domainTypes.HttpResp, err error) {
	return s.handleHttpRequestImpl(ctx, config.RpcServiceType_REST, urlPattern, method, reqParams)
}

func (s Service) HandleSoapRequest(
	ctx context.Context,
	urlPattern string,
	method string,
	reqParams domainTypes.HttpClientRequestParams,
) (*domainTypes.HttpResp, error) {
	return s.handleHttpRequestImpl(ctx, config.RpcServiceType_SOAP, urlPattern, method, reqParams)
}

func (s Service) handleHttpRequestImpl(
	ctx context.Context,
	rpcType config.RpcServiceType,
	urlPattern string,

	method string,
	reqParams domainTypes.HttpClientRequestParams,
) (*domainTypes.HttpResp, error) {
	behavior, err := s.repo.GetHttpServiceBehaviorByUrlMethod(ctx, rpcType, urlPattern, method)
	if err != nil {
		return nil, err
	}

	result, err := domainTypes.PerformHttpLogicItem(
		reqParams,
		config.RpcServiceType_REST,
		behavior,
	)
	if err != nil {
		return nil, err
	}

	// я хз зачем это сделал но пусть будет
	if result.Mock != nil {
		return result.Mock, nil
	}
	if result.Stub != nil {
		return result.Stub, nil
	}
	return nil, errors.New("no result to return")
}
