# User Auth

It is a server side web app using Go templates that allows user to signup, login, reset password, see and edit profile.
Persistance using MySql.

## Getting Started

### Prerequisites

- [Golang](http://golang.org/)(>11.0)
- [GNU Make](https://www.gnu.org/software/make/)
- [Docker](http://docker.com)

### Environment variables

```bash
	SESSION_KEY
	HOST
	PORT
	LOGGER_LEVEL
	EMAIL_FROM
	EMAIL_PASSWORD
	GOOGLE_KEY
	GOOGLE_SECRET
	GOOGLE_MAPS_API_KEY
	PLATFORM_URL
	DATABASE_URL
```

### Installing and running locally

```bash
#Run mongo locally
make run-mysql

#then 

#Run server locally 
make run-local

# or

#Run on docker
make run-docker
```

## Contributing

See [CONTRIBUTING](CONTRIBUTING.md) documentation for more details.
