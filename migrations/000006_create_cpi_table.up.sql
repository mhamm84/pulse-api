-- ####################################################################################################
-- cpi
-- ####################################################################################################
CREATE TABLE IF NOT EXISTS cpi (
    time TIMESTAMP WITH TIME ZONE NOT NULL PRIMARY KEY,
    value DOUBLE PRECISION NOT NULL
);

SELECT create_hypertable('cpi', 'time', chunk_time_interval => INTERVAL '1 year');

INSERT INTO economic_report (slug, display_name, description, unit, image, last_data_pull, initial_sync_delay_minutes) VALUES('cpi', 'CPI', 'Consumer Price Index (CPI) of the United States. CPI is widely regarded as the barometer of inflation levels in the broader economy.','index 1982-1984=100','images/inflation.jpeg', NOW() - INTERVAL '7 day', 1);