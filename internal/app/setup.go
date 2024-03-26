package core

import (
	"fl/my-portfolio/internal/configs"
	dbConnection "fl/my-portfolio/pkg/db_connection"

	"os"
	"strconv"
	"sync"
	"time"

	validator "github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

var (
	config     *configs.Config
	configOnce = new(sync.Once)

	logger     *logrus.Logger
	loggerOnce = new(sync.Once)

	v     *validator.Validate
	vOnce = new(sync.Once)

	database     *sqlx.DB
	databaseOnce = new(sync.Once)

	key     []byte
	keyOnce = new(sync.Once)

	location     *time.Location
	locationOnce = new(sync.Once)
)

func GetConfig() *configs.Config {
	configOnce.Do(
		func() {
			if err := godotenv.Load(); err != nil {
				panic(err)
			}

			cacherDB, err := strconv.Atoi(os.Getenv("CACHER_DB"))
			if err != nil {
				panic(err)
			}

			config = &configs.Config{
				Address:   os.Getenv("ADDRESS"),
				LogLevel:  os.Getenv("LOG_LEVEL"),
				SecretKey: os.Getenv("SECRET_KEY"),
				DB: configs.DBConfig{
					DBAddress: os.Getenv("DB_ADDRESS"),
					DBDialect: os.Getenv("DB_DIALECT"),
				},
				Cacher: configs.CacherConfig{
					CacherHost: os.Getenv("CACHER_HOST"),
					CacherPort: os.Getenv("CACHER_PORT"),
					CacherDB:   cacherDB,
				},
				Exchange: configs.ExchangeConfig{
					BinanceAPIKey: os.Getenv("BINANCE_API_KEY"),
					BinanceSecretKey: os.Getenv("BINANCE_SECRET_KEY"),
					BybitAPIKey: os.Getenv("BYBIT_API_KEY"),
					BybitSecretKey: os.Getenv("BYBIT_SECRET_KEY"),
					OkxAPIKey: os.Getenv("OKX_API_KEY"),
					OkxSecretKey: os.Getenv("OKX_SECRET_KEY"),
					OkxPassphrase: os.Getenv("OKX_PASSPHRASE"),
				},
				CMC: configs.CMCConfig{
					APIKey: os.Getenv("CMC_API_KEY"),
				},
			}
		},
	)

	return config
}

func GetLogger() *logrus.Logger {
	loggerOnce.Do(
		func() {
			level, err := logrus.ParseLevel(GetConfig().LogLevel)
			if err != nil {
				panic(err)
			}

			logger = logrus.New()
			logger.SetLevel(level)
		},
	)

	return logger
}

func GetValidator() *validator.Validate {
	vOnce.Do(
		func() {
			v = validator.New()
		},
	)

	return v
}

func GetDB() *sqlx.DB {
	databaseOnce.Do(func() {
		conn := dbConnection.NewPGConnection(GetConfig())
		db, err := conn.PostgreSQLConnection()
		if err != nil {
			panic(err)
		}

		database = db
	})

	return database
}

func GetKey() []byte {
	keyOnce.Do(func() {
		// file, err := os.Open("jwt-private-key.pem")
		// if err != nil {
		// 	panic(err)
		// }

		// defer file.Close()

		// fileData, err := io.ReadAll(file)
		// if err != nil {
		// 	panic(err)
		// }

		// block, _ := pem.Decode(fileData)

		// keyParse, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		// if err != nil {
		// 	panic(err)
		// }

		// k, ok := keyParse.(*rsa.PrivateKey)
		// if !ok {
		// 	panic(err)
		// }
		key = []byte(GetConfig().SecretKey)
	})

	return key
}

func GetLocation() *time.Location {
	locationOnce.Do(func() {
		l, err := time.LoadLocation("Europe/Moscow")
		if err != nil {
			panic(err)
		}

		location = l
	})

	return location
}
