package database

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"regexp"
	"strconv"
)

var (
	// our migration files must start with 3 integers followed by an underscore
	reMigrationPrefix = regexp.MustCompile(`^[0-9]{3}_.+`)
)

// Just a helper struct to represent each migration file we find
type migrationFile struct {
	AbsolutePath string
	Filename     string
	Version      int
}

// Scan our migrations folder for files that look like migrations and return
// migrationFile structs that represent them in proper order
func findPotentialMigrations(migrationDir string) []migrationFile {
	out := []migrationFile{}
	files, err := ioutil.ReadDir(migrationDir)
	Check(err, "Listing migration directory")
	// we'll use this to determine if we skipped a number or double defined one
	dedupeCounter := 0
	for _, f := range files { // these should be sorted by the ReadDir call above
		// find entries that look like our migrations and are just regular files
		if f.Mode().IsRegular() && reMigrationPrefix.MatchString(f.Name()) {
			version, err := strconv.Atoi(f.Name()[0:3])
			switch {
			case version == dedupeCounter:
				log.Fatal("FATAL: duplicate version number! ", version)
			case version != dedupeCounter+1:
				log.Fatal("FATAL: skipped a version after migration! ", dedupeCounter)
			default:
				dedupeCounter = version
			}
			Check(err, "converting", f.Name(), "to version integer")
			mf := migrationFile{
				AbsolutePath: filepath.Join(migrationDir, f.Name()),
				Filename:     f.Name(),
				Version:      version}
			out = append(out, mf)
		}
	}
	return out
}

// GetCurrentVersion returns the current schema version as an integer, or nil with an error
func (db *DB) GetCurrentVersion() int {
	currentVersion := 0
	db.DB.QueryRow(`
		SELECT max(current_version) FROM meta.migrations
	`).Scan(&currentVersion)
	return currentVersion
}

func (db *DB) isMigrationsTablePresent() (exists bool) {
	err := db.DB.QueryRow(`
		SELECT EXISTS (
			SELECT *
			FROM information_schema.tables
			WHERE
				table_schema = 'meta'
				AND
				table_name = 'migrations'
		)
	`).Scan(&exists)
	Check(err, "looking for migrations table")
	return
}

func (db *DB) bootstrap(migrationDir string) {
	if db.isMigrationsTablePresent() {
		return
	}
	log.Println("creating migrations table")
	// find the bootstrap.sql file and run it
	bs, err := ioutil.ReadFile(filepath.Join(migrationDir, "bootstrap.sql"))
	Check(err, "reading contents of bootstrap.sql")
	sql := string(bs)
	_, err = db.DB.Exec(sql)
	Check(err, "bootstrapping schema")
}

// Migrate ensures the db supports migration, loads bootstrap info, and applies
// as many migrations as it can in order until there are no migrations left to
// run, or a migration produces an error.
func (db *DB) Migrate(migrationDir string, dryrun bool) {
	if dryrun {
		log.Print("migrating (dry run)...")
	} else {
		log.Print("migrating...")
	}

	exists := db.isMigrationsTablePresent()
	currentVersion := 0 // assume 0 unless we find out otherwise

	if exists {
		log.Println("migrations table found")
		currentVersion = db.GetCurrentVersion()
	} else {
		if dryrun {
			log.Println("would bootstrap database")
		} else {
			db.bootstrap(migrationDir)
		}
	}

	migrations := findPotentialMigrations(migrationDir)
	count := 0
	for _, mf := range migrations {
		if mf.Version > currentVersion {
			if dryrun {
				log.Println("would apply migration", mf.Filename)
			} else {
				db.mustApplyMigration(mf)
			}
			count++
		}
	}
	if dryrun {
		log.Println("would apply", count, "migrations")
	} else {
		log.Println("applied", count, "migrations")
	}
}

// Apply a numbered migration to the database in a transaction
func (db *DB) mustApplyMigration(mf migrationFile) {
	log.Println("applying migration:", mf.Filename)

	bs, err := ioutil.ReadFile(mf.AbsolutePath)
	Check(err, "reading contents of", mf.AbsolutePath)
	sql := string(bs)

	tx, err := db.Begin()
	_, err = tx.Exec(sql) // the actual migration
	if err != nil {
		tx.Rollback()
		Check(err, fmt.Sprintf("during migration %d", mf.Version))
	}
	_, err = tx.Exec(`
		INSERT INTO
			meta.migrations (
				current_version,filename
			) VALUES (
				$1, $2
			);
		`, mf.Version, mf.Filename)
	if err != nil {
		tx.Rollback()
		Check(err, fmt.Sprintf("updating migration log %d", mf.Version))
	}
	tx.Commit()
}
