.PHONY: deps clean build

.PHONY: deps
deps:
	go get -u ./...

.PHONY: clean
clean:
	rm -rf ./api/api

.PHONY: build
build:
	GOOS=linux GOARCH=amd64 go build -o api/api ./api

.PHONY: local
local:
	sam local start-api

.PHONY: dev
dev:
	reflex -c reflex.conf

.PHONY: test
test:
	go test -v ./api

.PHONY: deploy
deploy:
	sam deploy \
        --template-file packaged.yaml \
        --stack-name shift-api-serverless-app-stack \
        --capabilities CAPABILITY_IAM \
        --region us-west-2