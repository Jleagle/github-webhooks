# Build image
FROM golang:1.10-alpine AS build-env
WORKDIR /go/src/github.com/Jleagle/github-webhooks/
COPY . /go/src/github.com/Jleagle/github-webhooks/
RUN apk update && apk add curl git openssh
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
RUN dep ensure
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo

# Runtime image
FROM alpine:3.8 AS runtime-env
WORKDIR /root/
RUN mkdir -p /root/scripts
COPY --from=build-env /go/src/github.com/Jleagle/github-webhooks/github-webhooks ./
CMD ["./github-webhooks"]
