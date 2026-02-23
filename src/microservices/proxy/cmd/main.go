package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"proxy/server"
	"strconv"
)

const(
	MOVIES_API_PREFIX = "/api/movies"
)

func main() {
	compositeHandler, err := provideProxyHandler()
	if err != nil {
		log.Fatalf("failed to create proxyHandler: %v\n", err)
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/api/movies", compositeHandler.Handle)
	log.Printf("Starting server on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}


func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvBool(key string, defaultValue bool) (bool, error) {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue, nil
	}
	res, err := strconv.ParseBool(value)
	if err != nil {
		return false, fmt.Errorf("failed to parse %s env as bool: %v", key, err)
	}

	return res, nil
}

func getEnvInt(key string, defaultValue int) (int, error) {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue, nil
	}
	res, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("failed to parse %s env as int: %v", key, err)
	}

	return res, nil
}

func provideProxyHandler() (*server.CompositeProxyHandler, error) {
	percentage, err :=  getEnvInt("MOVIES_MIGRATION_PERCENT", 0)
	if err != nil {
		return nil, err
	}
	monolithURL := getEnv("MONOLITH_URL", "")
	if monolithURL == "" {
		return nil, fmt.Errorf("uncpecified MONOLITH_URL")
	}
	gradualMigrationEnabled, err := getEnvBool("GRADUAL_MIGRATION", false)
	if err != nil {
		return nil, err
	}
	if !gradualMigrationEnabled {
		return  server.NewCompositeProxyHandler([]server.ProxyHandler {
			server.NewBasicProxyHandler(monolithURL, MOVIES_API_PREFIX),
		}), nil
	}

	moviesURL := getEnv("MOVIES_SERVICE_URL", "")
	if moviesURL == "" {
		return nil, fmt.Errorf("uncpecified MOVIES_SERVICE_URL")
	}

	return server.NewCompositeProxyHandler([]server.ProxyHandler {
		server.NewDualPercentageProxyHandler(
			moviesURL,
			monolithURL,
			MOVIES_API_PREFIX,
			percentage,
		),
	}), nil
}