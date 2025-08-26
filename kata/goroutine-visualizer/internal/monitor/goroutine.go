package monitor

import (
	"runtime"
	"sync"
	"time"
)

// GoroutineInfo 表示单个 goroutine 的信息
type GoroutineInfo struct {
	ID        int    `json:"id"`
	State     string `json:"state"`
	Function  string `json:"function"`
	File      string `json:"file"`
	Line      int    `json:"line"`
	Duration  int64  `json:"duration"`
	CreatedAt int64  `json:"created_at"`
}

// SystemInfo 表示系统信息
type SystemInfo struct {
	NumGoroutine int              `json:"num_goroutine"`
	NumCPU       int              `json:"num_cpu"`
	GOMAXPROCS   int              `json:"gomaxprocs"`
	MemStats     runtime.MemStats `json:"mem_stats"`

	Timestamp  int64           `json:"timestamp"`
	Goroutines []GoroutineInfo `json:"goroutines"`
}

// GoroutineMonitor 监控 goroutine 的状态
type GoroutineMonitor struct {
	mu          sync.RWMutex
	subscribers []chan SystemInfo
	ticker      *time.Ticker
	stopChan    chan struct{}
	running     bool
}

// NewGoroutineMonitor 创建新的监控器
func NewGoroutineMonitor() *GoroutineMonitor {
	return &GoroutineMonitor{
		subscribers: make([]chan SystemInfo, 0),
		stopChan:    make(chan struct{}),
	}
}

// Start 开始监控
func (m *GoroutineMonitor) Start() {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.running {
		return
	}

	m.running = true
	m.ticker = time.NewTicker(1 * time.Second) // 每秒收集一次数据

	go m.collectLoop()
}

// Stop 停止监控
func (m *GoroutineMonitor) Stop() {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.running {
		return
	}

	m.running = false
	close(m.stopChan)
	m.ticker.Stop()
}

// Subscribe 订阅系统信息更新
func (m *GoroutineMonitor) Subscribe() <-chan SystemInfo {
	m.mu.Lock()
	defer m.mu.Unlock()

	ch := make(chan SystemInfo, 10)
	m.subscribers = append(m.subscribers, ch)
	return ch
}

// collectLoop 收集数据的主循环
func (m *GoroutineMonitor) collectLoop() {
	for {
		select {
		case <-m.ticker.C:
			info := m.collectSystemInfo()
			m.broadcast(info)
		case <-m.stopChan:
			return
		}
	}
}

// collectSystemInfo 收集系统信息
func (m *GoroutineMonitor) collectSystemInfo() SystemInfo {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	// 获取 goroutine 堆栈信息
	buf := make([]byte, 1024*1024*2) // 2MB buffer
	stackSize := runtime.Stack(buf, true)
	goroutines := parseGoroutineStack(string(buf[:stackSize]))

	return SystemInfo{
		NumGoroutine: runtime.NumGoroutine(),
		NumCPU:       runtime.NumCPU(),
		GOMAXPROCS:   runtime.GOMAXPROCS(0),
		MemStats:     memStats,
		Timestamp:    time.Now().UnixNano() / int64(time.Millisecond),
		Goroutines:   goroutines,
	}
}

// parseGoroutineStack 解析 goroutine 堆栈信息
func parseGoroutineStack(stack string) []GoroutineInfo {
	// 这里是一个简化的解析器
	// 在实际应用中，您可能需要更复杂的解析逻辑
	goroutines := make([]GoroutineInfo, 0)

	// 计算当前时间戳
	now := time.Now().UnixNano() / int64(time.Millisecond)

	// 为演示目的，我们创建一些模拟的 goroutine 信息
	numGoroutines := runtime.NumGoroutine()
	for i := 0; i < numGoroutines && i < 50; i++ { // 限制显示数量
		goroutines = append(goroutines, GoroutineInfo{
			ID:        i + 1,
			State:     getRandomState(),
			Function:  "example.function",
			File:      "main.go",
			Line:      42 + i,
			Duration:  int64(i * 100),
			CreatedAt: now - int64(i*1000),
		})
	}

	return goroutines
}

// getRandomState 获取随机状态（用于演示）
func getRandomState() string {
	states := []string{"running", "runnable", "waiting", "blocked", "dead"}
	return states[time.Now().UnixNano()%int64(len(states))]
}

// broadcast 广播信息给所有订阅者
func (m *GoroutineMonitor) broadcast(info SystemInfo) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, ch := range m.subscribers {
		select {
		case ch <- info:
		default:
			// 如果 channel 满了，跳过这个订阅者
		}
	}
}
