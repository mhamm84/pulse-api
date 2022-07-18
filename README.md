# pulse-api
Backend for the Pulse application - to keep the pulse on your investing

## Description
The Pulse project consists of an API written in Go and a UI written with Next.JS. The purpose is to build up an investing portal to have data points, stats and information across all assets types (stocks, bonds, commodities, FX, Bitcoin, Stacks, Ethereum) . It will include economic indicators and analysis to gain macro insights, data, stats, analysis and to be able to create a personal macro model to help assess the best places to invest and strategies to engage in.

## Economic Data
- CPI
- Consumer Sentiment
- Retail Sales
- Treasury Yields

## Tech Details
Some tech notes which I shall formalize more as the projects evolves

### API Environment Vars

#### Postgresql
- PULSE_POSTGRES_DSN=postgres://USER:PWD@HOST/DB?sslmode=disable
#### Alpha Vantage API
- ALPHA_VANTAGE_BASE_URL=https://www.alphavantage.co/query
- ALPHA_VANTAGE_API_TOKEN=YOUR_TOKEN

### Database - Postgresql
Notes for DB setup

#### CITEXT type
Add extension to add citext to postgresql

-CREATE EXTENSION citext;
#### Timescaledb
- CREATE EXTENSION IF NOT EXISTS timescaledb CASCADE;

#### Go DB Migrations
##### Create
- migrate create -seq -ext=.sql -dir=./migrations name-of-file
##### Up
- migrate -path=./migrations -database=$$EXAMPLE_DSN up
##### Down
- migrate -path=./migrations -database=$EXAMPLE_DSN down

## Run/Build
- Install make ```brew install make```
### Audit project
- ```make audit```
### Run the API
- ```make api/run```
### Build
- ```make api/build```
