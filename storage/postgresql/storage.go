package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"boilgopher/storage/postgresutil"
)

const dbConnEnv = "DATABASE_CONNECTION"

type GetterContext interface {
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

// Storage provides a wrapper around an sql database and provides
// required methods for interacting with the database
type Storage struct {
	logger logrus.FieldLogger
	db     *sqlx.DB

	// dataEncryptionKey is the key used for protecting the secret fields
	dataEncryptionKey []byte
	defaultPageSize   int
}

// GetDBConn returns the underlying sql.DB object
func (s Storage) GetDBConn() *sql.DB {
	return s.db.DB
}

// SetDataEncryptionKey stores the DEK in the storage layer
// created/fetched by secret handler.
// This allows the table field to be encrypted and decrypted
func (s *Storage) SetDataEncryptionKey(key []byte) {
	s.dataEncryptionKey = key
}

func newTestStorage(t *testing.T) (*Storage, func()) {
	ddlConnStr := os.Getenv(dbConnEnv)
	if ddlConnStr == "" {
		t.Skipf("%s is not set", dbConnEnv)
	}

	db, teardown := postgresutil.MustNewDevelopmentDB(ddlConnStr, filepath.Join("..", "..", "migrations", "sql"))
	return &Storage{db: db, defaultPageSize: 10}, teardown
}

func envTest(t *testing.T) {
	config := viper.NewWithOptions(
		viper.EnvKeyReplacer(
			strings.NewReplacer(".", "_"),
		),
	)
	config.SetConfigFile("../../env/config")
	config.SetConfigType("ini")
	config.AutomaticEnv()
	if err := config.ReadInConfig(); err != nil {
		t.Fatalf("error loading configuration: %v", err)
	}
	dConfig, err := postgresutil.NewDBFromConfig(config)
	if err != nil {
		t.Fatalf("error configure db, please set your config first")
	}
	// DATABASE_CONNECTION="user={user} host={host} port={port} dbname={dbname} password={password} sslmode=disable"
	t.Setenv(dbConnEnv, fmt.Sprintf("user=%s host=%s port=%s dbname=%s password=%s sslmode=%s",
		dConfig.User,
		dConfig.Host,
		dConfig.Port,
		dConfig.DBName,
		dConfig.Password,
		dConfig.SSLMode,
	))
}
