package database

import (
	"fmt"
	"log"
	"os"

	"github.com/lib/pq"
)

var (
	testDatabaseURL = os.Getenv("TEST_DATABASE_URL")
	testDB          *TestDB
)

// TestDB is the same as DB, but provides clean up methods that should never be
// used in production.
type TestDB struct {
	*DB
}

var _ Transacter = TestDB{} // Interface check

// GetTestDB returns a connection to the TEST_DATABASE_URL, reusing the same
// connection if one already exists.
func GetTestDB() *TestDB {
	if testDB == nil {
		db := mustConnect(testDatabaseURL)
		testDB = &TestDB{db}
	}
	return testDB
}

// TruncateTable does horrible things which is why it's only allowed on test databases
func (db *TestDB) TruncateTable(tableName string) {
	query := fmt.Sprintf("TRUNCATE TABLE %s CASCADE", pq.QuoteIdentifier(tableName))
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}
