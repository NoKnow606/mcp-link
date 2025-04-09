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
						Value:   "omnimcp",
						Usage:   "MongoDB database name",
						EnvVars: []string{"MONGODB_DATABASE"},
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

					// 初始化API服务器配置
					return runServer(c.String("host"), c.Int("port"))
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func runServer(host string, port int) error {
	// Create server address
	addr := fmt.Sprintf("%s:%d", host, port)

	baseURL := fmt.Sprintf("http://%s", addr)

	// Configure the SSE server
	ss := utils.NewSSEServer()

	// Get MongoDB client
	mongoClient, err := mongo.GetDefaultClient()
	if err != nil {
		return fmt.Errorf("failed to get MongoDB client: %w", err)
	}

	// Initialize API server config repository (先初始化它，因为SSE服务需要用到它)
	apiServerConfigRepo, err := repositories.NewAPIServerConfigRepository(mongoClient)
	if err != nil {
		return fmt.Errorf("failed to create API server config repository: %w", err)
	}

	// Initialize SSE config repository
	sseConfigRepo, err := repositories.NewSSEConfigRepository(mongoClient)
	if err != nil {
		return fmt.Errorf("failed to create SSE config repository: %w", err)
	}

	// Initialize SSE config service with API server config repository
	sseConfigService := services.NewSSEConfigServiceWithAPIRepo(sseConfigRepo, apiServerConfigRepo)

	// Initialize SSE config controller
	sseConfigController := controllers.NewSSEConfigController(sseConfigService, ss, baseURL)

	// Initialize API server config service
	apiServerConfigService := services.NewAPIServerConfigService(apiServerConfigRepo)

	// Initialize API server config controller
	apiServerConfigController := controllers.NewAPIServerConfigController(apiServerConfigService)

	// Initialize router with both controllers
	apiRouter := router.NewRouter(sseConfigController, apiServerConfigController)

	// Create HTTP server with CORS middleware and router
	mux := http.NewServeMux()
	mux.Handle("/api/v1/", corsMiddleware(apiRouter))
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

	// 设置HTTP服务器配置
	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	// Channel to listen for interrupt signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		fmt.Printf("Starting server on %s with base URL %s\n", addr, baseURL)
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
