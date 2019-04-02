# builder image
FROM golang:1.12-alpine3.9 as builder

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
RUN mkdir /tmp/result/ && cp /bin/kube-conformity /tmp/result/kube-conformity
COPY mailtemplate.html /tmp/result/mailtemplate.html
COPY config.yaml /tmp/result/config.yaml

# final image
FROM gcr.io/distroless/base
MAINTAINER Stijn De Haes <stijndehaes@gmail.com>
COPY --from=builder /tmp/result/ /
CMD ["/kube-conformity"]
