package indexer

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/go-xorm/xorm"
)

type FileInfo struct {
	ID       int64  `xorm:"pk autoincr 'id'" json:"id"`
	DeviceID string `xorm:"notnull index 'device_id'" json:"device_id"`
	Path     string `xorm:"notnull index 'path'" json:"path"`
	Hash     string `xorm:"notnull index 'hash'" json:"hash"`
	Size     int64  `xorm:"'size'" json:"hash"`

	Created time.Time `xorm:"created" json:"created"`
	Updated time.Time `xorm:"updated" json:"updated"`
}

/// For API server
func (fi *FileInfo) Insert(orm *xorm.Engine) error {
	n, err := orm.Insert(fi)
	return checkOrmRet(n, err)
}

func (fi *FileInfo) Update(orm *xorm.Engine) error {
	n, err := orm.Where("id = ?", fi.ID).Update(fi)
	return checkOrmRet(n, err)
}

func (fi *FileInfo) Remove(orm *xorm.Engine) error {
	n, err := orm.Where("id = ?", fi.ID).Delete(fi)
	return checkOrmRet(n, err)
}

/// For Client
//
func (fi *FileInfo) Download(idx *Indexer) error {
	fp := filepath.Join(idx.chroot, fi.Path)

	u, _ := url.Parse(idx.addr)
	u.Path = filepath.Join("/api/fs/cat/", fi.Path)

	resp, err := idx.cli.Get(u.String())
	if err != nil {
		return err
	}

	if resp.StatusCode >= 300 {
		return fmt.Errorf(resp.Status)
	}

	_ = os.Remove(fp)
	f, err := os.OpenFile(fp, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer func() {
		resp.Body.Close()
		f.Close()
	}()

	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func (fi *FileInfo) RemoveLocal(idx *Indexer) error {
	fp := filepath.Join(idx.chroot, fi.Path)
	return os.Remove(fp)
}

func checkOrmRet(n int64, err error) error {
	if err != nil {
		return err
	}

	if n == 0 {
		return fmt.Errorf("control xorm not found")
	}

	return nil
}
