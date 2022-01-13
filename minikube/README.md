

curl -i http://127.0.0.1:8080

grpc_cli ls 127.0.0.1:8082
grpc_cli ls localhost:8082 addservice.AddService -l
grpc_cli call localhost:8082 addservice.AddService.Add "a: 5, b: 9"

minikube

go build -ldflags "-X minikube/version.Version=1.0.0" -o ./bin/app main.go

make build_container

docker run --rm --name go-web -p 8080:8080 -p 8082:8082 docker.io/ppdraga/simple_go_app
docker push docker.io/ppdraga/simple_go_app
