package configs

type DBConfig struct {
	DBAddress string
	DBDialect string
}

type CacherConfig struct {
	CacherHost string
	CacherPort string
	CacherDB   int
}

type ExchangeConfig struct {
	BinanceAPIKey    string
	BinanceSecretKey string

	BybitAPIKey    string
	BybitSecretKey string

	OkxAPIKey     string
	OkxSecretKey  string
	OkxPassphrase string
}

type CMCConfig struct {
	APIKey string
}

type Config struct {
	DB     DBConfig
	Cacher CacherConfig

	CMC      CMCConfig
	Exchange ExchangeConfig

	Address   string
	LogLevel  string
	SecretKey string
}
