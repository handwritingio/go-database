package database

// Scanner is a stand-in for sql.Row or sql.Rows. It allows for writing model
// methods that handle the call to Scan to populate the struct.
type Scanner interface {
	Scan(dest ...interface{}) error
}
