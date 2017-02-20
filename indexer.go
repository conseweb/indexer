package indexer

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/go-xorm/xorm"
	"github.com/spf13/viper"
)

var (
	orm *xorm.Engine
)

type Device struct {
	ID      string `xorm:"pk 'id'" json:"id"`
	Address string `xorm:"notnull index" json:"address"`
	Online  bool   `xorm:"'online'" json:"online"`
}

type Indexer struct {
	cli    *http.Client
	chroot string
	devID  string
	addr   string
}

func NewIndexer(addr, devID, localChroot string) (*Indexer, error) {
	_, err := url.Parse(addr)
	if err != nil {
		return nil, err
	}

	if devID == "" {
		return nil, fmt.Errorf("required device ID.")
	}

	return &Indexer{
		cli:    http.DefaultClient,
		addr:   addr,
		devID:  devID,
		chroot: localChroot,
	}, nil
}

func (idx *Indexer) ListLocalAll() ([]*FileInfo, error) {
	type walkArgs struct {
		fp   string
		info os.FileInfo
		e    error
	}
	files := []*FileInfo{}

	buf := make(chan *FileInfo, 100)
	wksbuf := make(chan *FileInfo, 100)
	done := make(chan struct{})

	MaxRunc := 4
	wg := new(sync.WaitGroup)

	for i := 0; i < MaxRunc; i++ {
		go func() {
			for {
				select {
				case file := <-wksbuf:
					f, err := os.Open(file.Path)
					if err != nil {
						log.Error(err)
						wg.Done()
						return
					}

					hash := sha256.New()
					_, err = io.Copy(hash, f)
					if err != nil {
						log.Error(err)
						wg.Done()
						f.Close()
						return
					}

					file.Hash = fmt.Sprintf("%x", hash.Sum(nil))
					f.Close()
					buf <- file

				case <-done:
					return
				}
			}
		}()
	}

	go func() {
		for {
			select {
			case fi := <-buf:
				wg.Done()
				files = append(files, fi)
			case <-done:
				return
			}
		}
	}()

	filepath.Walk(idx.chroot, func(fp string, info os.FileInfo, e error) error {
		if e != nil {
			log.Error(e)
			return e
		}
		if info.IsDir() {
			return nil
		}

		wg.Add(1)
		wksbuf <- &FileInfo{
			DeviceID: idx.devID,
			Path:     fp,
			Size:     info.Size(),
		}
		return nil
	})

	wg.Wait()
	close(done)

	return files, nil
}

func (idx *Indexer) ListRemoteAll() ([]*FileInfo, error) {
	u, _ := url.Parse(idx.addr)
	u.Path = fmt.Sprintf("/devices/%s/", idx.devID)

	resp, err := idx.cli.Get(u.String())
	if err != nil {
		return nil, err
	}

	if resp.StatusCode > 300 {
		bs, _ := ioutil.ReadAll(resp.Body)
		log.Errorf("Status: %v, body: %s", resp.Status, bs)
		return nil, fmt.Errorf("%s", bs)
	}

	remoteFiles := []*FileInfo{}
	err = json.NewDecoder(resp.Body).Decode(remoteFiles)
	if err != nil {
		return nil, err
	}

	return remoteFiles, nil
}

func (idx *Indexer) send(files []*FileInfo) error {
	u, _ := url.Parse(idx.addr)
	u.Path = fmt.Sprintf("/devices/%s/online", idx.devID)

	body := &bytes.Buffer{}
	err := json.NewEncoder(body).Encode(files)
	if err != nil {
		return err
	}
	resp, err := idx.cli.Post(u.String(), "application/json", body)
	if err != nil {
		return err
	}

	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	log.Infof("Send over. %s", bs)
	return nil
}

func (idx *Indexer) SyncLocalFS() error {
	remoteFiles, err := idx.ListRemoteAll()
	if err != nil {
		return err
	}

	localFiles, err := idx.ListLocalAll()
	if err != nil {
		return err
	}

	ret, err := idx.mergeFiles(localFiles, remoteFiles)
	if err != nil {
		return err
	}

	for _, fi := range ret.Add {
		err = fi.Download(idx)
		if err != nil {
			return fmt.Errorf("SyncLocalFS Download %s, %s", fi.Path, err)
		}
	}

	for _, fi := range ret.Update {
		err = fi.Download(idx)
		if err != nil {
			return fmt.Errorf("SyncLocalFS Update %s, %s", fi.Path, err)
		}
	}

	for _, fi := range ret.Del {
		err = fi.RemoveLocal(idx)
		if err != nil {
			return fmt.Errorf("SyncLocalFS Del %s, %s", fi.Path, err)
		}
	}

	return nil
}

func (idx *Indexer) Online() error {
	addr := viper.GetString("daemon.address")
	dev := &Device{
		Address: addr,
	}

	u, _ := url.Parse(idx.addr)
	u.Path = fmt.Sprintf("/devices/%s/online", idx.devID)

	buf := &bytes.Buffer{}
	err := json.NewEncoder(buf).Encode(dev)
	if err != nil {
		return err
	}

	resp, err := idx.cli.Post(u.String(), "application/json", buf)
	if err != nil {
		return err
	}

	if resp.StatusCode >= 300 {
		log.Errorf("status: %s", resp.Status)
		return fmt.Errorf("online failed.")
	}

	return nil
}

func (idx *Indexer) SyncToRemote() error {
	remoteFiles, err := idx.ListRemoteAll()
	if err != nil {
		return err
	}

	localFiles, err := idx.ListLocalAll()
	if err != nil {
		return err
	}

	ret, err := idx.mergeFiles(remoteFiles, localFiles)
	if err != nil {
		return err
	}

	err = idx.updateRemote(ret)
	if err != nil {
		return err
	}

	return nil
}

func (idx *Indexer) mergeFiles(bases, currents []*FileInfo) (*MergeResult, error) {
	ret := &MergeResult{
		Add:    []*FileInfo{},
		Update: []*FileInfo{},
		Del:    []*FileInfo{},
	}

	for i := 0; i < len(currents); i++ {
		local := currents[i]
		local.DeviceID = idx.devID
		exists := false
		var j int
	loop_remote:
		for ; j < len(bases); j++ {
			remote := bases[j]
			if remote.Hash == local.Hash && remote.Path == local.Path {
				exists = true
				break loop_remote
			} else if (remote.Hash != local.Hash && remote.Path == local.Path) ||
				(remote.Hash == local.Hash && remote.Path != local.Path) {
				exists = true
				local.ID = remote.ID
				ret.Update = append(ret.Update, local)
				break loop_remote
			}
		}
		if exists {
			if j+1 < len(bases) {
				bases = append(bases[:j], bases[j+1:]...)
			} else {
				bases = bases[:j]
			}
		} else {
			ret.Add = append(ret.Add, local)
		}
	}
	for _, rem := range bases {
		rem.DeviceID = idx.devID
		ret.Del = append(ret.Del, rem)
	}
	return ret, nil
}

func (idx *Indexer) updateRemote(upd *MergeResult) error {
	buf := &bytes.Buffer{}
	err := json.NewEncoder(buf).Encode(upd)
	if err != nil {
		return err
	}

	u, _ := url.Parse(idx.addr)
	u.Path = "/devices/files"

	resp, err := idx.cli.Post(u.String(), "application/json", buf)
	if err != nil {
		return err
	}

	if resp.StatusCode >= 300 {
		return fmt.Errorf("%s", resp.Status)
	}

	return nil
}
