package compose

import (
	"fmt"
	"os"

	"github.com/AugustineAurelius/eos/pkg/strings"
	"gopkg.in/yaml.v3"
)

func addApplication(m map[string]any) {
	var app map[string]any
	s := `
build:
  context: . 
  dockerfile: Dockerfile
ports:
  - "${APP_PORT}:8080"
restart: unless-stopped`

	if len(m) != 0 {
		s += "\ndepends_on: \n"

	}

	for k := range m {
		s += `  - ` + k + "\n"

	}

	yaml.Unmarshal([]byte(s), &app)

	m["app"] = app
}

func addAppDockerfile() {
	f, err := os.Create("Dockerfile")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}
	defer f.Close()

	b := strings.Builder{}

	b.WriteStringWithEnter("FROM golang:1.22-alpine AS builder").
		WriteStringWithEnter("RUN apk add --no-cache --update go").
		WriteStringWithEnter("WORKDIR /app").
		WriteStringWithEnter("COPY go.* ./").
		WriteStringWithEnter("RUN go mod download").
		WriteStringWithEnter("COPY . .").
		WriteStringWithEnter(`RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o main ./cmd/app/main.go`). //TODO sprintf path/cgo/os
		WriteEnter().
		WriteStringWithEnter("FROM alpine:latest").
		WriteStringWithEnter("RUN apk --no-cache add ca-certificates").
		WriteStringWithEnter("RUN adduser -D -g '' user"). //TODO user
		WriteStringWithEnter("WORKDIR /app").
		WriteStringWithEnter("COPY --from=builder /app/main .").
		WriteStringWithEnter("RUN chown user:user /app/main").
		WriteStringWithEnter("USER user").
		WriteStringWithEnter(`CMD ["./main"]`)

	f.Write(b.Bytes())
}
