package ws

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"goroutine-visualizer/internal/monitor"

	"github.com/gorilla/websocket"
)

// WebSocketHandler 处理 WebSocket 连接
type WebSocketHandler struct {
	upgrader  websocket.Upgrader
	clients   map[*websocket.Conn]bool
	clientsMu sync.RWMutex
	monitor   *monitor.GoroutineMonitor
	broadcast chan monitor.SystemInfo
	stopChan  chan struct{}
}

// NewWebSocketHandler 创建新的 WebSocket 处理器
func NewWebSocketHandler(mon *monitor.GoroutineMonitor) *WebSocketHandler {
	return &WebSocketHandler{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // 允许所有来源，生产环境中应该更严格
			},
		},
		clients:   make(map[*websocket.Conn]bool),
		monitor:   mon,
		broadcast: make(chan monitor.SystemInfo, 100),
		stopChan:  make(chan struct{}),
	}
}

// HandleWebSocket 处理 WebSocket 连接
func (h *WebSocketHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket 升级失败: %v", err)
		return
	}

	log.Printf("新的 WebSocket 连接建立")

	h.clientsMu.Lock()
	h.clients[conn] = true
	h.clientsMu.Unlock()

	// 处理连接
	go h.handleConnection(conn)
}

// handleConnection 处理单个 WebSocket 连接
func (h *WebSocketHandler) handleConnection(conn *websocket.Conn) {
	defer func() {
		h.clientsMu.Lock()
		delete(h.clients, conn)
		h.clientsMu.Unlock()
		conn.Close()
		log.Printf("WebSocket 连接关闭")
	}()

	// 设置读取截止时间
	conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	// 启动 ping 协程
	go h.pingHandler(conn)

	// 读取消息循环
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket 读取错误: %v", err)
			}
			break
		}
	}
}

// pingHandler 定期发送 ping 消息
func (h *WebSocketHandler) pingHandler(conn *websocket.Conn) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := conn.WriteControl(websocket.PingMessage, nil, time.Now().Add(10*time.Second)); err != nil {
				return
			}
		case <-h.stopChan:
			return
		}
	}
}

// StartBroadcaster 启动广播器
func (h *WebSocketHandler) StartBroadcaster() {
	// 订阅监控器的数据
	dataChan := h.monitor.Subscribe()

	go func() {
		for {
			select {
			case data := <-dataChan:
				h.broadcast <- data
			case <-h.stopChan:
				return
			}
		}
	}()

	// 广播循环
	go func() {
		for {
			select {
			case data := <-h.broadcast:
				h.broadcastToClients(data)
			case <-h.stopChan:
				return
			}
		}
	}()
}

// broadcastToClients 向所有客户端广播数据
func (h *WebSocketHandler) broadcastToClients(data monitor.SystemInfo) {
	h.clientsMu.RLock()
	defer h.clientsMu.RUnlock()

	// 序列化数据
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("JSON 序列化失败: %v", err)
		return
	}

	// 向所有客户端发送数据
	for client := range h.clients {
		select {
		case <-h.stopChan:
			return
		default:
			if err := client.WriteMessage(websocket.TextMessage, jsonData); err != nil {
				log.Printf("WebSocket 写入失败: %v", err)
				// 移除失败的客户端
				go func(c *websocket.Conn) {
					h.clientsMu.Lock()
					delete(h.clients, c)
					h.clientsMu.Unlock()
					c.Close()
				}(client)
			}
		}
	}
}

// Stop 停止 WebSocket 处理器
func (h *WebSocketHandler) Stop() {
	close(h.stopChan)

	h.clientsMu.Lock()
	defer h.clientsMu.Unlock()

	for client := range h.clients {
		client.Close()
	}
}
