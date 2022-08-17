# pulse-api
Backend for the Pulse investing application - to keep the pulse on your investing

## Description
The Pulse project consists of an API written in Go and a UI written with Next.JS. The purpose is to build up an investing portal to have data points, stats and information across all assets types (stocks, bonds, commodities, FX, Bitcoin, Stacks, Ethereum) . It will include economic indicators and analysis to gain macro insights, data, stats, analysis and to be able to create a personal macro model to help assess the best places to invest and strategies to engage in.

## Economic Data
- CPI
- Consumer Sentiment
- Retail Sales
- Treasury Yields
- Real GDP
- Federal Funds Rate

## Tech Details
Some tech notes which I shall formalize more as the projects evolves

### API Environment Vars

#### Example .env

```
#Database
POSTGRES_USER=pulse_user
POSTGRES_PASSWORD=password
POSTGRES_DB=pulse
PULSE_POSTGRES_DSN=postgres://pulse_user:password@pulse_timescale_db:5432/pulse?sslmode=disable

# Data sync with third party API's
PULSE_DATA_SYNC_ENABLE=false

#CORS
PULSE_CORS_TRUSTED_ORIGIN=http://localhost:9090

#Logging
PULSE_LOG_LEVEL=INFO

# AlphaVantage Financial Data API
ALPHA_VANTAGE_BASE_URL=https://www.alphavantage.co/query
ALPHA_VANTAGE_API_TOKEN=you_token

#SMTP
PULSE_SMTP_HOST=smtp.mailtrap.io
PULSE_SMTP_PORT=2525
PULSE_SMTP_USERNAME=your_username
PULSE_SMTP_PASSWORD=your_password
PULSE_SMTP_SENDER=no-reply@pulseinvesting.support.net
```

## Make

### Run/Build
- Install make ```brew install make```

In top level project dir type ```make``` to see list of available commands and descriptions to build, run