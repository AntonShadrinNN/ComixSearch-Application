##
## STEP 1 - BUILD
##
ARG REGISTRY=docker.io
FROM ${REGISTRY}/golang:alpine3.19 AS build-stage

WORKDIR /app

RUN apk --no-cache add ca-certificates=20240226-r0

COPY go.mod go.sum ./

RUN go mod tidy

COPY config.yml ./
COPY . ./

RUN adduser\
    --disabled-password\
    --home "/nonexistent"\
    --shell "/sbin/nologin"\
    --no-create-home\
    nonroot\
    && CGO_ENABLED=0 GOOS=linux go build -a -o main cmd/main.go

##
## STEP 2 - DEPLOY
##

FROM scratch

LABEL org.label-schema.vcs-url="https://github.com/AntonShadrinNN/ComixSearch-Application"

WORKDIR /

COPY --from=build-stage /etc/passwd /etc/passwd
COPY --from=build-stage /etc/group /etc/group
COPY --from=build-stage /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build-stage /app/config.yml ./
COPY --from=build-stage /app/main ./

EXPOSE 8000

USER nonroot:nonroot

ENTRYPOINT ["./main"]