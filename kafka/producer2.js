const { Kafka } = require("kafkajs");

const kafka = new Kafka({
  clientId: "my-producer",
  brokers: ["localhost:9092"], // Update this if necessary
});

const topic = "25-june";
const producer = kafka.producer();

const run = async () => {
  await producer.connect();

  for (let i = 0; i < 10; i++) {
    // Send 10 messages
    const message = `Message ${i}`;

    await producer.send({
      topic,
      messages: [{ key: null, value: message }],
    });

    console.log(`Sent: ${message}`);
  }

  await producer.disconnect();
};

run().catch(console.error);
