# builder image
FROM golang:1.11-alpine3.8 as builder

ENV CGO_ENABLED 0
RUN apk --no-cache add git
RUN go get github.com/golang/dep/cmd/dep
WORKDIR /go/src/github.com/stijndehaes/kube-conformity
COPY . .
RUN dep ensure -vendor-only
RUN go test -v ./...
ENV GOARCH amd64
RUN go build -o /bin/kube-conformity -v \
  -ldflags "-X main.version=$(git describe --tags --always --dirty) -w -s"

# final image
FROM alpine:3.6
MAINTAINER Stijn De Haes <stijndehaes@gmail.com>

RUN apk --no-cache add openssl ca-certificates dumb-init
COPY mailtemplate.html /etc/kube-conformity/mailtemplate.html
COPY config.yaml /etc/kube-conformity/config.yaml
COPY --from=builder /bin/kube-conformity /etc/kube-conformity/kube-conformity

USER 65534
ENTRYPOINT ["dumb-init", "--", "/etc/kube-conformity/kube-conformity"]
