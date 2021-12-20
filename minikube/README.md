

curl -i http://127.0.0.1:8080

grpc_cli ls 127.0.0.1:8082
grpc_cli ls localhost:8082 addservice.AddService -l
grpc_cli call localhost:8082 addservice.AddService.Add "a: 5, b: 9"


