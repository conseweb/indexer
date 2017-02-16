package indexer

import (
	log "github.com/Sirupsen/logrus"
	"github.com/go-xorm/xorm"
)

type MergeResult struct {
	Add    []*FileInfo `json:"add"`
	Update []*FileInfo `json:"update"`
	Del    []*FileInfo `json:"del"`
}

func (mr *MergeResult) UpdateData(orm *xorm.Engine) error {
	for _, v := range mr.Del {
		if err := v.Remove(orm); err != nil {
			log.Errorf("Remove FileInfo: %v, error: %s", v, err)
			return err
		}
	}

	for _, v := range mr.Update {
		if err := v.Update(orm); err != nil {
			log.Errorf("Update FileInfo: %v, error: %s", v, err)
			return err
		}
	}

	for _, v := range mr.Add {
		if err := v.Insert(orm); err != nil {
			log.Errorf("Add FileInfo: %v, error: %s", v, err)
			return err
		}
	}

	return nil
}
