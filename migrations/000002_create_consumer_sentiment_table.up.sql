-- ####################################################################################################
-- consumer_sentiment
-- ####################################################################################################
CREATE TABLE IF NOT EXISTS consumer_sentiment (
    time TIMESTAMP WITH TIME ZONE NOT NULL,
    value DOUBLE PRECISION NOT NULL
);

SELECT create_hypertable('consumer_sentiment', 'time', chunk_time_interval => INTERVAL '1 year');