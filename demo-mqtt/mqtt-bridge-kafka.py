import paho.mqtt.client as mqtt
from kafka import KafkaProducer
import json

producer = KafkaProducer(
    bootstrap_servers='160.191.49.50:9092',
    value_serializer=lambda v: v.encode('utf-8')
)

def on_message(client, userdata, msg):
    payload_str = msg.payload.decode('utf-8')
    print(f"→ Forwarding: {payload_str}")

    payload_dict = json.loads(payload_str)

    for sensor in payload_dict['sensors']:
        message = {
            "station_id": payload_dict["station_id"],
            "datetime": payload_dict["datetime"],
            "sensor_id": sensor["sensor_id"],
            "value": sensor["value"],
            "metric": sensor["metric"],
            "unit": sensor["unit"]
        }
        producer.send('sensor_data', value=json.dumps(message))
    # producer.send('sensor_data', payload)

client = mqtt.Client()
client.username_pw_set("test-mqtt-admin", "asd123")
client.connect("localhost", 1883, 60)
client.subscribe("stations/data")
client.on_message = on_message


print("Listening on MQTT and forwarding to Kafka...")
client.loop_forever()
