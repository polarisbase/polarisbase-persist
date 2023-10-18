package persist

type Persist interface {
	Connect() (err error)
	Disconnect() (err error)
	EnableLogging() (err error)
	TestConnection() (err error)
	MigrateUsing(models ...interface{}) (err error)
}
