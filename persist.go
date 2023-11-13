package persist

type Persist interface {
	Connect() (err error)
	Disconnect() (err error)
	EnableLogging() (err error)
	TestConnection() (err error)
	MigrateUsingWithNamespace(namespace string, models ...interface{}) (err error)
	MigrateUsing(models ...interface{}) (err error)
	Migrate(namespace string, model interface{}) (tableName string, err error)
}
