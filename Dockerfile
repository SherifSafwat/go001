FROM golang:1.24-alpine
# FROM golang:1.24-bookworm AS builder

# ENV GOPROXY=direct
# ENV GONOSUMDB=*

# RUN apk add --no-cache ca-certificates openssl

# RUN go env
# RUN go env -w GOPROXY=https://goproxy.cn,direct
# RUN go env -w GOSUMDB=off
# RUN go env -w GO111MODULE=on

# # RUN --rm -it golang:1.22-bookworm bash
# RUN apt-get update && apt-get install -y openssl ca-certificates curl --no-install-recommends
# # Then
# RUN echo | openssl s_client -showcerts -servername proxy.golang.org -connect proxy.golang.org:443 2>/dev/null | sed -ne '/-BEGIN CERTIFICATE-/,/-END CERTIFICATE-/p'


# RUN apk add --no-cache ca-certificates
# RUN apk add --no-cache ca-certificates && update-ca-certificates
# RUN echo | openssl s_client -connect proxy.golang.org:443 -showcerts 2>/dev/null
# ENV GOPROXY=direct
# ENV GOSUMDB=off
# ENV GOINSECURE=proxy.golang.org
# ENV GOPROXY=https://proxy.golang.org

# RUN apt-get update && apt-get install -y ca-certificates
# RUN update-ca-certificates

# RUN apk add --no-cache ca-certificates

# RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /go-app main.go

CMD [ "/go-app" ]
