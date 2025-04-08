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

	"github.com/anyisalin/mcp-openapi-to-mcp-adapter/controllers"
	"github.com/anyisalin/mcp-openapi-to-mcp-adapter/db/mongo"
	"github.com/anyisalin/mcp-openapi-to-mcp-adapter/repositories"
	"github.com/anyisalin/mcp-openapi-to-mcp-adapter/router"
	"github.com/anyisalin/mcp-openapi-to-mcp-adapter/services"
	"github.com/anyisalin/mcp-openapi-to-mcp-adapter/utils"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "mcp-link",
		Usage: "Convert OpenAPI to MCP compatible endpoints",
		Commands: []*cli.Command{
			{
				Name:  "serve",
				Usage: "Start the MCP Link server",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:    "port",
						Aliases: []string{"p"},
						Value:   8080,
						Usage:   "Port to listen on",
					},
					&cli.StringFlag{
						Name:    "host",
						Aliases: []string{"H"},
						Value:   "localhost",
						Usage:   "Host to listen on",
					},
					&cli.StringFlag{
						Name:    "mongodb-uri",
						Value:   "mongodb://localhost:47017",
						Usage:   "MongoDB connection URI",
						EnvVars: []string{"MONGODB_URI"},
					},
					&cli.StringFlag{
						Name:    "mongodb-database",
						Value:   "ominmcp",
						Usage:   "MongoDB database name",
						EnvVars: []string{"MONGODB_DATABASE"},
					},
					&cli.StringFlag{
						Name:    "base-url",
						Value:   "http://localhost:8080",
						Usage:   "Base URL for the server",
						EnvVars: []string{"BASE_URL"},
					},
				},
				Action: func(c *cli.Context) error {
					// Initialize MongoDB
					mongoConfig := &mongo.Config{
						URI:      c.String("mongodb-uri"),
						Database: c.String("mongodb-database"),
					}

					if err := mongo.InitMongoDBWithConfig(mongoConfig); err != nil {
						return fmt.Errorf("failed to initialize MongoDB: %w", err)
					}

					return runServer(c.String("host"), c.Int("port"), c.String("base-url"))
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func runServer(host string, port int, baseURL string) error {
	// Create server address
	addr := fmt.Sprintf("%s:%d", host, port)

	// Configure the SSE server
	ss := utils.NewSSEServer()

	// Get MongoDB client
	mongoClient, err := mongo.GetDefaultClient()
	if err != nil {
		return fmt.Errorf("failed to get MongoDB client: %w", err)
	}

	// Initialize repository
	sseConfigRepo, err := repositories.NewSSEConfigRepository(mongoClient)
	if err != nil {
		return fmt.Errorf("failed to create SSE config repository: %w", err)
	}

	// Initialize service
	sseConfigService := services.NewSSEConfigService(sseConfigRepo)

	// Initialize controller
	sseConfigController := controllers.NewSSEConfigController(sseConfigService, ss, baseURL)

	// Initialize router
	apiRouter := router.NewRouter(sseConfigController)

	// Create HTTP server with CORS middleware and router
	mux := http.NewServeMux()
	mux.Handle("/api/", corsMiddleware(apiRouter))
	mux.Handle("/sse/config", corsMiddleware(apiRouter))
	mux.Handle("/sse", corsMiddleware(ss))
	mux.Handle("/message", corsMiddleware(ss))
	
	// 添加健康检查端点
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		// 检查 MongoDB 连接状态
		mongoClient, err := mongo.GetDefaultClient()
		if err != nil {
			http.Error(w, "Database connection error", http.StatusServiceUnavailable)
			return
		}
		
		// 尝试 ping MongoDB
		if err := mongoClient.Ping(r.Context()); err != nil {
			http.Error(w, "Database ping failed", http.StatusServiceUnavailable)
			return
		}
		
		// 返回成功状态
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","message":"Service is healthy"}`))
	})

	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	// Channel to listen for interrupt signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		fmt.Printf("Starting server on %s\n", addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting server: %v\n", err)
		}
	}()

	// Wait for interrupt signal
	<-stop

	// Create a deadline for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown the server
	fmt.Println("Shutting down server...")
	if err := ss.Shutdown(ctx); err != nil {
		log.Fatalf("Error shutting down SSE server: %v\n", err)
	}

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Error shutting down HTTP server: %v\n", err)
	}

	fmt.Println("Server gracefully stopped")
	return nil
}

// corsMiddleware adds CORS headers to allow requests from any origin
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Pass the request to the next handler
		next.ServeHTTP(w, r)
	})
}
