# https://goswagger.io/install.html
.PHONY: gen-swagger
gen-swagger: 
	swagger generate spec -o ./api/dns/swagger.yml ./cmd/dns/

.PHONY: lint
lint:
	golangci-lint run -c .golangcilint.yml

.PHONY: test
test: 
	go test -coverpkg=./... -coverprofile=profile.cov ./...
	go tool cover -func profile.cov	| grep total

build:
	CGO_ENABLED=0 go build -o ./app/dns ./cmd/dns/

build-image:
	docker build -t cauchy2384/space2218:latest . -f ./deployment/Dockerfile