PROJECT_PATH := $(patsubst %/,%,$(dir $(abspath $(lastword $(MAKEFILE_LIST)))))

.PHONY: run-apicast
run-apicast:
	docker run --name apicast -d --rm -p 8080:8080 -v $(PROJECT_PATH)/testdata/conf.json:/opt/app/config.json -e APICAST_BACKEND_CACHE_HANDLER=none -e THREESCALE_CONFIG_FILE=/opt/app/config.json quay.io/3scale/apicast:3scale-2.8.1-GA

.PHONY: stop-apicast
stop-apicast:
	docker stop apicast

.PHONY: build-cache
build-cache:
	go build -o apisonator-cache /Users/philipgough/go/src/github.com/philipgough/apisonator-cache/main.go

.PHONY: run-cache
run-cache:
	$(PROJECT_PATH)/apisonator-cache --upstream=https://su1.3scale.net