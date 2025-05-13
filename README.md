# Board Microservice Codespace

## Setting up gRPC-Gateway

Ensure that you include http transcoding inside proto files (more examples in `/api/proto/v1/board.proto`):

```protobuf
service BoardService {
  // example of http transcoding
  rpc CreateBoard(CreateBoardRequest) returns (BoardResponse) {
    option (google.api.http) = {
      post: "/v1/boards"
      body: "*"
    };
  }

  // ...
}
```

For more information, see: https://cloud.google.com/endpoints/docs/grpc/transcoding

1. Install the gRPC plugins (ensure Golang's bin folder is in PATH):
```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
```
2. Install [buf](https://github.com/bufbuild/buf)
3. Copy `buf.gen.yaml` and `buf.yaml` to your codebase.
4. Adjust configs for your needs. (proto files path, out directories, etc. For better understanding, check buf's [docs](https://buf.build/docs/generate/tutorial))
5. Run `buf dep update` & `buf generate`
