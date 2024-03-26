package dbConnection

import (
	"fl/my-portfolio/internal/configs"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type PostgresDB struct {
	config *configs.Config
}

func NewPGConnection(config *configs.Config) *PostgresDB {
	return &PostgresDB{
		config: config,
	}
}

func (p *PostgresDB) PostgreSQLConnection() (*sqlx.DB, error) {
	db, err := sqlx.Connect(p.config.DB.DBDialect, p.config.DB.DBAddress)
	if err != nil {
		return nil, fmt.Errorf("error, not connected to database, %w", err)
	}

	if err = db.Ping(); err != nil {
		defer db.Close()
		return nil, fmt.Errorf("error, not sent ping to database, %w", err)
	}

	return db, nil
}
