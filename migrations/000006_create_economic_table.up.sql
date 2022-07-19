-- ####################################################################################################
-- cpi
-- ####################################################################################################
CREATE TABLE IF NOT EXISTS cpi (
    time TIMESTAMP WITH TIME ZONE NOT NULL,
    value DOUBLE PRECISION NOT NULL
);

SELECT create_hypertable('cpi', 'time', chunk_time_interval => INTERVAL '1 year');
