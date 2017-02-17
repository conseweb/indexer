package api

import (
	"encoding/json"
	"fmt"
	"github.com/Sirupsen/logrus"
	"strconv"

	"github.com/conseweb/indexer"
	"github.com/go-martini/martini"
	"github.com/go-xorm/xorm"
)

type FileWrapper struct {
	*indexer.FileInfo `json:",inline"`
	Addr              string `json:"address"`
}

// GET /indexer/address/:file_id
func GetFileAddr(ctx *RequestContext, orm *xorm.Engine, params martini.Params) {
	fileid, err := strconv.Atoi(params["file_id"])
	if err != nil {
		ctx.Error(400, "invalid params file_id.")
		return
	}

	files := make([]*indexer.FileInfo, 0)
	err = orm.Where("id = ?", fileid).Find(&files)
	if err != nil {
		ctx.Error(500, err)
		return
	}

	if len(files) == 0 {
		ctx.Error(404, "not found")
		return
	}

	file := files[0]
	devs := make([]*indexer.Device, 0)
	err = orm.Where("id = ?", file.DeviceID).Find(&devs)
	if err != nil {
		ctx.Error(500, err)
		return
	}

	if len(files) == 0 {
		ctx.Error(404, "not found running server.")
		return
	}

	ctx.rnd.JSON(200, FileWrapper{file, devs[0].Address})
}

// GetDeviceIndexer GET /indexer/devices/:device_id
func GetDeviceIndexer(ctx *RequestContext, orm *xorm.Engine, params martini.Params) {
	devID := params["device_id"]
	files := []*indexer.FileInfo{}
	err := orm.Where("device_id", devID).Find(&files)
	if err != nil {
		ctx.Error(500, err)
		return
	}

	ctx.rnd.JSON(200, files)
}

// SetFileIndex POST /indexer/devices/:device_id?clean=false clean old files in this deviceID
func SetFileIndex(ctx *RequestContext, orm *xorm.Engine, params martini.Params) {
	devID := params["device_id"]
	isClean, _ := strconv.ParseBool("clean")

	var files []*indexer.FileInfo
	err := json.NewDecoder(ctx.req.Body).Decode(&files)
	if err != nil {
		ctx.Error(400, err)
		return
	}

	if isClean {
		_, err := orm.Where("device_id = ?", devID).Delete(&indexer.FileInfo{})
		if err != nil {
			ctx.Error(500, err)
			return
		}
	} else {
		filePaths := []interface{}{}
		for _, file := range files {
			filePaths = append(filePaths, file.Path)
		}
		orm.Where("device_id = ?", devID).In("path", filePaths).Delete(&indexer.FileInfo{})
	}

	insrt := []interface{}{}
	for _, file := range files {
		file.DeviceID = devID
		insrt = append(insrt, file)
	}

	n, err := orm.Insert(insrt...)
	if err != nil {
		ctx.Error(500, err)
		return
	}

	ctx.Message(201, n)
}

func UpdateFileIndex(ctx *RequestContext, orm *xorm.Engine, params martini.Params) {
	body := indexer.MergeResult{}
	// file := &indexer.FileInfo{}
	err := json.NewDecoder(ctx.req.Body).Decode(&body)
	if err != nil {
		ctx.Error(400, err)
		return
	}
	err = body.UpdateData(orm)
	if err != nil {
		ctx.Error(500, err)
		return
	}

	ctx.Message(200, "ok")
}

func DeleteFileIndex(ctx *RequestContext, orm *xorm.Engine, params martini.Params) {
	fID, err := strconv.Atoi(params["file_id"])
	if err != nil {
		ctx.Error(400, "invalid params file_id")
		return
	}

	n, err := orm.Where("file_id = ?", int64(fID)).Delete(&indexer.FileInfo{})
	if err != nil {
		ctx.Error(500, err)
		return
	}

	if n == 0 {
		ctx.Error(404, fmt.Errorf("not found file %v", fID))
		return
	}

	ctx.Message(200, "ok")
}

func OnlineDevice(ctx *RequestContext, orm *xorm.Engine, params martini.Params) {
	devID := params["device_id"]

	var dev indexer.Device

	err := json.NewDecoder(ctx.req.Body).Decode(&dev)
	if err != nil {
		ctx.Error(400, "invalid address")
		return
	}

	_, err = orm.Where("device_id = ?", devID).Delete(&indexer.Device{})
	if err != nil {
		ctx.Error(500, err)
		return
	}

	dev.ID = devID

	_, err = orm.Insert(dev)
	if err != nil {
		ctx.Error(500, err)
		return
	}

	ctx.Message(201, "ok")
}

func OfflineDevice(ctx *RequestContext, orm *xorm.Engine, params martini.Params) {
	devID := params["device_id"]

	_, err := orm.Where("device_id = ?", devID).Delete(&indexer.Device{})
	if err != nil {
		ctx.Error(500, err)
		return
	}

	ctx.Message(200, "ok")
}

func ListView(ctx *RequestContext, orm *xorm.Engine, params martini.Params) {
	files := []*indexer.FileInfo{}
	err := orm.OrderBy("updated").Find(&files)
	if err != nil {
		ctx.Error(500, err)
		return
	}
	for _, file := range files {
		logrus.Infof("index; %+v", file)
	}

	ctx.rnd.HTML(200, "index", map[string][]*indexer.FileInfo{"Files": files})
}
