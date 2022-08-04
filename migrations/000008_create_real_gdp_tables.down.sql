DROP TABLE IF EXISTS real_gdp;
DROP TABLE IF EXISTS real_gdp_per_capita;

DELETE FROM economic_report WHERE slug='real_gdp';
DELETE FROM economic_report WHERE slug='real_gdp_per_capita';