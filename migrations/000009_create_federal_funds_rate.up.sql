
CREATE TABLE IF NOT EXISTS federal_funds_rate (
    time TIMESTAMP WITH TIME ZONE NOT NULL PRIMARY KEY,
    value DOUBLE PRECISION NOT NULL
);

SELECT create_hypertable('federal_funds_rate', 'time', chunk_time_interval => INTERVAL '1 year');

INSERT INTO economic_report (slug, display_name, description, unit, image, last_data_pull, initial_sync_delay_minutes) VALUES('federal_funds_rate', 'Federal Funds Rate', 'federal funds rate (interest rate) of the United States','percent','images/fed.jpeg', NOW() - INTERVAL '7 day', 7);