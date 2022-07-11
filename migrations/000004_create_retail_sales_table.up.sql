-- ####################################################################################################
-- retail_sales
-- ####################################################################################################
CREATE TABLE IF NOT EXISTS retail_sales (
    time TIMESTAMP WITH TIME ZONE NOT NULL,
    value DOUBLE PRECISION NOT NULL
);

SELECT create_hypertable('retail_sales', 'time', chunk_time_interval => INTERVAL '1 year');