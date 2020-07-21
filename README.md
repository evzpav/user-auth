# User Auth

It is a server side web app using Go templates that allows user to signup, login, reset password, see and edit profile.

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
#Run mongo locally
make run-mysql

#Run server locally or
make run-local

#Run on docker
make run-docker
```

## Deployment

### Build

```bash
make build
```

## Contributing

See [CONTRIBUTING](CONTRIBUTING.md) documentation for more details.
