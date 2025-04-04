# python 3.11

import random
import os

from paho.mqtt import client as mqtt_client
from dotenv import load_dotenv
load_dotenv()


broker = os.getenv('MQTT_BROKER_HOST')
port = int(os.getenv('MQTT_BROKER_PORT'))
topic = "t/test"
# Generate a Client ID with the publish prefix.
client_id = f'publish-{random.randint(0, 1000)}'
username = os.getenv('MQTT_BROKER_USERNAME')
password = os.getenv('MQTT_BROKER_PASSWORD')

def connect_mqtt() -> mqtt_client:
    def on_connect(client, userdata, flags, rc):
        if rc == 0:
            print("Connected to MQTT Broker!")
        else:
            print("Failed to connect, return code %d\n", rc)

    client = mqtt_client.Client(client_id = client_id)
    client.username_pw_set(username, password)
    client.on_connect = on_connect
    client.connect(broker, port)
    return client


def subscribe(client: mqtt_client):
    def on_message(client, userdata, msg):
        print(f"Received `{msg.payload.decode()}` from `{msg.topic}` topic")

    client.subscribe(topic)
    client.on_message = on_message


def run():
    client = connect_mqtt()
    subscribe(client)
    client.loop_forever()


if __name__ == '__main__':
    run()
