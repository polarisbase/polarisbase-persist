package bun

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/driver/sqliteshim"
	"github.com/uptrace/bun/extra/bundebug"
	"strings"
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

	return s.MigrateUsingWithNamespace("", models...)

}

// Migrate
func (s *Store) Migrate(namespace string, model interface{}) (tableName string, err error) {

	if namespace == "" {
		namespace = ""
	} else {
		namespace = namespace + "_"
	}

	// get the model name
	modelName := fmt.Sprintf("%T\n", model)
	// remove and '.' and replace with '_'
	modelName = strings.ReplaceAll(modelName, ".", "_")
	// remove the '*'
	modelName = strings.ReplaceAll(modelName, "*", "")

	// tableName
	tableName = namespace + strings.ToLower(modelName)

	// Create table.
	_, err = s._db.NewCreateTable().
		Model(model).
		ModelTableExpr(
			tableName,
		).
		Exec(context.Background())
	if err != nil {
		panic(err) // TODO: handle error and remove panic
	}

	return
}

// Migrate the database (apart of the Persist interface)
func (s *Store) MigrateUsingWithNamespace(namespace string, models ...interface{}) (err error) {

	prefix := namespace + "_"
	if namespace == "" {
		prefix = ""
	}

	_ = prefix

	for _, model := range models {

		// get the model name
		modelName := fmt.Sprintf("%T\n", model)
		// remove and '.' and replace with '_'
		modelName = strings.ReplaceAll(modelName, ".", "_")
		// remove the '*'
		modelName = strings.ReplaceAll(modelName, "*", "")

		fmt.Printf("Modle Name: %s", modelName)

		//_, err := s._db.NewCreateTable().Model(model).Exec(context.Background())

		// Create table.
		_, err := s._db.NewCreateTable().
			Model(model).
			ModelTableExpr(
				prefix + strings.ToLower(modelName),
			).
			Exec(context.Background())
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
