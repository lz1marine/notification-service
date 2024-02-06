# Notification Worker

### Configuration

TODO: make this pretty

You should have the following files in `/app/secrets`. If any files are missing, an error will occur.
* `username` - a single line containing username that will be user for the third-party software. Use your email.
* `password` - a single line containing the password.
* `queue_password` - a single line containing the password for the distributed queue the worker will be using.

You can set the following env vars.
* `HOST` - the third-party email server. Defaults to `smtp.gmail.com`.
* `PORT` - the third-party email server port. Defaults to `587`.
* `MAX_CONNECTIONS` - max goroutines that can run in parallel. Defaults to `200`.
* `QUEUE_ENDPOINT` - the endpoint of the distributed queue.