.PHONY: deps clean build

.PHONY: deps
deps:
	go get -u ./...

.PHONY: clean
clean:
	rm -rf ./api/api
	rm -rf ./census_acs5/census_acs5
	rm -rf ./geography/geography

.PHONY: build
build:
	GOOS=linux GOARCH=amd64 go build -o api/api ./api
	GOOS=linux GOARCH=amd64 go build -o census_acs5/census_acs5 ./census_acs5
	GOOS=linux GOARCH=amd64 go build -o geography/geography ./geography

.PHONY: local
local:
	sam local start-api --env-vars config.json 

.PHONY: dev
dev:
	reflex -c reflex.conf

.PHONY: test
test:
	go test -v ./api
	go test -v ./census_acs5
	go test -v ./goegraphy

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
