# News API

A microservice to store and get news articles.

## Technologies

There are HTTP server written in Go and PostgreSQL storage.

## Development

> Prerequisites: `docker`, `docker-compose`, `go@1.18`, `make` must be installed.

### Binary

Build the binary `bin/news-api`:

```shell
make build
```

Build docker image `news-api:latest`:

```shell
make docker
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

### Tests

Run tests:

```shell
make test
```

### CI

There are configured GitHub actions for build, lint, and run tests.
See [.github/workflows](.github/workflows) directory.