-- ####################################################################################################
-- consumer_sentiment
-- ####################################################################################################
CREATE TABLE IF NOT EXISTS consumer_sentiment (
    time TIMESTAMP WITH TIME ZONE NOT NULL PRIMARY KEY,
    value DOUBLE PRECISION NOT NULL
);

SELECT create_hypertable('consumer_sentiment', 'time', chunk_time_interval => INTERVAL '1 year');

INSERT INTO economic_report (slug, display_name, description, unit, image, last_data_pull, initial_sync_delay_minutes) VALUES('consumer_sentiment', 'Consumer Sentiment', 'Consumer sentiment and confidence data of the United States', 'index 1966:Q1=100','images/consumer.jpeg', NOW() - INTERVAL '7 day', 1);