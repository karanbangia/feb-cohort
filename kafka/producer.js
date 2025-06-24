const {Kafka} = require('kafkajs');

const kafka = new Kafka({
    clientId: 'my-producer',
    brokers: ['localhost:9092'], // Update this if necessary
});

const topic = 'test-topic-2';
const partitions = [0, 1, 2]; // Define the partitions (Update based on your setup)
const producer = kafka.producer();

const run = async () => {
    await producer.connect();

    for (let i = 0; i < 1000000; i++) { // Send 10 messages
        const partition = partitions[i % partitions.length]; // Round-robin partition selection
        const message = `Message ${i} to partition ${partition}`;

        await producer.send({
            topic,
            messages: [{key: `key-${i}`, value: message, partition}],
        });

        console.log(`Sent: ${message}`);
    }

    await producer.disconnect();
};

run().catch(console.error);
