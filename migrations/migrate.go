package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/pressly/goose"
	"github.com/spf13/viper"
)

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

func main() {
	if err := migrate(); err != nil {
		log.Fatal(err)
	}
}

func migrate() error {
	configPath := flag.String("config", "env/config", "config file")
	format := flag.String("format", "ini", "config format")

	flag.Usage = usage

	flag.Parse()
	args := flag.Args()

	config := viper.NewWithOptions(
		viper.EnvKeyReplacer(
			strings.NewReplacer(".", "_"),
		),
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
