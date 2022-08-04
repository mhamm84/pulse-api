-- ####################################################################################################
-- retail_sales
-- ####################################################################################################
CREATE TABLE IF NOT EXISTS retail_sales (
    time TIMESTAMP WITH TIME ZONE NOT NULL PRIMARY KEY,
    value DOUBLE PRECISION NOT NULL
);

SELECT create_hypertable('retail_sales', 'time', chunk_time_interval => INTERVAL '1 year');

INSERT INTO economic_report (slug, display_name, description, unit, image, last_data_pull, initial_sync_delay_minutes) VALUES('retail_sales', 'Retail Sales', 'Advance Retail Sales: Retail Trade data of the United States.','millions of dollars','images/retail.jpeg', NOW() - INTERVAL '7 day', 1);