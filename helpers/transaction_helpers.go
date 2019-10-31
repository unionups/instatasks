package helpers

import (
	"github.com/jinzhu/gorm"
	"sync"
)

type Transaction struct {
	once     sync.Once
	rollback bool
	tx       *gorm.DB
}

func (t *Transaction) Close() {
	t.once.Do(func() {
		if t.rollback {
			t.tx.Rollback()
		} else {
			t.tx.Commit()
		}
	})
}

func (t *Transaction) Fail() {
	t.rollback = true
}

func NewTransaction(db *gorm.DB) (*gorm.DB, *Transaction) {
	return db, &Transaction{tx: db}
}
