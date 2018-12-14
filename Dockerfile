FROM golang:alpine
WORKDIR /app
RUN apk add git \
    && go get github.com/digitalocean/godo \
    && go get golang.org/x/oauth2
COPY main.go .
RUN go build -o kubeget main.go

FROM alpine:latest
WORKDIR /work
COPY --from=0 /app/kubeget /usr/bin/kubeget
RUN apk --no-cache add curl ca-certificates \
    && curl -LO https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl \
    && chmod +x ./kubectl \
    && mv ./kubectl /usr/bin/kubectl \
    && apk del curl
