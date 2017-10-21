# builder image
FROM golang:1.9-alpine

RUN apk -U add git
WORKDIR /go/src/github.com/stijndehaes/kube-conformity
COPY . .
RUN go test -v ./...
RUN go build -o /bin/kube-conformity -v \
  -ldflags "-X main.version=$(git describe --tags --always --dirty) -w -s"

# final image
FROM alpine:3.6
MAINTAINER Stijn De Haes <stijndehaes@gmail.com>

RUN addgroup -S kube-conformity && adduser -S -g kube-conformity kube-conformity
COPY --from=0 /bin/kube-conformity /bin/kube-conformity

USER kube-conformity
ENTRYPOINT ["/bin/kube-conformity"]
