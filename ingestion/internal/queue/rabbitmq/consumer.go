package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/porcool/ingestion/internal/config"
)

// IngestionMessage represents the message structure received from RabbitMQ
// Message format: {successfullyIngestedFirestoreDocsID: string}
type IngestionMessage struct {
	SuccessfullyIngestedFirestoreDocsID string `json:"successfullyIngestedFirestoreDocsID"`
}

// MessageHandler is a function type for handling ingestion messages
type MessageHandler func(ctx context.Context, msg IngestionMessage) error

// Consumer represents a RabbitMQ consumer
type Consumer struct {
	conn       *amqp.Connection
	channel    *amqp.Channel
	cfg        config.RabbitMQConfig
	handler    MessageHandler
	stopChan   chan struct{}
	wg         sync.WaitGroup
	running    bool
	mu         sync.Mutex
}

// NewConsumer creates a new RabbitMQ consumer
func NewConsumer(cfg config.RabbitMQConfig, handler MessageHandler) (*Consumer, error) {
	conn, err := amqp.Dial(cfg.URI)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	// Declare the queue (creates it if it doesn't exist)
	_, err = ch.QueueDeclare(
		cfg.QueueName, // name
		true,          // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare queue: %w", err)
	}

	return &Consumer{
		conn:     conn,
		channel:  ch,
		cfg:      cfg,
		handler:  handler,
		stopChan: make(chan struct{}),
	}, nil
}

// Start starts consuming messages from the queue
func (c *Consumer) Start() error {
	c.mu.Lock()
	if c.running {
		c.mu.Unlock()
		return nil
	}
	c.running = true
	c.mu.Unlock()

	// Set prefetch count to 1 to process one message at a time
	err := c.channel.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		return fmt.Errorf("failed to set QoS: %w", err)
	}

	msgs, err := c.channel.Consume(
		c.cfg.QueueName, // queue
		"",              // consumer tag (auto-generated)
		false,           // auto-ack
		false,           // exclusive
		false,           // no-local
		false,           // no-wait
		nil,             // args
	)
	if err != nil {
		return fmt.Errorf("failed to register consumer: %w", err)
	}

	c.wg.Add(1)
	go c.consume(msgs)

	log.Printf("RabbitMQ consumer started, listening on queue: %s", c.cfg.QueueName)
	return nil
}

// consume processes incoming messages
func (c *Consumer) consume(msgs <-chan amqp.Delivery) {
	defer c.wg.Done()

	for {
		select {
		case <-c.stopChan:
			log.Println("Stopping RabbitMQ consumer...")
			return
		case msg, ok := <-msgs:
			if !ok {
				log.Println("RabbitMQ channel closed")
				return
			}

			c.processMessage(msg)
		}
	}
}

// processMessage processes a single message
func (c *Consumer) processMessage(msg amqp.Delivery) {
	log.Printf("Received message: %s", string(msg.Body))

	var ingestionMsg IngestionMessage
	if err := json.Unmarshal(msg.Body, &ingestionMsg); err != nil {
		log.Printf("Error unmarshaling message: %v", err)
		// Reject the message without requeue for invalid JSON
		msg.Nack(false, false)
		return
	}

	if ingestionMsg.SuccessfullyIngestedFirestoreDocsID == "" {
		log.Printf("Invalid message: missing successfullyIngestedFirestoreDocsID")
		msg.Nack(false, false)
		return
	}

	// Process the message using the handler
	ctx := context.Background()
	if err := c.handler(ctx, ingestionMsg); err != nil {
		log.Printf("Error processing message: %v", err)
		// Requeue the message on processing error
		msg.Nack(false, true)
		return
	}

	// Acknowledge the message on success
	if err := msg.Ack(false); err != nil {
		log.Printf("Error acknowledging message: %v", err)
	}

	log.Printf("Successfully processed message for document ID: %s", ingestionMsg.SuccessfullyIngestedFirestoreDocsID)
}

// Stop stops the consumer gracefully
func (c *Consumer) Stop() {
	c.mu.Lock()
	if !c.running {
		c.mu.Unlock()
		return
	}
	c.running = false
	c.mu.Unlock()

	close(c.stopChan)
	c.wg.Wait()
}

// Close closes the RabbitMQ connection
func (c *Consumer) Close() error {
	if c.channel != nil {
		if err := c.channel.Close(); err != nil {
			log.Printf("Error closing channel: %v", err)
		}
	}
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// IsRunning returns whether the consumer is running
func (c *Consumer) IsRunning() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.running
}
