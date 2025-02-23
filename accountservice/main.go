package main

import (
	"accountservice/configurations"
	"accountservice/database"
	"accountservice/handlers"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
)

func main() {

	fmt.Println("Starting to develop banking application")
	// logging app file
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

	// Connect to PostgreSQL
	dsn := "postgres://postgres:abcd@localhost:5432/accounts?sslmode=disable"
	db := database.NewPostgresPoolDB(dsn, 10, 2)
	ctx := context.Background()

	if err := db.Connect(ctx); err != nil {
		fmt.Printf("Connection failed: %v\n", err)
		return
	}
	defer db.Close(ctx)

	// handlers
	handler := handlers.NewUserHandler(db)

	// router
	router := mux.NewRouter()

	// register handlers
	handler.RegisterRoutes(router)

	opts := hclog.StandardLoggerOptions{
		InferLevels: true,
	}

	// //cerate a new server
	httpServer := http.Server{
		Addr:         *uri.GetAppURI(),
		Handler:      router,
		ErrorLog:     loggs.StandardLogger(&opts),
		ReadTimeout:  5 * time.Second,   // max time to read request from the client
		WriteTimeout: 10 * time.Second,  // max time to write response to the client
		IdleTimeout:  120 * time.Second, // max time for connections using TCP Keep-Alive
	}

	go func() {
		loggs.Info("starting server on ", "Port", uri.GetAppURI())
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
