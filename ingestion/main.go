package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/porcool/ingestion/internal/config"
	"github.com/porcool/ingestion/internal/database/mariadb"
	"github.com/porcool/ingestion/internal/database/mongodb"
	"github.com/porcool/ingestion/internal/ingestion"
	"github.com/porcool/ingestion/internal/logging"
	"github.com/porcool/ingestion/internal/queue/rabbitmq"
)

const serviceName = "porcool-ingestion-non-relational-database-to-relational-database"

func main() {
	log.Println("Starting PorCool Ingestion Service (porcool-ingestion-nosql-to-sql-database)...")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger with OpenSearch support
	logger := logging.NewLogger(cfg.OpenSearch, serviceName)
	if err := logger.Initialize(context.Background()); err != nil {
		log.Printf("Warning: Logger initialization issue: %v", err)
	}
	defer logger.Close()

	// Set as default logger and redirect standard log output
	logging.SetDefaultLogger(logger)
	logging.RedirectStdLog(logger)

	// Initialize MariaDB connection
	mariaDB, err := mariadb.NewConnection(cfg.MariaDB)
	if err != nil {
		log.Fatalf("Failed to connect to MariaDB: %v", err)
	}
	defer mariaDB.Close()

	log.Println("Connected to MariaDB successfully")

	// Run migrations
	if err := mariaDB.RunMigrations(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	log.Println("Database migrations completed successfully")

	// Seed domains
	if err := mariaDB.SeedDomains(); err != nil {
		log.Fatalf("Failed to seed domains: %v", err)
	}
	log.Println("Domain seeding completed successfully")

	// Initialize MongoDB connection
	mongoDB, err := mongodb.NewConnection(cfg.MongoDB)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer mongoDB.Close()

	log.Println("Connected to MongoDB successfully")

	// Initialize ingestion service
	svc := ingestion.NewService(mariaDB, mongoDB, cfg)

	// Create message handler that delegates to ingestion service
	messageHandler := func(ctx context.Context, msg rabbitmq.IngestionMessage) error {
		return svc.ProcessIngestionMessage(ctx, msg.SuccessfullyIngestedFirestoreDocsID)
	}

	// Initialize RabbitMQ consumer
	consumer, err := rabbitmq.NewConsumer(cfg.RabbitMQ, messageHandler)
	if err != nil {
		log.Fatalf("Failed to create RabbitMQ consumer: %v", err)
	}
	defer consumer.Close()

	log.Println("Connected to RabbitMQ successfully")

	// Start RabbitMQ consumer
	if err := consumer.Start(); err != nil {
		log.Fatalf("Failed to start RabbitMQ consumer: %v", err)
	}

	log.Printf("Ingestion service started, listening on queue: %s", cfg.RabbitMQ.QueueName)

	// Wait for shutdown signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down ingestion service...")
	consumer.Stop()
	log.Println("Ingestion service stopped")
}
