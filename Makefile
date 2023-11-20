build:
	cd scripts && sh build.sh

tests:
	go test ./... -coverprofile cover.out && go tool cover -func cover.out