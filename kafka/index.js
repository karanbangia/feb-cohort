const { Kafka } = require('kafkajs');

const kafka = new Kafka({
    clientId: 'my-consumer',
    brokers: ['localhost:9092'], // Update if using a different host/port
});

const topic = 'test-topic-2';
const consumer = kafka.consumer({ groupId: 'test-group-3' });

const run = async () => {
    await consumer.connect();
    await consumer.subscribe({ topic, fromBeginning: false });

    await consumer.run({
        eachMessage: async ({ topic, partition: msgPartition, message, }) => {
                console.log(`Received message from partition ${msgPartition}:`, message.value.toString());

        },
    });


};

run().catch(console.error);

