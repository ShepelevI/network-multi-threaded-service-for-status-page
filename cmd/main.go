package main

import (
	"github.com/go-chi/chi"
	"log"
	"net/http"
	"os"
	"os/signal"
	"skillbox_final/internal/handler"
	"skillbox_final/internal/json"
	"syscall"
)

const gate = "http://127.0.0.1:8383"

type App struct {
	router *chi.Mux
	done   chan os.Signal
}

func NewApp() *App {
	ret := &App{
		router: chi.NewRouter(),
		done:   make(chan os.Signal, 1),
	}
	signal.Notify(ret.done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	return ret
}

func (a *App) run() {
	json.StartBufferCleaner(30)
	a.router.Route("/api", func(r chi.Router) {
		r.Get("/", handler.New())
	})
	go func() {
		log.Println("Start")
		log.Fatal(http.ListenAndServe(":8282", a.router))

	}()
	<-a.done
	log.Println("Exit")
}
func main() {
	var app = NewApp()
	app.run()
}
