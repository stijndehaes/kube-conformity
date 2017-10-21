# builder image
FROM golang:1.9-alpine
MAINTAINER Stijn De Haes <stijndehaes@gmail.com>

RUN apk -U add git
WORKDIR /go/src/github.com/stijndehaes/kube-conformity
COPY . .
RUN go test -v ./...
RUN go build -o /bin/kube-conformity -v \
  -ldflags "-X main.version=$(git describe --tags --always --dirty) -w -s"

RUN addgroup -S kube-conformity && adduser -S -g kube-conformity kube-conformity
USER kube-conformity
ENTRYPOINT ["/bin/kube-conformity"]