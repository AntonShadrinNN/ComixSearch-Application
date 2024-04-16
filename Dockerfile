FROM golang:1.22.0-alpine3.19 as build

RUN apk add --no-cache ca-certificates

WORKDIR /src

COPY go.mod ./
COPY go.sum ./

RUN go mod tidy

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main cmd/main.go

FROM scratch

WORKDIR /

COPY --from=build /src/main .
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
EXPOSE 8080

ENTRYPOINT [ "./main" ]
