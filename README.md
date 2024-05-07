# KantSearchBackend

This is the Go backend for the kant-search project.

## Environment variables

- `KSDB_USER` - the name of the database user
- `KSDB_PASSWORD` - the password of the database user
- `KSDB_NAME` - the name of the database
- `KSGO_DB_HOST` - the host of the database
- `KSDB_PORT` - the port of the database
- `KSGO_DB_SSLMODE` - SSL mode of the database (see [here](https://www.postgresql.org/docs/current/libpq-ssl.html) for valid values)
- `KSGO_ALLOW_ORIGINS` - comma separted list of URLs allowed to communicate with the backend (use `*` to allow all)
- `KSGO_DISABLE_SSL` - set to "true" to disable SSL
- `KSGO_CERT_PATH` - path to the SSL certificate
- `KSGO_KEY_PATH` - path to the SSL key
- `KSGO_PYTHON_BIN_PATH` - path to the python executable
- `KSGO_PYTHON_SCRIPT_PATH` - path to the `split_text.py` python script
