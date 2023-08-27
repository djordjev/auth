# auth

Simple username/password authentication service

[![test](https://github.com/djordjev/auth/actions/workflows/test.yaml/badge.svg)](https://github.com/djordjev/auth/actions/workflows/test.yaml)

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
