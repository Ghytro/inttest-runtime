package internalErr

type ErrCode int

const (
	ErrCodeConfigurationParsing ErrCode = iota + 1
	ErrCodeServiceGet
	ErrCodeServiceCreate
	ErrCodeRestServiceCreate
	ErrCodeGrpcServiceCreate
	ErrCodeJSONParsing
)

var explains = []string{
	ErrCodeConfigurationParsing: "Error while parsing configuration",
	ErrCodeServiceGet:           "Error while getting info about service",
	ErrCodeServiceCreate:        "Error writing info about service",
	ErrCodeRestServiceCreate:    "Error writing info about REST service",
	ErrCodeGrpcServiceCreate:    "Error writing info about gRPC service",
	ErrCodeJSONParsing:          "Error while parsing JSON",
}
