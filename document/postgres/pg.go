package postgres

import (
	persist "github.com/polarisbase/polarisbase-persist"
	"github.com/upper/db/v4"
	"github.com/upper/db/v4/adapter/postgresql"
	"reflect"
)

// Store is a store.
type Store struct {
	sess     db.Session
	settings postgresql.ConnectionURL
}

// New creates a new store.
func New(host string, database string, user string, password string) *Store {
	s := &Store{}

	s.settings = postgresql.ConnectionURL{
		Host:     host,
		Database: database,
		User:     user,
		Password: password,
		Options: map[string]string{
			"sslmode": "disable",
		},
	}

	return s
}

// Connect to the database (part of the Persist interface)
func (s *Store) Connect() (err error) {
	sess, err := postgresql.Open(s.settings)
	if err != nil {
		return err
	}
	s.sess = sess
	return
}

// Collections returns a list of collections.
func (s *Store) Collections() (collections []db.Collection, err error) {
	collections, err = s.sess.Collections()
	if err != nil {
		return collections, err
	}
	return
}

// Collection returns a collection by name.
func (s *Store) Collection(name string, model interface{}) (collection db.Collection, err error) {

	// Create the collection if it doesn't exist as a table in the database if it doesn't exist
	//_, err = s.sess.SQL().Exec("CREATE TABLE IF NOT EXISTS " + name + " (id varchar(255) NOT NULL, PRIMARY KEY (id))")
	//if err != nil {
	//	return collection, err
	//}

	// reflect the model to get all the fields
	modelType := reflect.TypeOf(model)

	// if the model is a pointer, get the type of the pointer
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}

	// Create the collection if it doesn't exist as a table in the database if it doesn't exist and make the columns too
	_, err = s.sess.SQL().Exec("CREATE TABLE IF NOT EXISTS " + name + " (id varchar(255) NOT NULL, PRIMARY KEY (id))")
	if err != nil {
		return collection, err
	}

	// Loop through the fields and create the columns
	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)

		// get the 'db' tag value if it exists
		dbTag := field.Tag.Get("db")
		if dbTag == "" {
			continue
		}

		// get the type of the field
		fieldType := field.Type.String()
		// database type
		databaseType := "varchar(255)"

		// if type is postgresql.JSONB then set the type to jsonb
		if fieldType == "postgresql.JSONB" {
			databaseType = "jsonb"
		}

		if fieldType == "string" {
			databaseType = "varchar(255)"
		}

		if fieldType == "int" {
			databaseType = "int"
		}

		if fieldType == "int64" {
			databaseType = "int"
		}

		if fieldType == "float64" {
			databaseType = "float"
		}

		if fieldType == "bool" {
			databaseType = "boolean"
		}

		if fieldType == "time.Time" {
			databaseType = "timestamp"
		}

		if fieldType == "*time.Time" {
			databaseType = "timestamp"
		}

		_, err = s.sess.SQL().Exec("ALTER TABLE " + name + " ADD COLUMN IF NOT EXISTS " + dbTag + " " + databaseType + ";")
		if err != nil {
			return collection, err
		}
	}

	collection = s.sess.Collection(name)
	if err != nil {
		return collection, err
	}
	return
}

// NewBucket returns a new bucket.
func (s *Store) NewBucket(name string) (bucket persist.Bucket, err error) {
	bucket = &Bucket{
		Store:     s,
		Namespace: name,
	}
	return
}
