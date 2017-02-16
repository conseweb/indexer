package api

import (
	"github.com/conseweb/indexer"
	"github.com/go-martini/martini"
)

func SetIndexerDBMW(ctx *RequestContext, mc martini.Context) {
	orm, err := indexer.InitDB()
	if err != nil {
		ctx.Error(500, err)
		return
	}
	mc.Map(orm)
}
