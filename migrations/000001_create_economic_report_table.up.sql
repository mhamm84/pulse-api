CREATE EXTENSION IF NOT EXISTS citext;
CREATE EXTENSION IF NOT EXISTS timescaledb CASCADE;

-- ####################################################################################################
-- economic_report
-- ####################################################################################################
CREATE TABLE IF NOT EXISTS economic_report (
    slug TEXT PRIMARY KEY,
    display_name TEXT UNIQUE NOT NULL,
    description TEXT NOT NULL,
    image TEXT NOT NULL,
    last_data_pull DATE NOT NULL,
    initial_sync_delay_minutes INTEGER NOT NULL,
    extras JSONB
);

CREATE INDEX idx_economic_report_slug ON economic_report(slug);

INSERT INTO economic_report (slug, display_name, description, image, last_data_pull, initial_sync_delay_minutes) VALUES('cpi', 'CPI', 'Consumer Price Index (CPI) of the United States. CPI is widely regarded as the barometer of inflation levels in the broader economy.','images/inflation.jpeg', NOW() - INTERVAL '7 day', 1);
INSERT INTO economic_report (slug, display_name, description, image, last_data_pull, initial_sync_delay_minutes) VALUES('consumer_sentiment', 'Consumer Sentiment', 'Consumer sentiment and confidence data of the United States','images/consumer.jpeg', NOW() - INTERVAL '7 day', 1);
INSERT INTO economic_report (slug, display_name, description, image, last_data_pull, initial_sync_delay_minutes) VALUES('retail_sales', 'Retail Sales', 'Advance Retail Sales: Retail Trade data of the United States.','images/retail.jpeg', NOW() - INTERVAL '7 day', 1);
INSERT INTO economic_report (slug, display_name, description, image, last_data_pull, initial_sync_delay_minutes, extras) VALUES('treasury_yield_three_month', '3M Treasury Yield', 'The three month US treasury yield','images/treasury.jpeg', NOW() - INTERVAL '7 day', 3, '{"maturity": "3month"}');
INSERT INTO economic_report (slug, display_name, description, image, last_data_pull, initial_sync_delay_minutes, extras) VALUES('treasury_yield_two_year', '2Y Treasury Yield', 'The two year US treasury yield','images/treasury.jpeg', NOW() - INTERVAL '7 day', 3, '{"maturity": "2year"}');
INSERT INTO economic_report (slug, display_name, description, image, last_data_pull, initial_sync_delay_minutes, extras) VALUES('treasury_yield_five_year', '5Y Treasury Yield', 'The five year US treasury yield','images/treasury.jpeg', NOW() - INTERVAL '7 day', 3, '{"maturity": "5year"}');
INSERT INTO economic_report (slug, display_name, description, image, last_data_pull, initial_sync_delay_minutes, extras) VALUES('treasury_yield_seven_year', '7Y Treasury Yield', 'The seven year US treasury yield','images/treasury.jpeg', NOW() - INTERVAL '7 day', 5, '{"maturity": "7year"}');
INSERT INTO economic_report (slug, display_name, description, image, last_data_pull, initial_sync_delay_minutes, extras) VALUES('treasury_yield_ten_year', '10Y Treasury Yield', 'The ten year US treasury yield','images/treasury.jpeg', NOW() - INTERVAL '7 day', 5, '{"maturity": "10year"}');
INSERT INTO economic_report (slug, display_name, description, image, last_data_pull, initial_sync_delay_minutes, extras) VALUES('treasury_yield_thirty_year', '30Y Treasury Yield', 'The thirty year US treasury yield','images/treasury.jpeg', NOW() - INTERVAL '7 day', 5, '{"maturity": "30year"}');