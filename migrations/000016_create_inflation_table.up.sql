create table if not exists inflation (
    time timestamp with time zone NOT NULL PRIMARY KEY,
    value double precision NOT NULL
);

SELECT create_hypertable('inflation', 'time', chunk_time_interval => INTERVAL '10 year');

INSERT INTO economic_report (slug, display_name, description, unit, image, last_data_pull, initial_sync_delay_minutes)
VALUES('inflation', 'Inflation', 'The annual inflation rates (consumer prices) of the United States.','percent','images/fed.jpeg', NOW() - INTERVAL '7 day', 11);