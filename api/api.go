package api

import (
	"database/sql"
	stdlog "log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/cors"
	"github.com/martini-contrib/render"
	"github.com/spf13/viper"
)

const (
	API_PREFIX = "/api"
)

type RequestContext struct {
	params martini.Params
	mc     martini.Context

	req *http.Request
	rnd render.Render
	res http.ResponseWriter

	db *sql.DB
}

func notFound(gateways map[string]*url.URL) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		for way, _ := range gateways {
			if strings.HasPrefix(req.URL.Path, way) {
				w.WriteHeader(http.StatusNotFound)
				return
			}
		}
		for _, v := range []string{API_PREFIX} {
			if strings.HasPrefix(req.URL.Path, v) {
				w.WriteHeader(http.StatusNotFound)
				return
			}
		}
	}
}

func getGatewayRouters() (map[string]*url.URL, error) {
	gateways := map[string]*url.URL{}

	for k, v := range viper.GetStringMapString("daemon.gateway") {
		log.Infof("get setting proxy router, %s --> %s", k, v)
		way := "/" + strings.Trim(k, "/ \n")
		to, err := url.Parse(v)
		if err != nil {
			log.Errorf("formt URL<%s> failed, error: %s", v, err.Error())
			return nil, err
		}
		gateways[way] = to
	}

	return gateways, nil
}

func Serve(listenAddr string) error {
	m := NewMartini()

	m.Use(cors.Allow(&cors.Options{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "PATCH", "POST", "DELETE", "PUT"},
		AllowHeaders:     []string{"Limt", "Offset", "Content-Type", "Origin", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Record-Count", "Limt", "Offset", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           time.Second * 864000,
	}))

	m.Use(requextCtx)
	m.Get("/", SetIndexerDBMW, ListView)
	m.Group(API_PREFIX, func(r martini.Router) {
		r.Group("/indexer", func(r martini.Router) {
			r.Get("/devices/:device_id", GetDeviceIndexer)
			r.Post("/devices/:device_id/online", OnlineDevice)
			r.Post("/devices/:device_id/offline", OfflineDevice)
			r.Post("/devices/:device_id/files", SetFileIndex)

			r.Post("/files", UpdateFileIndex)
			r.Get("/files/:file_id", GetFileAddr)
			r.Delete("/files/:file_id", DeleteFileIndex)
		}, SetIndexerDBMW)
	})

	server := &http.Server{
		Handler:  m,
		Addr:     listenAddr,
		ErrorLog: stdlog.New(os.Stderr, "", 0),
	}

	log.Info("server is starting on ", listenAddr)
	return server.ListenAndServe()
}

func NewMartini() *martini.ClassicMartini {
	r := martini.NewRouter()
	m := martini.New()
	m.Use(martini.Recovery())
	m.Use(render.Renderer(render.Options{
		Directory: "api/templates",
		// Layout:          "layout",
		Extensions:      []string{".tmpl", ".html"},
		Delims:          render.Delims{"{{", "}}"},
		Charset:         "UTF-8",
		HTMLContentType: "text/html",
		// IndentJSON:      true,
		// IndentXML:       true,
		// Funcs:           []template.FuncMap{AppHelpers},
	}))
	m.MapTo(r, (*martini.Routes)(nil))
	m.Action(r.Handle)
	return &martini.ClassicMartini{Martini: m, Router: r}
}

func requextCtx(w http.ResponseWriter, req *http.Request, mc martini.Context, rnd render.Render) {
	ctx := &RequestContext{
		res:    w,
		req:    req,
		mc:     mc,
		rnd:    rnd,
		params: make(map[string]string),
	}

	req.ParseForm()
	if len(req.Form) > 0 {
		for k, v := range req.Form {
			ctx.params[k] = v[0]
		}
	}

	log.Debugf("[%s] %s", req.Method, req.URL.String())

	mc.Map(ctx)
	mc.Next()
}
