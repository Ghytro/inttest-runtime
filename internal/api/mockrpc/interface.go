package mockrpc

import (
	"context"
	domainTypes "inttest-runtime/internal/domain/types"
)

type IRestApiService interface {
	HandleRestRequest(ctx context.Context, urlPattern string, method string, reqParams domainTypes.HttpClientRequestParams) (response *domainTypes.HttpResp, err error)
}

type ISoapApiService interface {
	HandleSoapRequest(ctx context.Context, urlPattern string, method string, reqParams domainTypes.HttpClientRequestParams) (*domainTypes.HttpResp, error)
}
