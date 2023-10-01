# auth

Simple username/password authentication service

[![test](https://github.com/djordjev/auth/actions/workflows/test.yaml/badge.svg)](https://github.com/djordjev/auth/actions/workflows/test.yaml)

The app can be used in two ways.

1. Within other go application. Import package server and mount it on existing http mux on particular entrypoint
   (for example `/auth`). Then application becomes bound to particular endpoint
2. Run `main.go` file what will start up a new server and mount `auth` to home route `/`

In order to run properly application needs `postgresql` database running for storing users and `redis` server running
for storing sessions. In order to send emails (for forget password or verification) it needs to have Mailjet api key
provided through environment variables.

### Configuration through environment variables

Env variables

```
DB_HOST - PostgreSQL database host
DB_PASS - PostgreSQL database password
DB_USER - PostgreSQL database user
DB_NAME - PostgreSQL database name
DB_PORT - PostgreSQL port. Optional: default 5432
PORT - Port on which the app will run if not used as mount on existing app. Optional.
GO_ENV - string `development` or `production`
DOMAIN - Domain name where the app is hosted. Optional
REQUIRE_VERIFICATION - If set to `true` user will get an email with verification link. If `false` new accounts will be automatically verified. Optional: default false
VERIFICATION_LINK - If required to verify account this is a base of link to verify account that user will get in email. Required only if `REQUIRE_VERIFICATION` is true.
FORGET_PASSWORD_LINK - Base of the link that user will get in email to reset forgotten password.
SENDER - email address that will be used as `sender` of emails.
MAILJET_API_KEY - Mailjet api key
MAILJET_SECRET_KEY - Mailjet secret key
REDIS_DB - Redis database number. Optional: default 0
REDIS_HOST - Redis host
REDIS_PASSWORD - Redis password. Optional
REDIS_PORT - Redis port. Optional: default 6379
SESSION_COOKIE - Name of the cookie that will be set on login. Optional: default `_tkn`
```

## Setup

### Tools

Please install following tools

#### go-task

[Documentation](https://taskfile.dev/installation/)

```
brew install go-task/tap/go-task
```

#### Docker

Install [docker](https://docs.docker.com/get-docker/)

#### Vektra/Mockery

[GitHub repo](https://github.com/vektra/mockery)

```
brew install mockery
```

Note: see [installation page](https://github.com/vektra/mockery/wiki/Installation-Methods#go-install) on github.

#### gotestfmt

[Github repo](https://github.com/GoTestTools/gotestfmt#github-actions)

```
go install github.com/gotesttools/gotestfmt/v2/cmd/gotestfmt@latest
```

Note: Make sure to have `$GOPATH/bin` folder in path in order to be able to execute `gotestfmt` command.

#### Mailjet

This project uses [Mailjet](https://www.mailjet.com/) to send emails. In order to be able to send verification / forget passwords mails you'll
have to set up account on Mailjet. See documentation [here](https://documentation.mailjet.com/hc/en-us)

#### pg-mig

Download latest release executable from [pg-mig Releases](https://github.com/djordjev/pg-mig) page. Copy it to a folder
within `$PATH` so it's accessible from any directory.
