# Initial stage: download modules
FROM golang:1.17 as modules

ADD go.mod go.sum /m/
RUN cd /m && go mod download

# Intermediate stage: Build the binary
FROM golang:1.17 as builder

COPY --from=modules /go/pkg /go/pkg
RUN mkdir -p /app
ADD . /app
WORKDIR /app

# Build the binary with go build
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ./bin/router cmd/router/main.go
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ./bin/acl cmd/acl/main.go

# Final stage: Run the binary
FROM scratch

# and finally the binary
COPY --from=builder /app/bin/router /bin/router
COPY --from=builder /app/bin/acl /bin/acl

CMD ["/bin/router"]