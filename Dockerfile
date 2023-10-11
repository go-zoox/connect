# Builder
FROM --platform=${BUILDPLATFORM:-linux/amd64} whatwewant/builder-go:v1.20-1 as builder

ARG TARGETOS

ARG TARGETARCH

WORKDIR /app

COPY go.mod ./

COPY go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -ldflags="-w -s" -v -o connect cmd/connect/main.go

# Server
# FROM  scratch # x509: certificate signed by unknown authority
FROM --platform=${TARGETPLATFORM:-linux/amd64} whatwewant/alpine:v3.17-1

LABEL MAINTAINER="Zero<tobewhatwewant@gmail.com>"

WORKDIR /app

COPY --from=builder /app/connect /bin

EXPOSE 8080

COPY ./entrypoint.sh /entrypoint.sh

CMD /entrypoint.sh
