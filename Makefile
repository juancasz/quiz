build:
	cd scripts && sh build.sh

tests:
	go test -cover -coverpkg=./... -coverprofile=coverage.out ./... && go tool cover -html=coverage.out