package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"time"
	"youtube_lesson/internal/handlers/user"

	"github.com/julienschmidt/httprouter"
)

func IndexHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	name := params.ByName("name")
	w.Write([]byte(fmt.Sprintf("Hello %s", name)))
}

func main() {
	log.Println("create router")
	// создаем роутер
	router := httprouter.New()

	log.Println("register user handler")
	// создаем handler
	handler := user.NewHandler()
	// регистрируем handler в router
	handler.Register(router)

	start(router)
}

// стартует сервер на порту 1234, по протоколу tsp
func start(router *httprouter.Router) {
	log.Println("start application")

	listener, err := net.Listen("tcp", "0.0.0.0:1234")
	if err != nil {
		panic(err)
	}

	server := &http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println("server is listening port 1234")
	log.Fatalln(server.Serve(listener))
}
