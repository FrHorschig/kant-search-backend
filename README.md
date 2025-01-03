# KantSearchBackend

This is the Go backend for the kant-search project. It reads frontend data from a PostgreSQL database.

This backend also implements the endpoint for uploading data and writing it to the PostgreSQL database. This code will be separated more clearly from the read and search code in the future.

## Contributing

If you want to improve this codebase or add a feature, feel free to open a pull request. Make sure to explain any deviation from existing code conventions.

## Installation

### Using Docker

- pull the newest container with `docker pull ghcr.io/frhorschig/kant-search-backend`
- run the docker container with the appropriate environment variables (you can ignore the python environment variables):

```bash
docker run -d \
  -v /path/to/local/ssl/files:/ssl \
  -e KSGO_CERT_PATH='/ssl/<my-cert-name>.pem' \
  -e KSGO_KEY_PATH='/ssl/<my-key-name>.pem' \
  # ... more environment variables
  -p 3000:3000
  --name ks-go \
  frhorschig/kant-search-backend
```

### Using the binary

- ensure that Python and the `spacy` library are installed on the system (see [here](https://spacy.io/usage) for instructions)
- download the binary and the python script from a release
- set the appropriate environment variables and run the binary

### Environment variables

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
- `KSGO_PYTHON_SCRIPT_PATH` - path to the python scripts directory

## Development setup

Refer to the [parent project](https://github.com/FrHorschig/kant-search) for a general overview and scripts for helping with the development setup, including a script to start the backend locally together with the database and the frontend.
