package bun

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
	"github.com/uptrace/bun/extra/bundebug"
)

type Store struct {
	_sqldb *sql.DB
	_db    *bun.DB
}

func New() *Store {
	s := &Store{}
	return s
}

// Connect to the database (apart of the Persist interface)
func (s *Store) Connect() (err error) {
	return s.ConnectSql()
}

// Connect to the database (apart of the Persist interface)
func (s *Store) ConnectSql() (err error) {
	// Connect to the database
	sqldb, err := sql.Open(sqliteshim.ShimName, "file::memory:?cache=shared")
	if err != nil {
		panic(err)
	}
	// Set the sqlite database
	s._sqldb = sqldb
	// Create bun database
	db := bun.NewDB(sqldb, sqlitedialect.New())
	// Set the bun database
	s._db = db
	// Enable logging
	s.EnableLogging()
	// Test the connection
	s.TestConnection()
	// Return any error
	return
}

// Connect to the database (apart of the Persist interface)
// "postgres://postgres:@localhost:5432/test?sslmode=disable"
func (s *Store) ConnectPostgres(dsn string) (err error) {
	// Connect to the database
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	// Set the sqlite database
	s._sqldb = sqldb
	// Create bun database
	db := bun.NewDB(sqldb, pgdialect.New())
	// Set the bun database
	s._db = db
	// Enable logging
	s.EnableLogging()
	// Test the connection
	s.TestConnection()
	// Return any error
	return
}

// Disconnect from the database (apart of the Persist interface)
func (s *Store) Disconnect() (err error) {
	return
}

// Enable logging for the database (apart of the Persist interface)
func (s *Store) EnableLogging() (err error) {
	s._db.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithVerbose(true),
		bundebug.FromEnv("BUNDEBUG"),
	))
	return
}

// Test the connection to the database (apart of the Persist interface)
func (s *Store) TestConnection() (err error) {
	var num int
	err = s._db.NewSelect().ColumnExpr("1").Scan(context.Background(), &num)
	if err != nil {
		return err
	}
	if num != 1 {
		return errors.New(fmt.Sprint("expected 1 got %d", num))
	}
	return
}

// Migrate the database (apart of the Persist interface)
func (s *Store) MigrateUsing(models ...interface{}) (err error) {

	for _, model := range models {
		// Create table.
		_, err := s._db.NewCreateTable().Model(model).Exec(context.Background())
		if err != nil {
			panic(err) // TODO: handle error and remove panic
		}
	}

	return

}

// Get the bun database
func (s *Store) GetBun() *bun.DB {
	return s._db
}

// Get the sql database
func (s *Store) SqlDb() *sql.DB {
	return s._sqldb
}
