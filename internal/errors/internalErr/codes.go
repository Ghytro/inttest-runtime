package internalErr

type ErrCode int

const (
	ErrCodeConfigurationParsing ErrCode = 1
)

var explains = []string{
	ErrCodeConfigurationParsing: "Error while parsing configuration",
}
