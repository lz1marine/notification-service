# Build the binary
FROM golang:1.21-alpine as build

WORKDIR /
COPY ./api ./api
COPY ./pkg/api ./pkg/api
COPY ./go.mod .
COPY ./go.sum .


COPY ./cmd/notification-gc ./cmd/

# TODO: cherry-pick specific files
COPY ./pkg ./pkg/

RUN go mod download
RUN GO111MODULE=on go build -o notification-gc cmd/main.go

FROM alpine:3.18
WORKDIR /
RUN apk upgrade --update --no-cache && apk add curl

ARG USERNAME=admin
RUN addgroup -S $USERNAME \
    && adduser -S $USERNAME -G $USERNAME
USER $USERNAME

COPY --from=build /notification-gc /app/notification-gc

ENTRYPOINT ["/app/notification-gc"]