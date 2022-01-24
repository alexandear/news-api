# News API

A microservice to store and retrieve news posts.

## Technologies

There are HTTP server written in Go and PostgreSQL storage.

## Getting Started

> Only `docker` and `docker-compose` are needed.

Run HTTP server and database with single command:

```shell
docker-compose up
```

Server will be accessible via port **8080** by default.

## Development

> Prerequisites: `go@1.18`, `make` must be installed.

All necessary configuration stored in [docker-compose.yaml](docker-compose.yaml) and [.env](.env) files.

### Binary

Build the binary `bin/news-api`:

```shell
make build
```

Build docker image `news-api:latest`:

```shell
make docker
```

See also [Makefile](Makefile) for all available targets.

### Tests

Run tests:

```shell
make test
```

### Code style

Consistent code style enforced by `gofmt`, `EditorConfig` tools and `golangci-lint` linter.

Format code:

```shell
make format
```

Run linter:

```shell
make lint
```

### CI

There are configured GitHub actions for build, lint, and run tests.
See [.github/workflows](.github/workflows) directory.
