.PHONY: generate import api migration_create init


init:
	go mod tidy 
	make generate
	make fmt 
	go mod tidy 

serve:
	go run . serve
    
migrate:
	go run . migrate --cmd=up

build: build-frontend
	go build -o bin/

generate: api
	go generate ./...
	make import

import:
	goimports -w .

api:
	oapi-codegen --config=api/config.yaml api/api.yaml

fmt:
	go fmt ./...

 migration_create:
	goose create $(name) go --dir=db/postgres/migrations