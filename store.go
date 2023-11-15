package persist

import "github.com/upper/db/v4"

type Store interface {
	Connect() (err error)
	Collections() (collections []db.Collection, err error)
	Collection(name string, model interface{}) (collection db.Collection, err error)
	NewBucket(name string) (bucket Bucket, err error)
}

type Bucket interface {
	Collection(name string, model interface{}) (collection db.Collection, err error)
}
