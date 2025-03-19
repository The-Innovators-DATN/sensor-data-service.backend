SET allow_unrestricted_reads_from_keeper = 1; -- Allow unrestricted reads from ZooKeeper to avoid the error
SELECT * FROM system.clusters;

CREATE DATABASE station_metrics ON CLUSTER StationCluster;

drop table if exists station_metrics.messages_local  ON CLUSTER StationCluster;
drop table if exists station_metrics.messages_sharded  ON CLUSTER StationCluster;
drop table if exists station_metrics.sensors_to_kafka  ON CLUSTER StationCluster;
drop table if exists station_metrics.sensors_to_kafka_mv  ON CLUSTER StationCluster;

CREATE TABLE station_metrics.messages_local ON CLUSTER StationCluster
(
    value Float32,
    metric String,
    station_id Int64,
    timestamp DateTime64(0)
)
ENGINE = MergeTree()
PARTITION BY toYYYYMM(timestamp)  -- Optional, for partitioning by month
ORDER BY (station_id, timestamp); -- Define order for efficient querying

-- Create a distributed table referencing the local table on all nodes
-- The table will be created on all nodes and will be used to query the local table on all nodes
CREATE TABLE station_metrics.messages_sharded ON CLUSTER StationCluster AS station_metrics.messages_local
ENGINE = Distributed(StationCluster, station_metrics, messages_local, station_id);


CREATE TABLE station_metrics.sensors_to_kafka ON CLUSTER StationCluster (
    value Float32,
    metric String,
    station_id Int64,
    timestamp DateTime64(0)
) ENGINE = Kafka
SETTINGS kafka_broker_list = '160.191.49.50:9092',
         kafka_topic_list = 'sensor_data',
         kafka_format = 'JSONEachRow',
         kafka_group_name = 'clickhouse_consumer',
         kafka_skip_broken_messages = 1,
         kafka_num_consumers = 2;

CREATE MATERIALIZED VIEW station_metrics.sensors_to_kafka_mv ON CLUSTER StationCluster
TO station_metrics.messages_sharded AS
SELECT * FROM station_metrics.sensors_to_kafka;

-- Check the table is created and distributed properly on all nodes yet?
SELECT * FROM system.databases WHERE name = 'station_metrics';

-- Check the table is created and distributed properly on all nodes yet?
SELECT database, name, engine, cluster FROM system.tables 
WHERE database = 'station_metrics';

-- Check operations on Kafka
SELECT * FROM system.kafka_consumers;

-- Check cluster operations on all nodes
SELECT * FROM system.clusters WHERE cluster = 'StationCluster';



