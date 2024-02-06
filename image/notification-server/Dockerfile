# Build the binary
FROM golang:1.21-alpine as build

WORKDIR /
COPY ./api ./api
COPY ./pkg/api ./pkg/api
COPY ./go.mod .
COPY ./go.sum .


COPY ./cmd/notification-service ./cmd/

# TODO: cherry-pick specific files
COPY ./pkg ./pkg/

RUN go mod download
RUN GO111MODULE=on go build -o notification-server cmd/main.go

FROM alpine:3.18
WORKDIR /
RUN apk upgrade --update --no-cache && apk add curl

ARG USERNAME=admin
ARG PORT=8080
RUN addgroup -S $USERNAME \
    && adduser -S $USERNAME -G $USERNAME
USER $USERNAME

COPY --from=build /notification-server /app/notification-server

EXPOSE $PORT

ENTRYPOINT ["/app/notification-server"]