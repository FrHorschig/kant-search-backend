# KantSearchBackend

This is the Go backend for the kant-search project. It reads and writes frontend data to and from a PostgreSQL database. The Go application uses the `spacy` library in a python script for splitting texts into sentences.

## Installation

### Using Docker

- pull the newest container with `docker pull ghcr.io/frhorschig/kant-search-backend`
- run the docker container with the appropriate environment variables (you can ignore the python environment variables):

```bash
docker run -d \
  -v /path/to/local/ssl/files:/ssl \
  -e KSGO_CERT_PATH='/ssl/cert-file.pem' \
  -e KSGO_KEY_PATH='/ssl/key-file.pem' \
  # ... more environment variables
  -p 3000:3000
  --name ks-backend \
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
- `KSGO_PYTHON_SCRIPT_PATH` - path to the `split_text.py` python script

## Development setup

Refer to the [parent project](https://github.com/FrHorschig/kant-search) for information about the development setup.
