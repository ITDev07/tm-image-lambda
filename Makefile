.PHONY: build clean deploy

build:
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/tm_score main.go utils.go process_scores.go key_maps.go

clean:
	rm -rf ./bin

deploy: clean build
	sls deploy --verbose
