import paho.mqtt.client as mqtt
from kafka import KafkaProducer

producer = KafkaProducer(
    bootstrap_servers='160.191.49.50:9092',
    value_serializer=lambda v: v.encode('utf-8')
)

def on_message(client, userdata, msg):
    payload = msg.payload.decode('utf-8')
    print(f"→ Forwarding: {payload}")
    producer.send('sensor_data', payload)

client = mqtt.Client()
client.username_pw_set("test-mqtt-admin", "asd123")
client.connect("localhost", 1883, 60)
client.subscribe("stations/data")
client.on_message = on_message


print("Listening on MQTT and forwarding to Kafka...")
client.loop_forever()
