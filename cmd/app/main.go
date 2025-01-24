package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"

	dbx "github.com/go-ozzo/ozzo-dbx"
	_ "github.com/go-sql-driver/mysql"
	"github.com/juankair/docs_sign_be/internal/account"
	"github.com/juankair/docs_sign_be/internal/auth"
	"github.com/juankair/docs_sign_be/internal/public"
	"github.com/juankair/docs_sign_be/pkg/dbcontext"

	"github.com/juankair/docs_sign_be/internal/config"
	"github.com/juankair/docs_sign_be/pkg/log"

	"github.com/uptrace/bunrouter"
)

var flagConfig = flag.String("config", "./config/local.yml", "path to the config file")

type CORSHandler struct {
	Next         http.Handler
	HostFrontend string
}

func main() {
	flag.Parse()

	logger := log.New()

	cfg, err := config.Load(*flagConfig)
	if err != nil {
		logger.Error("failed to load application configuration: %s", err)
		os.Exit(-1)
	}

	db, err := dbx.MustOpen("mysql", cfg.DSN)
	if err != nil {
		logger.Error(err)
		os.Exit(-1)
	}

	defer func() {
		if err := db.Close(); err != nil {
			logger.Error(err)
		}
	}()

	address := fmt.Sprintf(":%v", cfg.ServerPort)
	configServer := &http.Server{
		Addr:    address,
		Handler: buildHandler(logger, dbcontext.New(db), cfg),
	}

	logger.Info(fmt.Sprintf("Server Is Running At https://localhost:%s", strconv.Itoa(cfg.ServerPort)))
	if err := configServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error(err)
		os.Exit(-1)
	}

}

func buildHandler(logger log.Logger, db *dbcontext.DB, cfg *config.Config) http.Handler {

	router := bunrouter.New()

	public.RegisterHandler(router)

	auth.RegisterHandler(router,
		auth.NewService(auth.NewRepository(db, logger), cfg.JWTSigningKey, cfg.JWTExpiration, logger),
		logger,
	)

	router.Use(
		auth.SecureMiddleware(cfg.JWTSigningKey),
	).WithGroup("/app", func(secure *bunrouter.Group) {
		account.RegisterHandler(secure,
			account.NewService(account.NewRepository(db, logger), logger),
			logger,
		)
	})

	handler := NewCORSHandler(router, cfg.HostFrontend)

	return handler
}

func NewCORSHandler(next http.Handler, hostFrontend string) CORSHandler {
	return CORSHandler{
		Next:         next,
		HostFrontend: hostFrontend,
	}
}

func (h CORSHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	allowedOrigins := h.HostFrontend

	header := w.Header()
	header.Set("Access-Control-Allow-Origin", allowedOrigins)
	header.Set("Access-Control-Allow-Credentials", "true")

	if req.Method == http.MethodOptions {
		header.Set("Access-Control-Allow-Methods", "GET,PUT,POST,DELETE")
		header.Set("Access-Control-Allow-Headers", "authorization,content-type")
		header.Set("Access-Control-Max-Age", "86400")

		return
	}

	h.Next.ServeHTTP(w, req)
}
