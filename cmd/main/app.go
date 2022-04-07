package main

import (
	"fmt"
	"net"
	"net/http"
	"rest-api-test/internal/handlers/user"
	"rest-api-test/pkg/logging"
	"time"

	"github.com/julienschmidt/httprouter"
)

func IndexHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	name := params.ByName("name")
	w.Write([]byte(fmt.Sprintf("Hello %s", name)))
}

func main() {
	logger := logging.GetLogger()
	logger.Info("create router")
	// создаем роутер
	router := httprouter.New()

	logger.Info("register user handler")
	// создаем handler
	handler := user.NewHandler(logger)
	// регистрируем handler в router
	handler.Register(router)

	start(router)
}

// стартует сервер на порту 1234, по протоколу tsp
func start(router *httprouter.Router) {
	logger := logging.GetLogger()
	logger.Info("start application")

	listener, err := net.Listen("tcp", "0.0.0.0:1234")
	if err != nil {
		panic(err)
	}

	server := &http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logger.Info("server is listening port 1234")
	logger.Fatal(server.Serve(listener))
}
