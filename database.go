package database

import (
	"database/sql"
	"log"
	"os"

	"github.com/lib/pq"
)

var (
	databaseURL = os.Getenv("DATABASE_URL")
	dbPool      *DB
)

// GetDB returns a db connection, reusing the same connection if one already exists.
func GetDB() *DB {
	if dbPool == nil {
		dbPool = mustConnect(databaseURL)
	}
	return dbPool
}

// DB wraps sql.DB
// FIXME Pointer type nesting is messy here...
type DB struct {
	*sql.DB
}

// Tx wraps sql.Tx
type Tx struct {
	*sql.Tx
}

// Begin starts a transaction
func (db *DB) Begin() (*Tx, error) {
	tx, err := db.DB.Begin()
	if err != nil {
		return nil, err
	}
	return &Tx{tx}, nil
}

// Connects to the database and pings it to ensure we're connected
func mustConnect(url string) *DB {
	db, err := sql.Open("postgres", url)
	Check(err, "Connection URL invalid")

	// NOTE: Open() doesn't actually connect, it only validates the connection
	// string and prepares a pool, so we Ping here to actually attempt
	// connection
	err = db.Ping()
	Check(err, "Could not connect to database")
	return &DB{db}
}

// Check is shorthand for dying on errors, but see if we can cast these errors to
// pq.Errors for more info
func Check(err error, extraInfo ...string) {
	if err, ok := err.(*pq.Error); ok {
		log.Printf("PQERROR %s (%s): %s\n", err.Code, err.Code.Name(), err.Message)
		log.Println(err.Severity)
		log.Println(err.Code)
		log.Println(err.Message)
		log.Println(err.Detail)
		log.Println(err.Hint)
		log.Println(err.Position)
		log.Println(err.InternalPosition)
		log.Println(err.InternalQuery)
		log.Println(err.Where)
		log.Println(err.Schema)
		log.Println(err.Table)
		log.Println(err.Column)
		log.Println(err.DataTypeName)
		log.Println(err.Constraint)
		log.Println(err.File)
		log.Println(err.Line)
		log.Println(err.Routine)
	}
	if err != nil {
		log.Panic(extraInfo, err)
	}
}
