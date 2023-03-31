# Generate Go code from Protobuf files
proto-gen:
	protoc -I ./proto --go_out=./pkg --go_opt=paths=source_relative --go-grpc_out=./pkg --go-grpc_opt=paths=source_relative ./proto/**/*.proto

run-server:
	cd app; go run main.go

run-server-test:
	go run test/*.go