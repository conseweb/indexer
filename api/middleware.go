package api

import (
	"github.com/Sirupsen/logrus"
	"github.com/conseweb/indexer"
	"github.com/go-martini/martini"
)

func SetIndexerDBMW(ctx *RequestContext, mc martini.Context) {
	logrus.Infof("Path: %s", ctx.req.URL.Path)
	orm, err := indexer.GetXorm()
	if err != nil {
		ctx.Error(500, err)
		return
	}
	mc.Map(orm)
}
