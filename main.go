package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/michaeltintiuc/xpenz-go/pkg/middleware"
	"github.com/michaeltintiuc/xpenz-go/pkg/routes"
)

func main() {
	var wait time.Duration
	var port string
	flag.StringVar(&port, "port", "8080", "the port to listen on - e.g. 8000")
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	r := mux.NewRouter().StrictSlash(true)
	middleware.Init(r)
	routes.Init(r)

	s := &http.Server{
		Handler:      r,
		Addr:         "0.0.0.0:" + port,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}

	log.Printf("Listening on port %s\n", port)

	go func() {
		if err := s.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	// Block until we receive our signal.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	// Doesn't block if no connections, but will otherwise wait until the timeout deadline
	s.Shutdown(ctx)

	log.Println("Shutting down")
	os.Exit(0)
}
