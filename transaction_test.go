package database

import (
	"errors"
	"fmt"
)

func ExampleDB_Transaction() {
	err := GetTestDB().Transaction(func(tx *Tx) error {
		q := `CREATE TABLE bogus_transaction_test (id integer)`
		_, err := tx.Exec(q)
		if err != nil {
			return err
		}

		q = `INSERT INTO bogus_transaction_test VALUES (42)`
		_, err = tx.Exec(q)
		if err != nil {
			return err
		}

		q = `SELECT id FROM bogus_transaction_test`
		r := tx.QueryRow(q)
		var i int64
		err = r.Scan(&i)
		if err != nil {
			return err
		}

		if i != 42 {
			return errors.New("i != 42")
		}

		return errors.New("OK. I've had my fun, but I don't really want that table")
	})

	fmt.Println(err)

	err = GetTestDB().Transaction(func(tx *Tx) error {
		q := `SELECT id FROM bogus_transaction_test`
		r := tx.QueryRow(q)
		var i int64
		err = r.Scan(&i)
		return err
	})
	fmt.Println(err)

	// Output:
	// OK. I've had my fun, but I don't really want that table
	// pq: relation "bogus_transaction_test" does not exist
}

func ExampleDB_Transaction_destructiveOperationRollback() {
	// Create table with one value
	err := GetTestDB().Transaction(func(tx *Tx) error {
		q := `CREATE TABLE bogus_transaction_test (id integer)`
		_, err := tx.Exec(q)
		if err != nil {
			return err
		}

		q = `INSERT INTO bogus_transaction_test VALUES (42)`
		_, err = tx.Exec(q)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Table created")

	// Truncate table in transaction that rolls back
	err = GetTestDB().Transaction(func(tx *Tx) error {
		q := `TRUNCATE TABLE bogus_transaction_test`
		_, err = tx.Exec(q)
		if err != nil {
			fmt.Println(err)
		}
		return errors.New("Truncate rolled back")
	})
	fmt.Println(err)

	// Verify data is still there
	err = GetTestDB().Transaction(func(tx *Tx) error {
		q := `SELECT id FROM bogus_transaction_test`
		r := tx.QueryRow(q)
		var i int64
		err = r.Scan(&i)
		if err != nil {
			return err
		}

		if i != 42 {
			return errors.New("i != 42")
		}

		return nil
	})

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Data still exists in table")
	}

	// Clean up
	err = GetTestDB().Transaction(func(tx *Tx) error {
		q := `DROP TABLE bogus_transaction_test`
		_, err = tx.Exec(q)
		if err != nil {
			fmt.Println(err)
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}

	// Output:
	// Table created
	// Truncate rolled back
	// Data still exists in table
}
