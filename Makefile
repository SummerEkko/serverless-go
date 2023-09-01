.PHONY: build clean deploy gomodgen

build: gomodgen
	export GO111MODULE=on
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/createItem createItem/main.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/deleteItem deleteItem/main.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/updateItem updateItem/main.go
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/getItem getItem/main.go

clean:
	rm -rf ./bin ./vendor go.sum

deploy: clean build
	sls deploy --verbose

gomodgen:
	chmod u+x gomod.sh
	./gomod.sh


# $env:GOARCH = "amd64"
# $env:GOOS = "linux"
# go build -ldflags "-s -w" -o bin/hello hello/main.go
# go build -ldflags "-s -w" -o bin/world world/main.go
