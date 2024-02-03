# notification-service
This is a notification service

To run the API:
```bash
go run cmd/notification-service/main.go
```

To test the API:
1. Get all of the available notifications and whether they are enabled or not:

```bash
curl http://localhost:8080/v1/notifications --request "GET"
```

2. Get all notifications channels that a specific user has subscribed for:

```bash
curl http://localhost:8080/v1/notifications/sub/1 --request "GET" 
```

3. Subscribe user to specific notification channels:

```bash
curl http://localhost:8080/v1/notifications/sub/1 \
--include \
--header "Content-Type: application/json" \
--request "PATCH" \
--data '{"channels":[{"name":"SMS","is_enabled":false},{"name":"Slack","is_enabled":false}]}'
```

To run the worker API:

**Note**: To enable sending a message if using two-factor, go to your google account in https://myaccount.google.com/. Then, search for "App password" and add yourself an app password, then copy it and use it (trim it first).

```bash
go run cmd/notification-worker/main.go
```

1. Send an email
```bash
curl http://localhost:8080/v1/temp/notifications/1 \
--include \
--header "Content-Type: application/json" \
--request "POST" \
--data '{"title":"A test email my friend","message":"This is a test message my friend","topic_id":"1","template_id":"1"}'
```
