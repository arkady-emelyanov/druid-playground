# Druid playground

## Generate events
```
go run ./generate_events.go
go run ./generate_merchants.go
```

## Start cluster
```
docker-compose up
```

## Load events to Kafka
```
kafkacat -z snappy -b 127.0.0.1:9093 -t merchants -P -l ./generated/merchants.json
kafkacat -z snappy -b 127.0.0.1:9093 -t events -P -l ./generated/events.json
```

## Setup ingestion

http://localhost:8888/unified-console.html#load-data
