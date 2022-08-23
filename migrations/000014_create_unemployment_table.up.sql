create table if not exists unemployment (
    time timestamp with time zone NOT NULL PRIMARY KEY,
    value double precision NOT NULL
);

SELECT create_hypertable('unemployment', 'time', chunk_time_interval => INTERVAL '1 year');

INSERT INTO economic_report (slug, display_name, description, unit, image, last_data_pull, initial_sync_delay_minutes)
VALUES('unemployment', 'Unemployment Rate', 'The monthly unemployment data of the United States. The unemployment rate represents the number of unemployed as a percentage of the labor force. Labor force data are restricted to people 16 years of age and older, who currently reside in 1 of the 50 states or the District of Columbia, who do not reside in institutions (e.g., penal and mental facilities, homes for the aged), and who are not on active duty in the Armed Forces','percent','images/fed.jpeg', NOW() - INTERVAL '7 day', 9);