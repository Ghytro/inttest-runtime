package mockrpc

import (
	"context"
	domainTypes "inttest-runtime/internal/domain/types"
)

type IRestApiService interface {
	HandleRestRequest(
		ctx context.Context,

		urlPattern string,
		method string,
		reqParams domainTypes.RestClientRequestParams,
	) (response *domainTypes.RestLogicResponse, err error)
}

type ISoapApiService interface {
	HandleSoapRequest(
		ctx context.Context,
		urlPattern string,
		method string,
		reqParams domainTypes.SoapClientRequestParams,
	) (*domainTypes.SoapLogicResponse, error)
}
