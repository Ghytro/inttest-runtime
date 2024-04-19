start.local:
	go run ./cmd/inttest-runtime --config ./config/conf.local.json

start.test-consumer:
	go run ./cmd/test-consumer
