SET allow_unrestricted_reads_from_keeper = 1; -- Allow unrestricted reads from ZooKeeper to avoid the error
SELECT * FROM system.clusters;

CREATE DATABASE station_metrics ON CLUSTER company_cluster;

drop table if exists station_metrics.messages_local  ON CLUSTER company_cluster;
drop table if exists station_metrics.messages_sharded  ON CLUSTER company_cluster;
drop table if exists station_metrics.sensors_to_kafka  ON CLUSTER company_cluster;
drop table if exists station_metrics.sensors_to_kafka_mv  ON CLUSTER company_cluster;

CREATE TABLE station_metrics.messages_local ON CLUSTER company_cluster
(
    value Float32,
    metric_id Int,
    station_id Int,
    unit String,
    predicted_value Float32,
    local_error Float32,
    trend_anomaly Bool,
    datetime DateTime64(0)
)
ENGINE = MergeTree()
PARTITION BY toYYYYMM(datetime)  -- Optional, for partitioning by month
ORDER BY (datetime); -- Define order for efficient querying

-- Create a distributed table referencing the local table on all nodes
-- The table will be created on all nodes and will be used to query the local table on all nodes
CREATE TABLE station_metrics.messages_sharded ON CLUSTER company_cluster AS station_metrics.messages_local
ENGINE = Distributed(company_cluster, station_metrics, messages_local, station_id);


CREATE TABLE station_metrics.sensors_to_kafka ON CLUSTER company_cluster (
    value Float32,
    metric_id Int,
    station_id Int,
    unit String,
    predicted_value Float32,
    local_error Float32,
    trend_anomaly Bool,
    datetime DateTime64(0)
) ENGINE = Kafka
SETTINGS kafka_broker_list = '103.172.79.28.:9092',
         kafka_topic_list = 'sensor_metric',
         kafka_format = 'JSONEachRow',
         kafka_group_name = 'clickhouse_consumer',
         kafka_skip_broken_messages = 1,
         kafka_num_consumers = 2;

CREATE MATERIALIZED VIEW station_metrics.sensors_to_kafka_mv ON CLUSTER company_cluster
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
SELECT * FROM system.clusters WHERE cluster = 'company_cluster';



