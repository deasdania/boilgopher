package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"log"
	"strings"

	_ "github.com/lib/pq"
	"github.com/pressly/goose"
	"github.com/spf13/viper"
)

func main() {
	if err := migrate(); err != nil {
		log.Fatal(err)
	}
}

const driver = "postgres"

type DBConfig struct {
	User     string `mapstructure:"user"`
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	DBName   string `mapstructure:"dbname"`
	Password string `mapstructure:"password"`
	SSLMode  string `mapstructure:"sslMode"`
}

func usage() {
	const (
		usageRun      = `goose [OPTIONS] COMMAND`
		usageCommands = `
Commands:
    up                   Migrate the DB to the most recent version available
    up-to VERSION        Migrate the DB to a specific VERSION
    down                 Roll back the version by 1
    down-to VERSION      Roll back to a specific VERSION
    redo                 Re-run the latest migration
    status               Dump the migration status for the current DB
    version              Print the current version of the database
    create NAME [sql|go] Creates new migration file with next version`
	)
	fmt.Println(usageRun)
	flag.PrintDefaults()
	fmt.Println(usageCommands)
}

func migrate() error {
	configPath := flag.String("config", "env/config", "config file")
	format := flag.String("format", "ini", "config format")

	flag.Usage = usage

	flag.Parse()
	args := flag.Args()

	config := viper.NewWithOptions(
		viper.EnvKeyReplacer(strings.NewReplacer(".", "_")),
	)
	config.SetConfigFile(*configPath)
	config.SetConfigType(*format)
	config.AutomaticEnv()
	if err := config.ReadInConfig(); err != nil {
		return fmt.Errorf("error loading configuration for migration: %v", err)
	}

	db, err := Open(config)
	if err != nil {
		return fmt.Errorf("error opening db connection: %v", err)
	}
	defer func() { _ = db.Close() }()

	if err = goose.SetDialect(driver); err != nil {
		return fmt.Errorf("failed to set goose dialect: %v", err)
	}

	if len(args) == 0 {
		return errors.New("expected at least one arg")
	}

	command := args[0]

	migrationDir := config.GetString("database.migrationDir")
	if err = goose.Run(command, db, migrationDir, args[1:]...); err != nil {
		return fmt.Errorf("goose run: %v", err)

	}
	return db.Close()
}

func Open(config *viper.Viper) (*sql.DB, error) {
	dbString, err := NewDBStringFromConfig(config)
	if err != nil {
		return nil, err
	}

	db, err := sql.Open(driver, dbString)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// NewDBStringFromConfig build database connection string from config file.
func NewDBStringFromConfig(config *viper.Viper) (string, error) {
	var allConfig struct {
		Database DBConfig `mapstructure:"database"`
	}
	if err := config.Unmarshal(&allConfig); err != nil {
		return "", fmt.Errorf("cannot unmarshal db config: %w", err)
	}

	return NewDBStringFromDBConfig(allConfig.Database)
}

func NewDBStringFromDBConfig(config DBConfig) (string, error) {
	var dbParams []string
	dbParams = append(dbParams, fmt.Sprintf("user=%s", config.User))
	dbParams = append(dbParams, fmt.Sprintf("host=%s", config.Host))
	dbParams = append(dbParams, fmt.Sprintf("port=%s", config.Port))
	dbParams = append(dbParams, fmt.Sprintf("dbname=%s", config.DBName))
	if password := config.Password; password != "" {
		dbParams = append(dbParams, fmt.Sprintf("password=%s", password))
	}
	dbParams = append(dbParams, fmt.Sprintf("sslmode=%s",
		config.SSLMode))
	return strings.Join(dbParams, " "), nil
}
