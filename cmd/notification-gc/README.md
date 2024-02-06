# Notification Garbage Collector

### Configuration

TODO: make this pretty

You should have the following files in `/app/secrets`. If any files are missing, an error will occur.
* `database_username` - a single line containing the database username.
* `database_password` - a single line containing the database password.
* `queue_password` - a single line containing the password for the distributed queue the worker will be using.

You can set the following env vars.
* `DB_LOCATION` - the database host. Defaults to `localhost`.
* `DB_PORT` - the database port. Defaults to `3306`.
* `QUEUE_ENDPOINT` - the endpoint of the distributed queue,