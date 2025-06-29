# KantSearchBackend

This is the Go backend for the kant-search project. It inserts and reads data to/from an Elasticsearch database using the [Elastic Go Client](https://www.elastic.co/docs/reference/elasticsearch/clients/go).

## Contributing

If you want to improve this codebase or add a feature, feel free to open a pull request.

## Installation

### Using Docker

- pull the newest container with `docker pull ghcr.io/frhorschig/kant-search-backend`
- run the docker container with the appropriate environment variables:

```bash
docker run -d \
  -v /path/to/local/ssl/files:/ssl \
  -v /path/to/config/directory:/config \
  -e KSGO_CERT_PATH='/ssl/<my-cert-name>.pem' \
  -e KSGO_KEY_PATH='/ssl/<my-key-name>.pem' \
  -e KSGO_CONFIG_PATH='/config' \
  # ... add remaining environment variables
  -p 3000:3000
  --name ks-go \
  frhorschig/kant-search-backend
```

### Using the binary

- download the binary and the configuration file from a release
- set the environment variables and run the binary
- it is recommended to forward all stdout and stderr output to adequate log file location

### Environment variables

These environment variables are necessary for the application to function properly:
- `KSGO_URL` - the host URL of the database
- `KSDB_PORT` - the port of the database
- `KSDB_USER` - the name of the elasticsearch user
- `KSDB_PWD` - the password of the elasticsearch user
- `KSDB_CERT_HASH` - the hash of the certificate of the elasticsearch application
- `KSGO_ALLOW_ORIGINS` - comma separated list of URLs allowed to communicate with the backend (use `*` to allow all)
- `KSGO_CONFIG_PATH` - path to the configuration directory
- `KSGO_DISABLE_SSL` - set to "true" to disable SSL
- `KSGO_CERT_PATH` - path to the SSL certificate
- `KSGO_KEY_PATH` - path to the SSL key

### Configuration

The configuration file `volume-metadata.json` contains metadata of the volumes and works of the Akademie-Ausgabe that is missing from or incomplete in the Akademie-Ausgabe texts, e.g. the Siglum or the publication year of some works. The application expects to find the `volume-metadata.json` file in the `KSGO_CONFIG_PATH` directory.

## Development setup

Refer to the [parent project](https://github.com/FrHorschig/kant-search) for a general overview and scripts for helping with the development setup, including a script to start the backend locally together with the database and the frontend.
