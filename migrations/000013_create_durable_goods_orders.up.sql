create table if not exists durable_goods_orders (
    time timestamp with time zone NOT NULL PRIMARY KEY,
    value double precision NOT NULL
);

SELECT create_hypertable('durable_goods_orders', 'time', chunk_time_interval => INTERVAL '1 year');

INSERT INTO economic_report (slug, display_name, description, unit, image, last_data_pull, initial_sync_delay_minutes)
VALUES('durable_goods_orders', 'Durable Goods Orders', 'the monthly manufacturers'' new orders of durable goods in the United States','millions of dollars','images/fed.jpeg', NOW() - INTERVAL '7 day', 9);