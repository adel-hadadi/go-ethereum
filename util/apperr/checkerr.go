package apperr

import (
	"github.com/mattn/go-sqlite3"
)

const ErrSQLDuplicateEntryCode = "23505"

func IsSQLDuplicateEntry(err error) bool {
	if sqliteErr, ok := err.(sqlite3.Error); ok {
		// Check if the error code is SQLITE_CONSTRAINT
		if sqliteErr.Code == sqlite3.ErrConstraint {
			return true
		}
	}

	return false
}
