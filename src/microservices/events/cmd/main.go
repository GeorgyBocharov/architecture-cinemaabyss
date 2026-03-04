package main

import (
	"context"
	"events/app"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	kafkaAdapter "events/adapters/kafka"
)


type ConsumerInfo struct {
	Name     string
	Consumer *kafkaAdapter.Consumer
}

func main() {
	config := loadConfig()
	
	container := initContainer(config)
	defer container.Producer.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := setupSignalHandler()

	wg := startServices(ctx, container)

	<-sigChan
	log.Println("Received shutdown signal. Initiating graceful shutdown...")

	gracefulShutdown(ctx, cancel, container, wg)
}

func loadConfig() *app.Config {
	kafkaServers := os.Getenv("KAFKA_BROKERS")
	
	return &app.Config{
		PaymentsTopic: os.Getenv("PAYMENT_TOPIC"),
		UsersTopic:    os.Getenv("USER_TOPIC"),
		MoviesTopic:   os.Getenv("MOVIE_TOPIC"),
		Port: getPort(),
		ConsumerPollTimeout: 100, 
		
		PaymentsConsumerConfig: map[string]interface{}{
			"bootstrap.servers": kafkaServers,
			"group.id":          "payments-consumer",
			"enable.auto.commit": true,
			"auto.offset.reset": "earliest",
		},
		UsersConsumerConfig: map[string]interface{}{
			"bootstrap.servers": kafkaServers,
			"group.id":          "users-consumer",
			"enable.auto.commit": true,
			"auto.offset.reset": "earliest",
		},
		MoviesConsumerConfig: map[string]interface{}{
			"bootstrap.servers": kafkaServers,
			"group.id":          "movies-consumer",
			"enable.auto.commit": true,
			"auto.offset.reset": "earliest",
		},
		ProducerConfig: map[string]interface{}{
			"bootstrap.servers":  kafkaServers,
			"message.timeout.ms": 10000,
			"acks":               "all",
		},
	}
}

func initContainer(config *app.Config) *app.Container {
	container := &app.Container{
		Config: config,
	}

	if err := container.Init(); err != nil {
		log.Fatalf("failed to init container: %v", err)
	}

	return container
}

func setupSignalHandler() <-chan os.Signal {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	return sigChan
}

func startServices(ctx context.Context, container *app.Container) *sync.WaitGroup {
	var wg sync.WaitGroup

	startConsumers(ctx, &wg, container)

	startHTTPServer(&wg, container)

	log.Println("All services started. Press Ctrl+C to stop...")
	return &wg
}

func startConsumers(ctx context.Context, wg *sync.WaitGroup, container *app.Container) {
	consumers := []ConsumerInfo{
		{"Payments", container.PaymentsConsumer},
		{"Movies", container.MoviesConsumer},
		{"Users", container.UsersConsumer},
	}

	for _, ci := range consumers {
		if ci.Consumer == nil {
			log.Printf("Warning: %s consumer is nil, skipping", ci.Name)
			continue
		}

		wg.Add(1)
		go runConsumer(ctx, wg, ci.Name, ci.Consumer)
	}
}

func runConsumer(ctx context.Context, wg *sync.WaitGroup, name string, consumer *kafkaAdapter.Consumer) {
	defer wg.Done()
	
	log.Printf("Starting %s consumer...", name)
	
	if err := consumer.Consume(ctx); err != nil {
		log.Printf("%s consumer stopped with error: %v", name, err)
	} else {
		log.Printf("%s consumer stopped gracefully", name)
	}
}

func startHTTPServer(wg *sync.WaitGroup, container *app.Container) {
	port := getPort()
	
	http.HandleFunc("/api/events/user", container.UsersHandler.Handle)
	http.HandleFunc("/api/events/movie", container.MoviesHandler.Handle)
	http.HandleFunc("/api/events/payment", container.PaymentsHandler.Handle)
	http.HandleFunc("/api/events/health", container.HealthcheckHandler.Handle)

	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Printf("Starting HTTP server on port %s", port)
		
		if err := container.HTTPServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("HTTP server error: %v", err)
		}
	}()

}

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return port
}

func gracefulShutdown(ctx context.Context, cancel context.CancelFunc, container *app.Container, wg *sync.WaitGroup) {
	shutdownHTTPServer(container)

	cancel()

	waitForShutdown(wg)
}

func shutdownHTTPServer(container *app.Container) {
	if container.HTTPServer == nil {
		return
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	log.Println("Shutting down HTTP server...")
	if err := container.HTTPServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("HTTP server shutdown error: %v", err)
	}
}

func waitForShutdown(wg *sync.WaitGroup) {
	done := make(chan struct{})
	
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		log.Println("All services closed successfully")
	case <-time.After(10 * time.Second):
		log.Println("Timeout waiting for services to close")
	}
}