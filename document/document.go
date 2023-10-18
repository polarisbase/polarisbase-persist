package document

import (
	"github.com/polarisbase/polarisbase-persist"
	"github.com/uptrace/bun"
)

type Store interface {
	persist.Persist
	GetBun() *bun.DB
}
