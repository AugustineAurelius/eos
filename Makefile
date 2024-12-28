install: 
	@go install github.com/AugustineAurelius/eos


generate:
	@go generate ./...


test: 
	@go test ./... -v -count=1



