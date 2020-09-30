package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"git.sr.ht/~tristan957/noteriety/noteriety-ls/services"
	"github.com/gorilla/mux"
	"github.com/gorilla/rpc/v2"
	"github.com/gorilla/rpc/v2/json2"
)

func main() {
	r := mux.NewRouter()
	s := rpc.NewServer()
	r.Handle("/jsonrpc", s)
	s.RegisterCodec(json2.NewCodec(), "application/json")
	s.RegisterService(new(services.HelloService), "HelloService")
	srv := http.Server{
		Addr:         "0.0.0.0:10260",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	<-c

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	srv.Shutdown(ctx)
}
