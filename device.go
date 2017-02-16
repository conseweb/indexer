package indexer

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-xorm/xorm"
	"github.com/spf13/viper"
)

func InitDB() (*xorm.Engine, error) {
	if orm != nil {
		if err := orm.Ping(); err != nil {
			return nil, err
		}
		return orm, nil
	}

	path := viper.GetString("peer.fileSystemPath")
	fi, err := os.Stat(filepath.Join(path, "indexer.db"))
	var Orm *xorm.Engine
	if err != nil && os.IsExist(err) {
		return nil, err
	}
	if fi.IsDir() {
		return nil, fmt.Errorf("%s is a directory", path)
	}

	Orm, err = xorm.NewEngine("sqlite3", path)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	Orm.ShowSQL(true)

	Orm.Sync2(&FileInfo{}, &Device{})

	orm = Orm

	return Orm, nil
}
