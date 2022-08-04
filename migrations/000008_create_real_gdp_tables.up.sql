
CREATE TABLE IF NOT EXISTS real_gdp (
    time TIMESTAMP WITH TIME ZONE NOT NULL PRIMARY KEY,
    value DOUBLE PRECISION NOT NULL
);

SELECT create_hypertable('real_gdp', 'time', chunk_time_interval => INTERVAL '1 year');

CREATE TABLE IF NOT EXISTS real_gdp_per_capita (
    time TIMESTAMP WITH TIME ZONE NOT NULL PRIMARY KEY,
    value DOUBLE PRECISION NOT NULL
);

SELECT create_hypertable('real_gdp_per_capita', 'time', chunk_time_interval => INTERVAL '1 year');

INSERT INTO economic_report (slug, display_name, description, unit, image, last_data_pull, initial_sync_delay_minutes) VALUES('real_gdp', 'Real GDP', 'Real Gross Domestic Product of the United States','billions of dollars','images/real_gdp.jpeg', NOW() - INTERVAL '7 day', 7);
INSERT INTO economic_report (slug, display_name, description, unit, image, last_data_pull, initial_sync_delay_minutes) VALUES('real_gdp_per_capita', 'Real GDP Per Capita', 'Real GDP per Capita data of the United States','chained 2012 dollars','images/real_gdp_per_capita.jpeg', NOW() - INTERVAL '7 day', 7);