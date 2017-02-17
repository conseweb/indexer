package indexer

import (
	"encoding/json"
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
	testData(orm)
	return Orm, nil
}

func testData(o *xorm.Engine) {
	d := []interface{}{
		&Device{"aaaaaaaa", "127.0.0.1:1234"},
		&Device{"bbbbbbbb", "127.0.0.2:1234"},
		&Device{"cccccccc", "127.0.0.3:1234"},

		&FileInfo{
			DeviceID: "aaaaaaaa",
			Path:     "/a",
			Hash:     "1234",
			Size:     1234,
		},
		&FileInfo{
			DeviceID: "aaaaaaaa",
			Path:     "/b",
			Hash:     "2345",
			Size:     1234,
		},
		&FileInfo{
			DeviceID: "aaaaaaaa",
			Path:     "/c",
			Hash:     "3456",
			Size:     1234,
		},
	}

	bs, _ := json.Marshal(d)
	logrus.Infof("body: %s", string(bs))

	n, err := o.Insert(d...)
	logrus.Infof("insert %v, %v", n, err)
}
