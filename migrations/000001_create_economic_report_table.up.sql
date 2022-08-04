CREATE EXTENSION IF NOT EXISTS citext;
CREATE EXTENSION IF NOT EXISTS timescaledb CASCADE;

-- ####################################################################################################
-- economic_report
-- ####################################################################################################
CREATE TABLE IF NOT EXISTS economic_report (
    slug TEXT PRIMARY KEY,
    display_name TEXT UNIQUE NOT NULL,
    description TEXT NOT NULL,
    unit TEXT NOT NULL,
    image TEXT NOT NULL,
    last_data_pull DATE NOT NULL,
    initial_sync_delay_minutes INTEGER NOT NULL,
    extras JSONB
);

CREATE INDEX idx_economic_report_slug ON economic_report(slug);
