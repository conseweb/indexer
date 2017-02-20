package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/conseweb/indexer/api"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	err := api.Serve(":8080")
	if err != nil {
		logrus.Error(err)
		return
	}
}
