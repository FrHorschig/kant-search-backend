# KantSearchBackend

This is the Go backend for the kant-search project. It inserts and reads data to/from an Elasticsearch database using the [Elastic Go Client](https://www.elastic.co/docs/reference/elasticsearch/clients/go).

## Contributing

If you want to improve this codebase or add a feature, feel free to open a pull request.

## Installation

Please refer to the [parent project](https://github.com/FrHorschig/kant-search).

## Configuration

The configuration file `volume-metadata.json` contains metadata of the volumes and works of the Akademie-Ausgabe that is missing from or incomplete in the Akademie-Ausgabe texts, e.g. the Siglum or the publication year of some works. The application expects to find the `volume-metadata.json` file in the `KSGO_CONFIG_PATH` directory.

### Environment variables

These environment variables are necessary for the application to function properly:
- `KSGO_URL` - the host URL of the database
- `KSDB_PORT` - the port of the database
- `KSDB_USERNAME` - the name of the elasticsearch user
- `KSDB_PASSWORD` - the password of the elasticsearch user
- `KSDB_CERT` - the path to the elasticsearch http certificate
- `KSGO_ALLOW_ORIGINS` - comma separated list of URLs allowed to communicate with the backend (use `*` to allow all)
- `KSGO_CONFIG_PATH` - path to the configuration directory
- `KSGO_CERT` - path to the SSL certificate
- `KSGO_KEY` - path to the SSL key
- `KSGO_DISABLE_SSL` - set to "true" to disable SSL (only for development!)

## Development setup

Refer to the [parent project](https://github.com/FrHorschig/kant-search) for a general overview and scripts for helping with the development setup, including a script to start the backend locally together with the database and the frontend.
