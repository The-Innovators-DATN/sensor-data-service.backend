# python 3.11

import random
import time
from datetime import datetime
import os
import json
from paho.mqtt import client as mqtt_client
from dotenv import load_dotenv
load_dotenv()


broker = os.getenv('MQTT_BROKER_HOST')
port = int(os.getenv('MQTT_BROKER_PORT'))
topic = "stations/data"
# Generate a Client ID with the publish prefix.
client_id = f'publish-12'
username = os.getenv('MQTT_BROKER_USERNAME')
password = os.getenv('MQTT_BROKER_PASSWORD')

def connect_mqtt():
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


def publish(client):
    msg_count = 1
    metric_list = ['temperature', 'humidity', 'pressure', 'wind_speed', 'wind_direction', 'rainfall']
    while True:
        time.sleep(1)
        #TypeError: payload must be a string, bytearray, int, float or None.
        msg = {
            "metric_value": random.uniform(0,100),
            "metric": random.choice(metric_list),
            "station_id": random.randint(0, 10),
            "datetime": datetime.now().strftime('%Y-%m-%d %H:%M:%S')
        }
        result = client.publish(topic, json.dumps(msg))
        # result: [0, 1]
        status = result[0]
        if status == 0:
            print(f"Send `{msg}` to topic `{topic}`")
        else:
            print(f"Failed to send message to topic {topic}")
        msg_count += 1
        if msg_count > 1000:
            break


def run():
    client = connect_mqtt()
    client.loop_start()
    publish(client)
    client.loop_stop()


if __name__ == '__main__':
    run()
