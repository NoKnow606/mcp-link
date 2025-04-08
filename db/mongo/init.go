package mongo

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var initOnce sync.Once
var shutdownHookRegistered bool
var mu sync.Mutex

// InitMongoDB 初始化MongoDB连接并注册优雅关闭钩子
func InitMongoDB() error {
	mu.Lock()
	defer mu.Unlock()

	var initErr error
	initOnce.Do(func() {
		// 从环境变量加载配置
		config := LoadConfigFromEnv()

		// 设置连接超时上下文
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// 初始化默认客户端
		if err := InitDefaultClient(ctx, config); err != nil {
			initErr = fmt.Errorf("failed to initialize MongoDB client: %w", err)
			return
		}

		// 测试连接
		client, err := GetDefaultClient()
		if err != nil {
			initErr = fmt.Errorf("failed to get MongoDB client: %w", err)
			return
		}

		if err := client.Ping(ctx); err != nil {
			initErr = fmt.Errorf("failed to ping MongoDB: %w", err)
			return
		}

		log.Println("Successfully connected to MongoDB at", config.URI)

		// 注册关闭钩子
		if !shutdownHookRegistered {
			registerShutdownHook()
			shutdownHookRegistered = true
		}
	})

	return initErr
}

// InitMongoDBWithConfig 使用指定配置初始化MongoDB连接
func InitMongoDBWithConfig(config *Config) error {
	mu.Lock()
	defer mu.Unlock()

	// 设置连接超时上下文
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 初始化默认客户端
	if err := InitDefaultClient(ctx, config); err != nil {
		return fmt.Errorf("failed to initialize MongoDB client with custom config: %w", err)
	}

	// 测试连接
	client, err := GetDefaultClient()
	if err != nil {
		return fmt.Errorf("failed to get MongoDB client: %w", err)
	}

	if err := client.Ping(ctx); err != nil {
		return fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	log.Println("Successfully connected to MongoDB at", config.URI)

	// 注册关闭钩子
	if !shutdownHookRegistered {
		registerShutdownHook()
		shutdownHookRegistered = true
	}

	return nil
}

// registerShutdownHook 注册程序关闭时清理MongoDB连接的钩子
func registerShutdownHook() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		<-c
		log.Println("Received shutdown signal, closing MongoDB connections...")
		
		// 设置关闭超时上下文
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		
		// 关闭MongoDB连接
		if err := CloseDefaultClient(ctx); err != nil {
			log.Printf("Error closing MongoDB connections: %v\n", err)
		} else {
			log.Println("MongoDB connections closed successfully")
		}
	}()
} 