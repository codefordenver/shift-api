# SHIFT_API

## Development

#### Make sure Docker daemon is running

#### Start local server API
- `make local`

#### Fetch dependencies and starts file watcher (hot-reload) 
- `make dev`

## Prerequisites

#### Install:

* [AWS CLI](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-install.html) configured with Administrator permission
* [AWS SAM CLI](https://aws.amazon.com/serverless/sam/)
* [Docker installed](https://www.docker.com/community-edition)
* [Golang](https://golang.org)
* [Reflex](https://github.com/cespare/reflex)

#### Add to your environment variables:
```
# GOPATH
export GOPATH="${HOME}/go"
export PATH="$GOPATH/bin:$PATH"
export GO111MODULE=on
```

#### Add a file to store your environment variables
You need a `config.json` file.  You can copy the `config.example.json` file (`mv config.example.json config.json`) and add environment variables as we need more in the app. 

#### Optional:
- [Docker desktop](https://www.docker.com/products/docker-desktop) for running Docker daemon