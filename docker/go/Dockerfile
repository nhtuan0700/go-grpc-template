FROM golang:1.22.4

WORKDIR /app

RUN apt-get update && apt-get install -y \
  make \
  protobuf-compiler \
  git

RUN go install github.com/bufbuild/buf/cmd/buf@v1.35.0

RUN go install github.com/google/wire/cmd/wire@latest
RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.59.0

ARG user=appuser
ARG group=appuser
ARG uid=1000
ARG gid=1000
RUN groupadd -g ${gid} ${group}
RUN useradd -u ${uid} -g ${group} -s /bin/sh -m ${user} 
# USER ${uid}:${gid}
# RUN CGO_ENABLED=0 go install -ldflags "-s -w -extldflags '-static'" github.com/go-delve/delve/cmd/dlv@latest

CMD ["bash"]
