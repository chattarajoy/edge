package mysql

import (
	"database/sql"
	"fmt"
)

type DbLog struct {
	Id 		  string `json:"id" required:"true"`
	Path	  string `json:"path" required:"true"`
}

func (dblog *DbLog) Insert(db *sql.DB) error {
	insertStatement := fmt.Sprintf("INSERT INTO db_log (path) VALUES ('%s')", dblog.Path)

	insert, _ := db.Query(insertStatement)
	return insert.Close()
}