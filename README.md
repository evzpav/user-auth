# Documents

Documents is a REST API which main goal is an example of a CRUD written in Golang with some best practices and domain driven design.

## Architecture Diagrams

To make the project easier to understand, some other diagrams have been made available in the following directory [`docs/diagrams`](docs/diagrams).

### Generate Diagrams

In order to facilitate the documentation process, a command is available to generate `.png` images from `.puml` files.

```bash
make create-diagrams
```

Files are generated and available inside `docs/diagrams`.

_Diagrams are defined using a simple and intuitive language [PlantUML](http://plantuml.com/)._


## REST API Documentation

All available REST API documentation exposed by the project was documented using the [OpenAPI](https://www.openapis.org/) standard.

To view this documentation locally use the following command:
```bash
make run-swagger
```
The same documentation is available online and can be accessed by clicking [here](https://gitlab.com/evzpav/user-auth/api/index.html).


_API routes documentation generate by [Swagger-UI](https://github.com/swagger-api/swagger-ui)._

## Getting Started

### Prerequisites

- [Golang](http://golang.org/)(>11.0)
- [GNU Make](https://www.gnu.org/software/make/)
- [Docker](http://docker.com)

### Environment variables

```bash
DOCUMENTS_HOST=localhost
DOCUMENTS_PORT=8080
LOGGER_LEVEL=error #values: error; warn; info; debug;
MONGO_URL=localhost:27017

```

### Installing and running locally

```bash
#Install dependencies
make install

#Run mongo locally
make env

#Run server locally
make run

#Run server locally with custom host and port
MONGO_URL=localhost:27017 \
DOCUMENTS_HOST=localhost \
DOCUMENTS_PORT=5002 \
LOGGER_LEVEL=error \
make run-local
```

## Setting up git hooks

We have a hook to help with the development process to keep standards up.
To set them up just run:

```bash
make git-config
```

## Running acceptance tests

To run acceptance tests locally use the following commands:

```bash
make acceptance-tests
```

## Running the tests and coverage report

To view report of tests locally use the following command:

```bash
make test
```

The same report is available online and can be accessed by clicking [here](https://gitlab.com/evzpav/user-auth/coverage/index.html).

_Coverage report generate by [Gocov](https://github.com/axw/gocov)._

## Running the lint verification

```bash
make lint
```
_Lint report generate by [GolangCI-lint](https://github.com/golangci/golangci-lint)._

## Deployment

### Build

```bash
make build
```

### Create release image, add tag and push

```bash
make image tag push
```

### Run registry image locally

```bash
make run-docker

make remove-docker
```

## Inspiration

### Package organization

The package structure used in this project was inspired by the [golang-standards](https://github.com/golang-standards/project-layout) project.

### Project layers organization

The project layers structure used in this project was inspired by the **Hexagonal Architecture** (Ports & Adapters).

## Contributing

See [CONTRIBUTING](CONTRIBUTING.md) documentation for more details.

## Changelog

See [CHANGELOG](CHANGELOG.md) documentation for more details.
