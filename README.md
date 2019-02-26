# SHIFT_API

## Development

#### Starts local server API
- `make local`

#### Fetches dependencies and starts file watcher (hot-reload) 
- `make dev`

## Prerequisites

* [AWS CLI](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-install.html) configured with Administrator permission
* [AWS SAM CLI](https://aws.amazon.com/serverless/sam/)
* [Docker installed](https://www.docker.com/community-edition)
* [Golang](https://golang.org)
* [Reflex](https://github.com/cespare/reflex)

Add to your environment variables:
```
# GOPATH
export GOPATH="${HOME}/go"
export PATH="$GOPATH/bin:$PATH"
export GO111MODULE=on
```