create table if not exists inflation_expectation (
    time timestamp with time zone NOT NULL PRIMARY KEY,
    value double precision NOT NULL
);

SELECT create_hypertable('inflation_expectation', 'time', chunk_time_interval => INTERVAL '1 year');

INSERT INTO economic_report (slug, display_name, description, unit, image, last_data_pull, initial_sync_delay_minutes)
VALUES('inflation_expectation', 'Inflation Expectation', 'The monthly inflation expectation data of the United States, as measured by the median expected price change next 12 months according to the Surveys of Consumers by University of Michigan (Inflation Expectation© [MICH]), retrieved from FRED, Federal Reserve Bank of St. Louis.','percent','images/fed.jpeg', NOW() - INTERVAL '7 day', 11);