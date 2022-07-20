module github.com/mhamm84/pulse-api

go 1.18

replace github.com/mhamm84/gofinance-alpha => ../../gofinance/gofinance-alpha

require (
	github.com/common-nighthawk/go-figure v0.0.0-20210622060536-734e95fb86be
	github.com/jmoiron/sqlx v1.3.5
	github.com/julienschmidt/httprouter v1.3.0
	github.com/lib/pq v1.10.6
	github.com/mhamm84/gofinance-alpha v0.0.0-00010101000000-000000000000
	github.com/pkg/errors v0.9.1
	github.com/shopspring/decimal v1.3.1
	github.com/stretchr/testify v1.8.0
	golang.org/x/time v0.0.0-20220609170525-579cf78fd858
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/kr/pretty v0.1.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/objx v0.4.0 // indirect
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
