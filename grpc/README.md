

1. Install protocol buffer compiler:
   sudo apt install -y protobuf-compiler

2. Install the protocol compiler plugins for Go using the following commands:
   go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.26
   go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1

3. Update your PATH so that the protoc compiler can find the plugins:
   export PATH="$PATH:$(go env GOPATH)/bin"

4. Compile module:
   protoc --go_out=. --go-grpc_out=. addressbook.proto


link: https://github.com/tensor-programming/grpc_tutorial



Тестирование grpc без клиента:
grpc_cli ls localhost:4040
grpc_cli type localhost:4040 addservice.Request
grpc_cli call localhost:4040 Add "a: 5,b: 4"
