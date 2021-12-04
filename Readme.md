# Readme

## Assumptions

* Request properties are Strings, response property is a Number
* Response location is rounded to 2 decimal places
* Healthcheck & readiness criterias are yet to be discussed

## Build & Run
* Local: 
    ```
    make build
    ./app/dns
    ```
* Dockerized: 
    ```
    make build-image
    docker-compose -f ./deployment/docker-compose.yml up
    ```

## Tests
* Lint:
    ```
    make lint
    ```
* Unit: 
    ```
    make test
    ```
* Manual integration: in ./test/dns.http