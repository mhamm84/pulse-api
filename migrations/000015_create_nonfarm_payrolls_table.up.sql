create table if not exists nonfarm_payrolls (
    time timestamp with time zone NOT NULL PRIMARY KEY,
    value double precision NOT NULL
);

SELECT create_hypertable('nonfarm_payrolls', 'time', chunk_time_interval => INTERVAL '1 year');

INSERT INTO economic_report (slug, display_name, description, unit, image, last_data_pull, initial_sync_delay_minutes)
VALUES('nonfarm_payrolls', 'Nonfarm Payroll', 'The monthly US All Employees: Total Nonfarm (commonly known as Total Nonfarm Payroll), a measure of the number of U.S. workers in the economy that excludes proprietors, private household employees, unpaid volunteers, farm employees, and the unincorporated self-employed.','thousands of persons','images/fed.jpeg', NOW() - INTERVAL '7 day', 9);