package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/storvik/goshrt/http"
	"github.com/storvik/goshrt/postgres"
	"github.com/storvik/goshrt/token"
)

// Serve serves api in graceful manner
func (a *application) Serve() error {
	s := http.NewServer(a.errorLog, a.cfg.Server.Port)

	// Setup server and attach interfaces
	// TODO: Connect to postgresql database
	auth := token.NewAuth(a.cfg.Server.Key)
	db := postgres.NewClient(a.cfg.Database.DB, a.cfg.Database.User, a.cfg.Database.Password, a.cfg.Database.Address)
	err := db.Open()
	if err != nil {
		return err
	}

	s.Auth = auth
	s.ShrtStore = db

	// Server run context
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	// Listen for syscall signals for process to interrupt/quit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig

		// Shutdown signal with grace period of 30 seconds
		shutdownCtx, _ := context.WithTimeout(serverCtx, 30*time.Second)

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()

		// Trigger graceful shutdown
		err := s.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal(err)
		}
		err = s.ShrtStore.Close()
		if err != nil {
			log.Fatal(err)
		}
		serverStopCtx()
	}()

	// Run the server
	err = s.ListenAndServe()
	if err != nil {
		return err
	}

	// Wait for server context to be stopped
	<-serverCtx.Done()

	return nil
}
