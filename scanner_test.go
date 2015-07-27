package database

import (
	"errors"
	"fmt"
)

type Model struct {
	ID   int
	Name string
}

func (m *Model) scan(s Scanner) error {
	return s.Scan(&m.ID, &m.Name)
}

func ExampleScanner() {
	err := GetTestDB().Transaction(func(tx *Tx) error {
		q := `CREATE TABLE models (id integer, name text)`
		_, err := tx.Exec(q)
		if err != nil {
			return err
		}

		q = `INSERT INTO models VALUES (23, 'skidoo')`
		_, err = tx.Exec(q)
		if err != nil {
			return err
		}

		q = `SELECT id, name FROM models`
		r := tx.QueryRow(q)
		var m Model
		err = m.scan(r)
		if err != nil {
			return err
		}

		if m.ID != 23 {
			return errors.New("m.ID != 23")
		}
		if m.Name != "skidoo" {
			return errors.New("m.Name != skidoo")
		}

		return errors.New("Rolling back.")
	})
	fmt.Println(err)

	// Output:
	// Rolling back.
}
