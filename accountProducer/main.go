package main

import (
	"accountProducer/configurations"
	"accountProducer/handlers"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	"github.com/joho/godotenv"
)

// producer work is to produce data and publish it to the kafka topic
// our gateawy can act as a producer

func main() {

	fmt.Println("Starting to develop banking application")
	// logging app file
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file:", err)
		// Optionally continue with default values or exit
	}
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Some error occured in creating or opeing a log file", err)
		os.Exit(1)
	}

	// logger
	loggs := hclog.New(&hclog.LoggerOptions{
		Name:       "Banking App",
		Output:     logFile,
		Level:      hclog.Debug,
		JSONFormat: false,
	})

	// get the URI
	uri, err := configurations.NewAppConfig()
	if err != nil {
		loggs.Error("Not able to create Retrieve App Configurations")
		os.Exit(1)
	}

	// handlers
	handler := handlers.NewUserHandler()

	// router
	router := mux.NewRouter()

	// register handlers
	handler.RegisterRoutes(router)

	opts := hclog.StandardLoggerOptions{
		InferLevels: true,
	}

	//cerate a new server
	httpServer := http.Server{
		Addr:         *uri.GetAppURI(),
		Handler:      router,
		ErrorLog:     loggs.StandardLogger(&opts),
		ReadTimeout:  5 * time.Second,   // max time to read request from the client
		WriteTimeout: 10 * time.Second,  // max time to write response to the client
		IdleTimeout:  120 * time.Second, // max time for connections using TCP Keep-Alive
	}

	go func() {
		loggs.Info("starting server on ", "Port", *uri.GetAppURI())
		fmt.Println(*uri.GetAppURI())
		err := httpServer.ListenAndServe()
		if err != nil {
			loggs.Info("Error starting server", err)
			os.Exit(1)
		}
	}()

	//gracefully shutting down the server
	signalChannel := make(chan os.Signal)
	signal.Notify(signalChannel, os.Interrupt)
	signal.Notify(signalChannel, os.Kill)

	//blocking statement
	waitingForChanel := <-signalChannel

	loggs.Info("System Interruptions Received", waitingForChanel)

	//gracefully shutting down the server

	newctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	httpServer.Shutdown(newctx)
}
