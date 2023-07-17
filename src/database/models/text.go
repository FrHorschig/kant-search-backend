package models

import (
	"database/sql"
)

type Text struct {
	Id   int32
	Text sql.NullString
}
