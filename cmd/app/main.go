package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/juankair/go_api_boilerplate/internal/config"
	"github.com/juankair/go_api_boilerplate/pkg/log"

	"github.com/uptrace/bunrouter"
)

var flagConfig = flag.String("config", "./config/local.yml", "path to the config file")

func main() {
	flag.Parse()

	logger := log.New()

	cfg, err := config.Load(*flagConfig, logger)
	if err != nil {
		logger.Error("failed to load application configuration: %s", err)
		os.Exit(-1)
	}

	address := fmt.Sprintf(":%v", cfg.ServerPort)
	configServer := &http.Server{
		Addr:    address,
		Handler: buildHandler(logger),
	}

	logger.Info(fmt.Sprintf("Server Is Running At https://localhost:%s", strconv.Itoa(cfg.ServerPort)))
	if err := configServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error(err)
		os.Exit(-1)
	}

}

func buildHandler(logger log.Logger) http.Handler {
	router := bunrouter.New()
	router.GET("/", func(w http.ResponseWriter, req bunrouter.Request) error {
		fmt.Println(req.Method, req.Route(), req.Params().Map())
		return nil
	})

	return router
}
