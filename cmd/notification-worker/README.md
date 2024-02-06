# Notification Worker

### Configuration

TODO: make this pretty

You should have the following files in `/app/secrets`. If any files are missing, an error will occur.
* `queue_password` - a single line containing the password for the distributed queue the worker will be using.

You can set the following env vars.
* `TYPE` - the type of channel we are starting. Either `email`, `sms`, `slack`. Mandatory.
* `MAX_CONNECTIONS` - max goroutines that can run in parallel. Defaults to `200`.
* `QUEUE_ENDPOINT` - the endpoint of the distributed queue. Defaults to `localhost:6379`

If setting an email channel, please add the following 2 files:
* `email_username` - a single line containing username that will be user for the third-party software.
* `email_password` - a single line containing the password.

If setting an email channel, please add the following 2 env vars:
* `HOST` - the third-party email server. Defaults to `smtp.gmail.com`.
* `PORT` - the third-party email server port. Defaults to `587`.

If setting an sms channel, please add the following 3 files:
* `sms_username` - a single line containing username that will be user for the third-party software.
* `sms_password` - a single line containing the password.
* `sms_phone_sender` - a single line containing the phone number.
