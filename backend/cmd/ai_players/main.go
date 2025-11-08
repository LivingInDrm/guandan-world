package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"guandan-world/backend/test"
)

func main() {
	// 解析命令行参数
	serverURL := flag.String("server", "localhost:8080", "Backend server URL (host:port)")
	roomID := flag.String("room-id", "", "Room ID to join (required)")
	verbose := flag.Bool("verbose", false, "Enable verbose logging")
	numPlayers := flag.Int("num-players", 3, "Number of AI players to create (default: 3)")
	usernamePrefix := flag.String("username-prefix", "ai_player", "Username prefix for AI players")
	password := flag.String("password", "ai123456", "Password for AI players")

	flag.Parse()

	// 验证必需参数
	if *roomID == "" {
		log.Fatal("Error: -room-id is required")
	}

	// 验证玩家数量
	if *numPlayers < 1 || *numPlayers > 3 {
		log.Fatal("Error: -num-players must be between 1 and 3")
	}

	log.Printf("=== AI Players Test Client ===")
	log.Printf("Server: %s", *serverURL)
	log.Printf("Room ID: %s", *roomID)
	log.Printf("Number of players: %d", *numPlayers)
	log.Printf("Verbose: %v", *verbose)
	log.Printf("==============================")
	log.Println()

	// 设置信号处理
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// 创建AI玩家客户端
	clients := make([]*test.AIPlayerClient, *numPlayers)
	for i := 0; i < *numPlayers; i++ {
		username := fmt.Sprintf("%s_%d", *usernamePrefix, i+1)
		clients[i] = test.NewAIPlayerClient(*serverURL, *roomID, username, *password, *verbose)
	}

	// 使用WaitGroup等待所有客户端启动
	var wg sync.WaitGroup
	errChan := make(chan error, *numPlayers)

	// 启动所有AI客户端
	for i, client := range clients {
		wg.Add(1)
		go func(idx int, c *test.AIPlayerClient) {
			defer wg.Done()

			// 添加短暂延迟，避免同时注册导致竞争
			time.Sleep(time.Duration(idx) * 200 * time.Millisecond)

			if err := c.Start(); err != nil {
				errChan <- fmt.Errorf("client %d failed to start: %w", idx+1, err)
				return
			}

			log.Printf("[Client %d] Started successfully", idx+1)
		}(i, client)
	}

	// 等待所有客户端启动完成
	go func() {
		wg.Wait()
		close(errChan)
	}()

	// 检查启动错误
	startErrors := make([]error, 0)
	for err := range errChan {
		startErrors = append(startErrors, err)
		log.Printf("Error: %v", err)
	}

	if len(startErrors) > 0 {
		log.Printf("Failed to start %d client(s), shutting down...", len(startErrors))
		for _, client := range clients {
			client.Stop()
		}
		os.Exit(1)
	}

	log.Println()
	log.Printf("All %d AI players connected and ready!", *numPlayers)
	log.Println("Waiting for game to start and complete...")
	log.Println("Press Ctrl+C to exit")
	log.Println()

	// 等待信号或所有客户端完成
	doneChan := make(chan struct{})
	go func() {
		for _, client := range clients {
			client.Wait()
		}
		close(doneChan)
	}()

	select {
	case <-sigChan:
		log.Println()
		log.Println("Received interrupt signal, shutting down...")
		for i, client := range clients {
			log.Printf("Stopping client %d...", i+1)
			client.Stop()
		}
		log.Println("All clients stopped")

	case <-doneChan:
		log.Println()
		log.Println("All games completed!")
	}

	log.Println("AI Players Test Client exited")
}

