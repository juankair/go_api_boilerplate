package main

import (
	"fmt"
	"github.com/juankair/go_api_boilerplate/pkg/log"
	"net/http"
	"os"

	"github.com/uptrace/bunrouter"
)

func main() {
	router := bunrouter.New()
	logger := log.New()

	router.GET("/", func(w http.ResponseWriter, req bunrouter.Request) error {
		fmt.Println(req.Method, req.Route(), req.Params().Map())
		return nil
	})

	logger.Info("Server Is Running At https://localhost:8181")
	if err := http.ListenAndServe(":8181", router); err != nil && err != http.ErrServerClosed {
		logger.Error(err)
		os.Exit(-1)
	}

}
