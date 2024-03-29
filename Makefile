all: test build

vendor:
	GO111MODULE=on go mod vendor

build: vendor
	GO111MODULE=on go build -o out/clails main/*.go

test: vendor
	GO111MODULE=on go test github.com/hekonsek/clails/clails

docker-build: build
	docker build . -t hekonsek/clails

docker-push: docker-build
	docker push hekonsek/clails