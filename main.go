package main

import (
	route "album.com/route"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
	"github.com/rs/cors"
)

func main() {

	//Configuring logger for the app
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	router := route.Config()

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Origin", "X-Requested-With", "Content-Type", "Accept", "Authorization", "z-api-key"},
	})

	server := http.Server{
		Addr:           fmt.Sprintf(":%d", 9000),
		Handler:        c.Handler(router),
		ReadTimeout:    100 * time.Second,
		WriteTimeout:   100 * time.Second,
		MaxHeaderBytes: 1 << 16,
	}

	done := make(chan bool)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		log.Println("Server is shutting down...")

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		//Shutdown server
		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(ctx); err != nil {
			log.Fatalf("Unable to gracefully shutdown the server: %v\n", err)
		}

		//Close channels
		close(quit)
		close(done)
	}()

	log.Printf("Server is listening on: %d", 9000)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("Unable to listen: %v\n", err)
	}

	<-done
	log.Fatal("Server stopped")

}
