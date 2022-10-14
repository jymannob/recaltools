ifneq (,$(wildcard ./.env))
	include .env
	export
	ENV_FILE_PARAM = --env-file .env
endif
ifneq (,$(wildcard ./.env.local))
	include .env.local
	export
	ENV_FILE_PARAM = --env-file .env.local
endif
PROJECT_VERSION ?= development
BUILD_DATE = ${shell date -u +%Y-%m-%d.%H:%M:%S}
LDFLAGS = -ldflags='-X main.buildDate=$(BUILD_DATE) -X main.buildVersion=$(PROJECT_VERSION) -X main.buildCommit=${shell git rev-parse --short HEAD}'
RASPBERRY_FLAGS = CGO_ENABLED=0 GOARCH=arm GOOS=linux

.PHONY: run test clean-build

dep: go.sum
	@go get -d

go.sum: 
	go mod tidy

run: dep
	go run $(APP_MAIN) $(filter-out $@,$(MAKECMDGOALS))

tool: dep
	go run ./cmd/recaltools/main.go $(filter-out $@,$(MAKECMDGOALS))

# PROJECT_VERSION=1.0.0 make build
build: ./bin/$(EXEC) ./bin/rpi/$(EXEC) ./bin/$(RECALTOOLS_EXEC) ./bin/rpi/$(RECALTOOLS_EXEC)

test: dep
	go test -cover $(APP_MODULE)... $(filter-out $@,$(MAKECMDGOALS))

clean-build:
	rm -rf ./bin

./bin/$(EXEC): clean-build dep
	go build $(LDFLAGS) -o ./bin/$(EXEC) $(APP_MAIN)

./bin/rpi/$(EXEC): clean-build dep
	$(RASPBERRY_FLAGS) go build $(LDFLAGS) -o ./bin/rpi/$(EXEC) $(APP_MAIN)

./bin/$(RECALTOOLS_EXEC): clean-build dep
	go build $(LDFLAGS) -o ./bin/$(RECALTOOLS_EXEC) $(RECALTOOLS_MAIN)

./bin/rpi/$(RECALTOOLS_EXEC): clean-build dep
	$(RASPBERRY_FLAGS) go build $(LDFLAGS) -o ./bin/rpi/$(RECALTOOLS_EXEC) $(RECALTOOLS_MAIN)