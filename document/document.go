package document

import (
	"github.com/polarisbase/polaris-sdk/v3/lib/persist"
	"github.com/uptrace/bun"
)

type Store interface {
	persist.Persist
	GetBun() *bun.DB
}
