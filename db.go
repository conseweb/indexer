package indexer

import (
	"path/filepath"

	"github.com/Sirupsen/logrus"
	"github.com/go-xorm/xorm"
	_ "github.com/mattn/go-sqlite3"
)

// func GetDB() *sql.DB {
// 	if stddb != nil {
// 		return stddb
// 	}

// 	path := "/tmp/farmer.db"
// 	var err error
// 	stddb, err = loaddb(path)
// 	if err != nil {
// 		log.Errorf("GetDB: %s", err.Error())
// 		return nil
// 	}

// 	return stddb
// }

// func loaddb(path string) (*sql.DB, error) {
// 	db, err := sql.Open("sqlite3", path)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if err := db.Ping(); err != nil {
// 		return nil, err
// 	}

// 	return db, nil
// }

func GetXorm() (*xorm.Engine, error) {
	if orm != nil {
		return orm, nil
	}

	return InitDB()
}

func InitDB() (*xorm.Engine, error) {
	if orm != nil {
		if err := orm.Ping(); err != nil {
			return nil, err
		}
		return orm, nil
	}

	path := filepath.Join("/tmp", "indexer.db")

	Orm, err := xorm.NewEngine("sqlite3", path)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	Orm.ShowSQL(true)

	Orm.Sync2(&FileInfo{}, &Device{})

	orm = Orm

	return Orm, nil
}
