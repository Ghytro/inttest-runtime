package rest

import (
	"context"

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
	reqParams domainTypes.RestClientRequestParams,
) (response *domainTypes.RestLogicResponse, err error) {
	behavior, err := s.repo.GetHttpServiceBehaviorByUrlMethod(ctx, config.RpcServiceType_REST, urlPattern, method)
	if err != nil {
		return nil, err
	}

	resp, err := domainTypes.PerformRestLogic(reqParams, behavior)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s Service) HandleSoapRequest(
	ctx context.Context,
	urlPattern string,
	method string,
	reqParams domainTypes.SoapClientRequestParams,
) (*domainTypes.SoapLogicResponse, error) {
	behavior, err := s.repo.GetHttpServiceBehaviorByUrlMethod(ctx, config.RpcServiceType_REST, urlPattern, method)
	if err != nil {
		return nil, err
	}

	resp, err := domainTypes.PerformSoapLogic(reqParams, behavior)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
