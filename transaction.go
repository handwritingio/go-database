package database

import "fmt"

// Transacter interface allows us to specify that we don't care if it's a DB or a Tx
type Transacter interface {
	Transaction(func(*Tx) error) error
}

var _ Transacter = DB{} // Interface check
var _ Transacter = Tx{} // Interface check

// Transaction runs your function in a new transaction
func (db DB) Transaction(f func(*Tx) error) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			switch p := p.(type) {
			case error:
				err = p
			default:
				err = fmt.Errorf("%s", p)
			}
		}

		if err != nil {
			tx.Rollback()
			return
		}

		err = tx.Commit()
	}()

	err = tx.Transaction(f)
	return
}

// Transaction runs your function in the current transaction
// allowing you to ignore the difference between a DB and a Tx
func (tx Tx) Transaction(f func(*Tx) error) error {
	return f(&tx)
}
