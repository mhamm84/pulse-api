package config

type ApiConfig struct {
	Port         int
	Env          string
	DB           DbConfig
	AlphaVantage struct {
		BaseUrl string
		Token   string
	}
	Limiter struct {
		RPS     float64
		Burst   int
		Enabled bool
	}
	Cors struct {
		TrustedOrigins []string
	}
	DataSync bool
	LogLevel string
	SMTP     struct {
		Host     string
		Port     int
		Username string
		Password string
		Sender   string
	}
}

type DbConfig struct {
	Dsn          string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  string
}
