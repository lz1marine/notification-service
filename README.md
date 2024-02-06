# notification-service
This is a notification service

You can find our high-level overview documents in our [documents directory](/documents/README.md).

### Prerequisites

1. make
2. [go](https://go.dev/doc/install) 1.21+
3. [docker](https://docs.docker.com/engine/install/) latest
4. [swagger](github.com/swaggo/swag/cmd/swag) v1.16.3+
5. [golangci-lint](github.com/golangci/golangci-lint/cmd/golangci-lint) 
6. [gofumpt](mvdan.cc/gofumpt@latest)
7. [goimports](golang.org/x/tools/cmd/goimports@latest)
8. [mysql](https://dev.mysql.com/downloads/installer/)
9. [redis](https://redis.io/docs/install/install-redis/)

### Preparing the system

In the [Makefile](Makefile) there are recipes to more easily handle the project:

Make sure you have installed everything go needs
```bash
make install-requirements
```

Fix or init your vendor dir
```bash
make deps
```

To format your code
```bash
make fmt
```

To then vet it
```bash
make vet
```

To then run linters
```bash
make lint
```

To then create the docs
```bash
make docs
```

You can manually validate the docs by starting a swagger service. Navigate to http://localhost and you are good to go!
```bash
docker pull swaggerapi/swagger-ui
docker run -p 80:8080 -e SWAGGER_JSON=$(pwd)/openapi/swagger.json -v $(pwd)/openapi:$(pwd)/openapi swaggerapi/swagger-ui
```

To then build it, note you can build a specific image, i.e. `build-notification-server` `build-notification-worker`
```bash
make build-all
```

To then test it
```bash
make tests
```

You can also run the entire ci flow
```bash
make ci
```

### Running the system

Before running the system make sure you have it [configured first](#configuration). You can run the system either locally or in a container.

1. Make sure you have mysql and redis running.
```bash
# Add databases and tables
mysql
CREATE DATABASE IF NOT EXISTS `users`;
CREATE DATABASE IF NOT EXISTS `notifications`;
exit
mysql users < scripts/mysql/01-init-users.sql
mysql notifications < scripts/mysql/01-init-notifications.sql
```

Install the redis stack
```bash 
curl -fsSL https://packages.redis.io/gpg | sudo gpg --dearmor -o /usr/share/keyrings/redis-archive-keyring.gpg
sudo chmod 644 /usr/share/keyrings/redis-archive-keyring.gpg
echo "deb [signed-by=/usr/share/keyrings/redis-archive-keyring.gpg] https://packages.redis.io/deb $(lsb_release -cs) main" | sudo tee /etc/apt/sources.list.d/redis.list
sudo apt-get update
sudo apt-get install redis-stack-server
```

2. Run the notification server:
```bash
# Run locally
go run cmd/notification-service/main.go
# Or in a container
VERSION=server-v0.0.7
docker network create my_network
docker run -p 12345:8080 --network=my_network -e QUEUE_ENDPOINT=host.docker.internal:6379 -e DB_LOCATION=host.docker.internal -e DB_PORT=3306 -v /app/secrets:/app/secrets lz1marine/notification-service:${VERSION}
```

3. Run the notification worker:
```bash
export TYPE=email # or sms
# Run locally
go run cmd/notification-worker/main.go
# Or in a container
VERSION=worker-v0.0.3
docker network create my_network
docker run --network=my_network -e QUEUE_ENDPOINT=host.docker.internal:6379 -e DB_LOCATION=host.docker.internal -e DB_PORT=3306 -v /app/secrets:/app/secrets lz1marine/notification-service:${VERSION}
```

### Examples

After having [run the system](#running-the-system), you can now play around with it. Here are some examples (note, change the port to 12345 if running the [example containers](#running-the-system)):

* Add a template in redis
```bash
email_template='{
    "template_id": "1",
    "template": "<!DOCTYPE html>\n<html>\n<body>\n    <h3>Name:</h3><span>Hello {{.Name}}</span><br/><br/>\n    <h3>Email:</h3><span>{{.Email}}</span><br/>\n    <h3>Message:</h3><span>{{.Message}}</span><br/>\n</body>\n</html>",
    "is_enabled": true
}'
redis-cli -n 10 SET "1" "$email_template"
```

* Get all of the available notifications and whether they are enabled or not
```bash
curl http://localhost:8080/api/v1/notifications --request "GET"
```

* Get all notifications channels that a specific user has subscribed for

```bash
curl http://localhost:8080/api/v1/notifications/sub/1 --request "GET" 
```

* Subscribe user to specific notification channels

```bash
curl http://localhost:8080/api/v1/notifications/sub/1 \
--include \
--header "Content-Type: application/json" \
--request "PATCH" \
--data '{"channels":[{"name":"sms","is_enabled":false},{"name":"slack","is_enabled":false}]}'
```

* Subscribe user to specific notification channels

```bash
curl http://localhost:8080/api/v1/notifications/sub/1 \
--include \
--header "Content-Type: application/json" \
--request "PATCH" \
--data '{"channels":[{"name":"sms","is_enabled":false},{"name":"slack","is_enabled":false}]}'
```

* Send an email
```bash
curl http://localhost:8080/api/v1/internal/notifications/123456789 \
--include \
--header "Content-Type: application/json" \
--request "POST" \
--data '{"channel":"email","subject":"A test email","message":"This is a test message","topic_id":"1","template_id":"1"}'
```

* Send an sms
```bash
curl http://localhost:8080/api/v1/internal/notifications/987654321 \
--include \
--header "Content-Type: application/json" \
--request "POST" \
--data '{"channel":"sms","message":"This is a test sms message","topic_id":"1"}'
```


**Note**: To enable sending a message if using two-factor, go to your google account in https://myaccount.google.com/. Then, search for "App password" and add yourself an app password, then copy it and use it (trim it first).

### Configuration

View the notification server [configuration](cmd/notification-service/README.md#configuration).

View the notification worker [configuration](cmd/notification-worker/README.md#configuration).

### Next steps

Before going to production we have to make sure the following are also completed:
* We should add kubernetes helm charts to be able to deploy the system
* We should add tests to cover at least a percentage of our code/85%?
* We should add Slack channel
* We should handle the state of the message in the db
* We should change our fmt.Print statements with an actual logger
* We should add a Jenkins pipeline
