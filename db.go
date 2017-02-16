package indexer

import (
	"database/sql"

	log "github.com/Sirupsen/logrus"
)

var (
	stddb *sql.DB
)

func GetDB() *sql.DB {
	if stddb != nil {
		return stddb
	}

	path := "/tmp/farmer.db"
	var err error
	stddb, err = loaddb(path)
	if err != nil {
		log.Errorf("GetDB: %s", err.Error())
		return nil
	}

	return stddb
}

func loaddb(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
