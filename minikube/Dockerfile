# Init args
ARG GIT_COMMIT
ARG VERSION
ARG PROJECT

# Initial stage: download modules
FROM golang:1.17 as modules

ADD go.mod go.sum /m/
RUN cd /m && go mod download

# Intermediate stage: Build the binary
FROM golang:1.17 as builder

ARG GIT_COMMIT
ENV GIT_COMMIT=$GIT_COMMIT
ARG VERSION
ENV VERSION=$VERSION
ARG PROJECT
ENV PROJECT=$PROJECT

ENV GOSUMDB=off
ENV GO111MODULE=on

COPY --from=modules /go/pkg /go/pkg

# add a non-privileged user
RUN useradd -u 10001 app_user

RUN mkdir -p /app
ADD . /app
WORKDIR /app

# Build the binary with go build
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
    go build -ldflags "-X ${PROJECT}/version.Version=1.0.0 \
    	-X ${PROJECT}/version.Commit=${GIT_COMMIT} " -o ./bin/app main.go

# Final stage: Run the binary
FROM scratch

EXPOSE 8080
EXPOSE 8082

# don't forget /etc/passwd from previous stage
COPY --from=builder /etc/passwd /etc/passwd
USER app_user

# and finally the binary
COPY --from=builder /app/bin/app /app

CMD ["/app"]
