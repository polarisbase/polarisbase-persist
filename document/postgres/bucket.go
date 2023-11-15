package postgres

import "github.com/upper/db/v4"

type Bucket struct {
	Store     *Store
	Namespace string
}

// Collection returns a collection by name.
func (b *Bucket) Collection(name string, model interface{}) (collection db.Collection, err error) {
	collection, err = b.Store.Collection(b.Namespace+"_"+name, model)
	if err != nil {
		return collection, err
	}
	return
}
