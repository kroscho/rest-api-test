package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"rest-api-test/internal/author"
	authorRepo "rest-api-test/internal/author/db"
	"rest-api-test/internal/config"
	"rest-api-test/pkg/client/postgresql"
	"rest-api-test/pkg/logging"
	"time"

	"github.com/julienschmidt/httprouter"
)

func main() {
	logger := logging.GetLogger()
	logger.Info("create router")
	// создаем роутер
	router := httprouter.New()

	cfg := config.GetConfig()

	postgreSQLClient, err := postgresql.NewClient(context.TODO(), cfg.Storage)
	if err != nil {
		logger.Fatalf("%v", err)
	}
	repository := authorRepo.NewRepository(postgreSQLClient, logger)

	logger.Info("register author handler")
	authorHandler := author.NewHandler(repository, logger)
	authorHandler.Register(router)

	start(router, cfg)
}

// стартует сервер на порту 1234, по протоколу tsp
func start(router *httprouter.Router, cfg *config.Config) {
	logger := logging.GetLogger()
	logger.Info("start application")

	var listener net.Listener
	var listenError error

	if cfg.Listen.Type == "sock" {
		// /path/to/binary
		// Dir() -- /path/to
		logger.Info("detect app path")
		appDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			logger.Fatal(err)
		}
		logger.Info("create socket")
		socketPath := path.Join(appDir, "app.sock")

		logger.Info("listen unix socket")
		listener, listenError = net.Listen("unix", socketPath)
		logger.Infof("server is listening unix socket: %s", socketPath)
	} else {
		logger.Info("Listen tcp")
		listener, listenError = net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.Listen.BindIp, cfg.Listen.Port))
		logger.Infof("server is listening port %s:%s", cfg.Listen.BindIp, cfg.Listen.Port)
	}

	if listenError != nil {
		logger.Fatal(listenError)
	}

	server := &http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logger.Fatal(server.Serve(listener))
}
