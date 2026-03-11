package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"github.com/lumi4x/news/src/common/config"
	"github.com/pkg/errors"
)

func ApplyMigrations(ctx context.Context, d source.Driver, databaseName string) error {
	err := createDatabaseIfNotExists(ctx, databaseName)
	if err != nil {
		return errors.Wrapf(err, "failed creating database")
	}

	m, err := migrate.NewWithSourceInstance("iofs", d, config.PostgresConnectionString(databaseName))
	if err != nil {
		return errors.Wrapf(err, "failed creating new migrations source")
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return errors.Wrapf(err, "failed applying migrations. if error says fix and force, fix your migration, manually run the migration commands that didn't run and set dirty=false in the database migrations table. more info")
	}

	return nil
}

func createDatabaseIfNotExists(ctx context.Context, databaseName string) error {
	db, err := sql.Open("postgres", config.PostgresConnectionStringWithoutDatabase())
	if err != nil {
		return errors.Wrapf(err, "failed connecting to postgres")
	}

	_, err = db.ExecContext(ctx, fmt.Sprintf("CREATE DATABASE %s", pq.QuoteIdentifier(databaseName)))
	if err != nil {
		if err, ok := err.(*pq.Error); ok && err.Code.Name() != "duplicate_database" {
			return errors.Wrapf(err, "failed creating database '%s'", databaseName)
		}
	}

	return nil
}
