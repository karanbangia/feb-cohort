# Kafka Null vs Empty String Test with IBM Sarama

This project demonstrates the difference between **null values** and **empty strings** in Kafka messages using IBM Sarama Go client.

## Prerequisites

1. **Docker & Docker Compose** - to run Kafka locally
2. **Go 1.21+** - to run the Sarama producer and consumer

## Setup

### 1. Start Kafka

```bash
docker-compose up -d
```

This will start:

- Zookeeper on port 2181
- Kafka on port 9092
- Kafka UI on port 8080

### 2. Install Dependencies

```bash
go mod tidy
```

## Running the Test

### Step 1: Run the Producer

The producer will send 5 different test messages to demonstrate null vs empty string behavior:

```bash
go run producer_sarama.go
```

This will send:

1. Message with NULL value
2. Message with empty string value
3. Message with normal string value
4. Message with empty key and NULL value
5. Message with empty key and empty string value

### Step 2: Run the Consumer (in a separate terminal)

To see the difference between null and empty string values:

```bash
# Build and run consumer in a separate terminal
go build -o consumer_test consumer_test.go
./consumer_test test-group

# Or run directly:
go run consumer_test.go test-group
```

### Step 3: Alternative - Use Kafka Console Consumer

You can also use the built-in Kafka console consumer to see the raw messages:

```bash
# Get Kafka container ID
docker ps

# Run console consumer
docker exec -it <kafka-container-id> /opt/kafka/bin/kafka-console-consumer.sh \
  --bootstrap-server localhost:9092 \
  --topic test-topic \
  --from-beginning \
  --property print.key=true \
  --property print.value=true \
  --property print.timestamp=true \
  --property key.separator='|'
```

## Key Differences to Observe

### NULL Value vs Empty String Value

| Aspect           | NULL Value              | Empty String                   |
| ---------------- | ----------------------- | ------------------------------ |
| Sarama Value     | `nil`                   | `[]byte{}` (zero-length slice) |
| Wire Protocol    | No value bytes          | Zero-length byte array         |
| Console Consumer | Nothing after separator | Empty content but present      |
| Go Detection     | `message.Value == nil`  | `len(message.Value) == 0`      |

### Test Results

When you run the producer and consumer, you'll see output like:

```
📨 MESSAGE RECEIVED
Topic: test-topic | Partition: 0 | Offset: 0
Key: 'test-key-1' (length: 10)
Value: NULL (nil) ⚠️
Value Analysis: This is a true null value in Kafka
Raw Key Bytes: [116 101 115 116 45 107 101 121 45 49]
Raw Value Bytes: <nil>

📨 MESSAGE RECEIVED
Topic: test-topic | Partition: 0 | Offset: 1
Key: 'test-key-2' (length: 10)
Value: EMPTY STRING (zero-length byte slice) 📝
Value Analysis: This is an empty string, not null
Raw Key Bytes: [116 101 115 116 45 107 101 121 45 50]
Raw Value Bytes: []
```

## Code Structure

- **`producer_sarama.go`** - Sends test messages with different null/empty combinations
- **`consumer_test.go`** - Detailed consumer that analyzes message structure
- **`go.mod`** - Go module with Sarama dependency

## Understanding the Wire Protocol

At the Kafka wire protocol level:

- **NULL value**: No value bytes are sent (length = -1 in protocol)
- **Empty string**: Zero-length byte array is sent (length = 0 in protocol)

This distinction is important for:

- Schema evolution
- Data processing pipelines
- Compatibility with other Kafka clients
- Proper null handling in downstream systems

## Cleanup

```bash
# Stop Kafka
docker-compose down

# Remove built binaries
rm -f consumer_test
```

## Additional Notes

- Sarama preserves the distinction between null and empty values
- This behavior is consistent with the Java Kafka client
- Always check for `nil` before checking length to avoid panics
- Consider your downstream consumers when deciding between null vs empty
