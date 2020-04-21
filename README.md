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

For Druid, Kafka service available as `kafka:9092`

## Sample query

```
curl -s -X 'POST' -H 'Content-Type:application/json' \
    -d @sample_query.json \
    http://localhost:8888/druid/v2/sql | jq
```
