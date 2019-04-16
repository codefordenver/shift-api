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
	sam local start-api --env-vars config.json 

.PHONY: dev
dev:
	reflex -c reflex.conf

.PHONY: test
test:
	go test -v ./api

.PHONY: package
package:
	sam package \
        --template-file template.yaml \
        --output-template-file packaged.yaml \
        --s3-bucket codefordenver

.PHONY: deploy
deploy:
	sam deploy \
        --template-file packaged.yaml \
        --stack-name shift-api-serverless-app-stack \
        --capabilities CAPABILITY_IAM \
        --region us-east-2

validate-circleci:
	circleci config validate