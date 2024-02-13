# DOCUMENTATION
## INSTALLATION

### 1. install go
https://golang.org/doc/install

### 2. install redis
https://redis.io/download

### 3. install rabbitmq
https://www.rabbitmq.com/download.html

### 4. install mongodb
https://docs.mongodb.com/manual/installation/

## CONFIGURATION
- make sure redis is running on localhost:6379
- make sure rabbitmq is running on localhost:5672
- make sure mongodb is running on localhost:27017

## RUN
- clone from github.com/yudaph/pubsub
- run redis, rabbitmq, and mongodb
- run consumer with `go run consumer/main.go`
- run publisher with `go run publisher/main.go`
