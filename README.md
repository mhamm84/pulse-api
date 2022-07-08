# pulse-api
Backend for the Pulse application - to keep the pulse on your investing

## Environment Vars

### Postgresql
PULSE_POSTGRES_DSN=postgres://USER:PWD@HOST/DB?sslmode=disable

### Alpha Vantage API
ALPHA_VANTAGE_BASE_URL=https://www.alphavantage.co/query
ALPHA_VANTAGE_API_TOKEN=YOUR_TOKEN


## Economic Data

### CPI
- Get the latest CPI value
- Get the % change from the previous value
- Get the last 12 months od CPI values

## Timescaledb
- CREATE DATABASE name
- CREATE EXTENSION IF NOT EXISTS timescaledb CASCADE;

## CORS
go run ./cmd/api/ -cors-trusted-origins="http://localhost:9090"
