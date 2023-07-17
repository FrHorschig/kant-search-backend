# KantSearchBackend

## Serving locally

For local development, run `source deployment/local_env.bash` to set all necessary environment variables and `go run .` in the `src` directory to start the server locally.

## Serving from a Docker container

To create a docker image, run `docker build -f ./deployment/Dockerfile -t kant-search-backend .` to build the docker image. You can start a container with this image by running `docker run -p 3000:3000 --network=kant-search kant-search-backend`.
