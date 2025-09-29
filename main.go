package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	//http.HandleFunc("/", homePage)
	//http.HandleFunc("/ws", handleConnections)

	//go handleMessages()
	manager := NewManager()
	http.HandleFunc("/ws", manager.wsHandler)
	fmt.Println("Websocket Server started on :8080")
	s := &http.Server{
		Addr: ":8080",
		//Handler:        rou,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	//log.Fatal(s.ListenAndServe())

	go func() {
		log.Println("Starting server on :8080")
		err := s.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("%v error has occured", err)
		}

	}()
	//This is a channel that is on standby to receice an event stream such as a ctrl+c
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	// This is to handle graceful shutown. It accounts for both a failure and an Intentional Shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		log.Fatalf("Shutdown failed: %v", err)
	}

	log.Println("Server exited cleanly")
}
