package common

import (
	"database/sql"
	"github.com/go-openapi/strfmt"
)

// SqlNullStringToPtr helper function converts sql.NullString to *string
func SqlNullStringToPtr(value sql.NullString) *string {
	if value.Valid {
		return &value.String
	}
	return nil
}

// SqlNullTimeToPtr helper function converts sql.NullTime to *strfmt.DateTime
func SqlNullTimeToPtr(value sql.NullTime) *strfmt.DateTime {
	if value.Valid {
		v := strfmt.DateTime(value.Time)
		return &v
	}
	return nil
}
