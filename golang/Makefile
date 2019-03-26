CURRENT_DIR := $(shell pwd)
# CURRENT_TAG := $(shell git describe --tags --abbrev=0 2> /dev/null || echo master)
# CURRENT_COMMIT := $(shell git rev-list -1 HEAD)

clean:
	rm -rf dist

build: 
	mkdir -p dist
	# echo "building $(CURRENT_TAG) ($(CURRENT_COMMIT))"
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -ldflags '-extldflags "-static" -s -w' -o dist/main
	chmod +x dist/main
	zip dist/main.zip dist/main
	rm -rf dist/main


all: clean build deploy

deploy: 
	aws lambda update-function-code \
		--function-name "$(AWS_FUNCTION_ARN)" \
		--zip-file fileb://$(CURRENT_DIR)/dist/main.zip