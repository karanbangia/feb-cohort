docker exec -it 437d2f54780c /opt/kafka/bin/kafka-console-consumer.sh \
--bootstrap-server localhost:9092 \
--topic test-topic \
--group my-group \
--from-beginning


docker exec -it 437d2f54780c /opt/kafka/bin/kafka-console-producer.sh \
--broker-list localhost:9092 \
--topic test-topic \
--property parse.key=true \
--property key.separator=:
