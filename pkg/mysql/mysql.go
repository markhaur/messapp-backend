package mysql

import (
	"context"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/pkg/errors"
)

const dbName = "messapp"

func NewDB(ctx context.Context, dataSourceName string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return nil, err
	}

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil
}

func Migrate(migrationURL string, db *sql.DB) error {
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		return errors.Wrap(err, "could not create postgres driver instance")
	}

	migration, err := migrate.NewWithDatabaseInstance(migrationURL, dbName, driver)
	if err != nil {
		return errors.Wrap(err, "could not create migration instance")
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		return errors.Wrap(err, "failed to run migrate up")
	}

	return nil
}
