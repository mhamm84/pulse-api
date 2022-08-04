-- ####################################################################################################
-- treasury_yield
-- ####################################################################################################
CREATE TABLE IF NOT EXISTS treasury_yield_three_month (
    time TIMESTAMP WITH TIME ZONE NOT NULL PRIMARY KEY,
    value DOUBLE PRECISION NOT NULL
);

SELECT create_hypertable('treasury_yield_three_month', 'time', chunk_time_interval => INTERVAL '1 month');

CREATE TABLE IF NOT EXISTS treasury_yield_two_year (
    time TIMESTAMP WITH TIME ZONE NOT NULL,
    value DOUBLE PRECISION NOT NULL
);

SELECT create_hypertable('treasury_yield_two_year', 'time', chunk_time_interval => INTERVAL '1 month');

CREATE TABLE IF NOT EXISTS treasury_yield_five_year (
    time TIMESTAMP WITH TIME ZONE NOT NULL,
    value DOUBLE PRECISION NOT NULL
);

SELECT create_hypertable('treasury_yield_five_year', 'time', chunk_time_interval => INTERVAL '1 month');

CREATE TABLE IF NOT EXISTS treasury_yield_seven_year (
    time TIMESTAMP WITH TIME ZONE NOT NULL,
    value DOUBLE PRECISION NOT NULL
);

SELECT create_hypertable('treasury_yield_seven_year', 'time', chunk_time_interval => INTERVAL '1 month');

CREATE TABLE IF NOT EXISTS treasury_yield_ten_year (
    time TIMESTAMP WITH TIME ZONE NOT NULL,
    value DOUBLE PRECISION NOT NULL
);

SELECT create_hypertable('treasury_yield_ten_year', 'time', chunk_time_interval => INTERVAL '1 month');

CREATE TABLE IF NOT EXISTS treasury_yield_thirty_year (
    time TIMESTAMP WITH TIME ZONE NOT NULL,
    value DOUBLE PRECISION NOT NULL
);

SELECT create_hypertable('treasury_yield_thirty_year', 'time', chunk_time_interval => INTERVAL '1 month');


INSERT INTO economic_report (slug, display_name, description, unit, image, last_data_pull, initial_sync_delay_minutes, extras) VALUES('treasury_yield_three_month', '3M Treasury Yield', 'The three month US treasury yield','percent','images/treasury.jpeg', NOW() - INTERVAL '7 day', 3, '{"maturity": "3month"}');
INSERT INTO economic_report (slug, display_name, description, unit, image, last_data_pull, initial_sync_delay_minutes, extras) VALUES('treasury_yield_two_year', '2Y Treasury Yield', 'The two year US treasury yield','percent','images/treasury.jpeg', NOW() - INTERVAL '7 day', 3, '{"maturity": "2year"}');
INSERT INTO economic_report (slug, display_name, description, unit, image, last_data_pull, initial_sync_delay_minutes, extras) VALUES('treasury_yield_five_year', '5Y Treasury Yield', 'The five year US treasury yield','percent','images/treasury.jpeg', NOW() - INTERVAL '7 day', 3, '{"maturity": "5year"}');
INSERT INTO economic_report (slug, display_name, description, unit, image, last_data_pull, initial_sync_delay_minutes, extras) VALUES('treasury_yield_seven_year', '7Y Treasury Yield', 'The seven year US treasury yield','percent','images/treasury.jpeg', NOW() - INTERVAL '7 day', 5, '{"maturity": "7year"}');
INSERT INTO economic_report (slug, display_name, description, unit, image, last_data_pull, initial_sync_delay_minutes, extras) VALUES('treasury_yield_ten_year', '10Y Treasury Yield', 'The ten year US treasury yield','percent','images/treasury.jpeg', NOW() - INTERVAL '7 day', 5, '{"maturity": "10year"}');
INSERT INTO economic_report (slug, display_name, description, unit, image, last_data_pull, initial_sync_delay_minutes, extras) VALUES('treasury_yield_thirty_year', '30Y Treasury Yield', 'The thirty year US treasury yield','percent','images/treasury.jpeg', NOW() - INTERVAL '7 day', 5, '{"maturity": "30year"}');
