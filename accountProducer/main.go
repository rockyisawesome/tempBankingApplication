package main

import (
	"accountProducer/configurations" // Importing configurations for app and MongoDB settings
	"accountProducer/database"       // Importing database package for MongoDB operations
	"accountProducer/handlers"       // Importing handlers for HTTP request handling
	"context"                        // Importing context for request-scoped operations and timeouts
	"fmt"                            // Importing fmt for formatted output
	"net/http"                       // Importing http for running the HTTP server
	"os"                             // Importing os for file operations and system signals
	"os/signal"                      // Importing signal for handling OS interrupts
	"time"                           // Importing time for timeout and duration settings

	_ "accountProducer/docs" // Importing docs package (Swagger) as a side effect for documentation

	"github.com/gorilla/mux"                        // Importing mux for HTTP routing
	"github.com/hashicorp/go-hclog"                 // Importing hclog for structured logging
	"github.com/joho/godotenv"                      // Importing godotenv for loading environment variables
	httpSwagger "github.com/swaggo/http-swagger/v2" // Importing http-swagger for Swagger UI
)

// @title Swagger Example API
// @version 1.0
// @description This is a sample server Petstore server.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host petstore.swagger.io
// @BasePath /v2

// main is the entry point of the banking application.
// It sets up logging, configurations, database connections, HTTP handlers, and the server,
// then runs the server with graceful shutdown support.
func main() {
	// Log the application startup
	fmt.Println("Starting to develop banking application")

	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file:", err) // Log error but continue with defaults
	}

	// Open or create a log file for application logs
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Some error occurred in creating or opening a log file", err)
		os.Exit(1) // Exit if log file cannot be created/opened
	}

	// Initialize a logger with debug level and file output
	loggs := hclog.New(&hclog.LoggerOptions{
		Name:       "Banking App", // Logger name for identification
		Output:     logFile,       // Direct logs to the file
		Level:      hclog.Debug,   // Set log level to debug for detailed output
		JSONFormat: false,         // Use plain text format instead of JSON
	})

	// Retrieve application configuration (e.g., server URI)
	uri, err := configurations.NewAppConfig()
	if err != nil {
		loggs.Error("Not able to create Retrieve App Configurations")
		os.Exit(1) // Exit if app config cannot be loaded
	}

	// Retrieve MongoDB configuration from environment variables
	mongodbconfig, err := configurations.NewMongoDbConfig()
	if err != nil {
		loggs.Error("Not able to create Retrieve App Configurations")
		os.Exit(1) // Exit if MongoDB config cannot be loaded
	}

	// Initialize MongoDB instance with config and logger
	mongodb := database.NewMongoDB(mongodbconfig, &loggs)
	ctx := context.Background() // Create a background context for database operations

	// Connect to MongoDB
	if err := mongodb.Connect(ctx); err != nil {
		fmt.Printf("Connection failed: %v\n", err)
		return // Exit if database connection fails
	}
	defer mongodb.Disconnect(ctx) // Ensure MongoDB disconnects when main exits

	// Create a new handler instance with MongoDB and logger
	handler := handlers.NewUserHandler(mongodb, &loggs)

	// Initialize the HTTP router
	router := mux.NewRouter()

	// Set up Swagger UI endpoint for API documentation
	router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"), // Point to the Swagger JSON file
	))

	// Register handler routes with the router
	handler.RegisterRoutes(router)

	// Configure standard logger options for HTTP server error logging
	opts := hclog.StandardLoggerOptions{
		InferLevels: true, // Automatically infer log levels from messages
	}

	// Create and configure the HTTP server
	httpServer := http.Server{
		Addr:         *uri.GetAppURI(),            // Set server address from config
		Handler:      router,                      // Use the mux router for handling requests
		ErrorLog:     loggs.StandardLogger(&opts), // Direct server errors to the logger
		ReadTimeout:  5 * time.Second,             // Max time to read client request
		WriteTimeout: 10 * time.Second,            // Max time to write response to client
		IdleTimeout:  120 * time.Second,           // Max time for idle TCP connections
	}

	// Start the server in a goroutine to allow concurrent signal handling
	go func() {
		loggs.Info("starting server on ", "Port", *uri.GetAppURI())
		fmt.Println(*uri.GetAppURI()) // Log the server address
		err := httpServer.ListenAndServe()
		if err != nil {
			loggs.Info("Error starting server", err)
			os.Exit(1) // Exit if server fails to start
		}
	}()

	// Set up signal handling for graceful shutdown
	signalChannel := make(chan os.Signal)
	signal.Notify(signalChannel, os.Interrupt) // Capture Ctrl+C
	signal.Notify(signalChannel, os.Kill)      // Capture kill signals

	// Block until a signal is received
	waitingForChanel := <-signalChannel
	loggs.Info("System Interruptions Received", waitingForChanel)

	// Gracefully shut down the server with a 30-second timeout
	newctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	httpServer.Shutdown(newctx) // Shutdown the server cleanly
}
