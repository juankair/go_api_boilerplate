package main

import (
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/uptrace/bunrouter"
)

var log = logrus.New()

func main() {
	router := bunrouter.New()

	router.GET("/", func(w http.ResponseWriter, req bunrouter.Request) error {
		fmt.Println(req.Method, req.Route(), req.Params().Map())
		return nil
	})

	log.WithFields(logrus.Fields{
		"param":  1231,
		"param2": "data",
	}).Info("This is log info")
	log.WithFields(logrus.Fields{
		"paramDebug":  1231,
		"paramDebug2": "data",
	}).Debug("This is log DEBUG")
	log.Warn("This is log warning")
	log.Error("This is log error")
	http.ListenAndServe(":8181", router)

}
