# Golang x RabbitMQ Example

An example implementation of Rabbit Message Queue in Golang

## How to Demo

### 1. Run RabbitMQ

```
docker-compose up -d
```

### 2. Receive Message
```
go run receive_message.go
```

### 3. Send Message
```
go run send_message.go
```