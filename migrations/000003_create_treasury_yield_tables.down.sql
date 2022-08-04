DROP TABLE IF EXISTS treasury_yield_three_month;
DROP TABLE IF EXISTS treasury_yield_two_year;
DROP TABLE IF EXISTS treasury_yield_five_year;
DROP TABLE IF EXISTS treasury_yield_seven_year;
DROP TABLE IF EXISTS treasury_yield_ten_year;
DROP TABLE IF EXISTS treasury_yield_thirty_year;

DELETE FROM economic_report WHERE slug='treasury_yield_three_month';
DELETE FROM economic_report WHERE slug='treasury_yield_two_year';
DELETE FROM economic_report WHERE slug='treasury_yield_five_year';
DELETE FROM economic_report WHERE slug='treasury_yield_seven_year';
DELETE FROM economic_report WHERE slug='treasury_yield_ten_year';
DELETE FROM economic_report WHERE slug='treasury_yield_thirty_year';
