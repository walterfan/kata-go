package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"goroutine-visualizer/internal/monitor"
	"goroutine-visualizer/internal/ws"
)

func main() {
	// 创建监控器
	monitor := monitor.NewGoroutineMonitor()

	// 创建 WebSocket 处理器
	wsHandler := ws.NewWebSocketHandler(monitor)

	// 启动监控器
	go monitor.Start()

	// 启动 WebSocket 广播器
	go wsHandler.StartBroadcaster()

	// 设置路由
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", wsHandler.HandleWebSocket)

	// 静态文件服务器
	fs := http.FileServer(http.Dir("web/"))
	http.Handle("/web/", http.StripPrefix("/web/", fs))

	// 启动模拟任务
	go startSimulatedTasks()

	port := ":8080"
	fmt.Printf("服务器启动在 http://localhost%s\n", port)
	fmt.Printf("WebSocket 连接地址: ws://localhost%s/ws\n", port)

	// 优雅关闭
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		fmt.Println("\n正在关闭服务器...")
		monitor.Stop()
		os.Exit(0)
	}()

	// 启动服务器
	log.Fatal(http.ListenAndServe(port, nil))
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "web/index.html")
}

// 启动一些模拟任务来产生可观察的 goroutine 行为
func startSimulatedTasks() {
	// 任务 1: 周期性任务
	go func() {
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				// 模拟一些计算密集型任务
				go func() {
					start := time.Now()
					sum := 0
					for i := 0; i < 1000000; i++ {
						sum += i
					}
					fmt.Printf("周期性任务完成，耗时: %v\n", time.Since(start))
				}()
			}
		}
	}()

	// 任务 2: 批量处理
	go func() {
		for {
			// 创建一批 goroutine
			for i := 0; i < 5; i++ {
				go func(id int) {
					time.Sleep(time.Duration(100+id*50) * time.Millisecond)
					fmt.Printf("批量任务 %d 完成\n", id)
				}(i)
			}
			time.Sleep(3 * time.Second)
		}
	}()

	// 任务 3: 网络模拟
	go func() {
		for {
			go func() {
				// 模拟网络请求
				time.Sleep(time.Duration(500+time.Now().UnixNano()%1000) * time.Millisecond)
				fmt.Println("网络请求完成")
			}()
			time.Sleep(1 * time.Second)
		}
	}()
}
