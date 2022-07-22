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
