# affise-test

## Tested on
* macOS 14.3.1
* docker 25.0.3

## Run server
1. Run docker
```
LOG_LEVEL=DEBUG make run
```
2. Send an example request
```
curl --location 'localhost:8080' \
--header 'Content-Type: application/json' \
--data '{
    "urls": ["https://google.com/", "https://github.com/", "https://www.wikipedia.org/"]
}
'
```
## Run tests
```
make test
```

## Run linter
```
make lint
```
