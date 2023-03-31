# Protocol Recruiting Challenge

## Implement Osmosis Proxy Server

This project demonstrates the implementation of an Osmosis proxy server using the Tendermint and Cosmos SDK packages.

## Implementation

The implementation relies on typical protobuf and gRPC packages, Tendermint, and Cosmos SDK packages. The required proto files have been copied from Cosmos SDK and Tendermint, and generated in the `/pkg` folder. A proxy server and tracker module have been built based on these proto files.

While the proxy server and tracker can be registered as microservices with a more sophisticated structure, it's not considered critical for this project.

## How to Check the Result

1. Run the server:

`make run-server`
After starting the server, open another terminal.

2. Run the server test:
   `make run-server-test`

The test checks two endpoints specified in the problem statement. Once the server collects 5 blocks, it registers the data in the `app/test_result.json` file.

## How to Call with grpcurl

- List All methods:

```
grpcurl -plaintext
-import-path ./proto
-import-path ./proto/gogoproto
-import-path ./proto/tendermint/types
-import-path ./proto/tendermint/version
-import-path ./proto/google/protobuf
-import-path ./proto/amino
-proto proto/cosmos/base/tendermint/v1beta1/query.proto
localhost:9000 list cosmos.base.tendermint.v1beta1.Service
```

- Call `GetNodeInfo` endpoint:
  ```
  grpcurl -plaintext
  -import-path ./proto
  -import-path ./proto/gogoproto
  -import-path ./proto/tendermint/types
  -import-path ./proto/tendermint/version
  -import-path ./proto/google/protobuf
  -import-path ./proto/amino
  -proto proto/cosmos/base/tendermint/v1beta1/query.proto
  -d '{}'
  localhost:9000 cosmos.base.tendermint.v1beta1.Service/GetNodeInfo
  ```
- Call `GetLatestBlock`:

```
grpcurl -plaintext
-import-path ./proto
-import-path ./proto/gogoproto
-import-path ./proto/tendermint/types
-import-path ./proto/tendermint/version
-import-path ./proto/google/protobuf
-import-path ./proto/amino
-proto proto/cosmos/base/tendermint/v1beta1/query.proto
-d '{}'
localhost:9000 cosmos.base.tendermint.v1beta1.Service/GetLatestBlock
```

## Note

If you want to generate proto file from scratch, please run follow command.
`make proto-gen`
It will generate types ./pkg folder.

This project was built for mock. There are many improvement points. 

