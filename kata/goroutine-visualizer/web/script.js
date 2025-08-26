// 全局变量
let ws;
let isConnected = false;
let goroutineChart;
let memoryChart;
let stateChart;
let goroutineData = [];
let memoryData = [];
let maxDataPoints = 50;

// DOM 元素
const connectionStatus = document.getElementById('connection-status');
const connectionText = document.getElementById('connection-text');
const goroutineCount = document.getElementById('goroutine-count');
const cpuCount = document.getElementById('cpu-count');
const gomaxprocs = document.getElementById('gomaxprocs');
const memoryUsage = document.getElementById('memory-usage');
const goroutineList = document.getElementById('goroutine-list');
const autoScroll = document.getElementById('auto-scroll');
const maxDisplay = document.getElementById('max-display');

// 初始化
document.addEventListener('DOMContentLoaded', function() {
    initializeCharts();
    connectWebSocket();
    setupEventListeners();
});

// 连接 WebSocket
function connectWebSocket() {
    const wsProtocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const wsUrl = `${wsProtocol}//${window.location.host}/ws`;
    
    console.log('正在连接到:', wsUrl);
    ws = new WebSocket(wsUrl);
    
    ws.onopen = function() {
        console.log('WebSocket 连接已建立');
        isConnected = true;
        updateConnectionStatus(true);
    };
    
    ws.onmessage = function(event) {
        try {
            const data = JSON.parse(event.data);
            handleSystemData(data);
        } catch (error) {
            console.error('解析数据失败:', error);
        }
    };
    
    ws.onclose = function() {
        console.log('WebSocket 连接已关闭');
        isConnected = false;
        updateConnectionStatus(false);
        
        // 尝试重连
        setTimeout(connectWebSocket, 3000);
    };
    
    ws.onerror = function(error) {
        console.error('WebSocket 错误:', error);
        isConnected = false;
        updateConnectionStatus(false);
    };
}

// 更新连接状态
function updateConnectionStatus(connected) {
    if (connected) {
        connectionStatus.className = 'status-dot connected';
        connectionText.textContent = '已连接';
    } else {
        connectionStatus.className = 'status-dot disconnected';
        connectionText.textContent = '未连接';
    }
}

// 处理系统数据
function handleSystemData(data) {
    updateOverviewMetrics(data);
    updateCharts(data);
    updateGoroutineList(data.goroutines || []);
}

// 更新概览指标
function updateOverviewMetrics(data) {
    goroutineCount.textContent = data.num_goroutine || 0;
    cpuCount.textContent = data.num_cpu || 0;
    gomaxprocs.textContent = data.gomaxprocs || 0;
    
    // 格式化内存使用
    if (data.mem_stats) {
        const memMB = Math.round(data.mem_stats.Alloc / 1024 / 1024);
        memoryUsage.textContent = `${memMB} MB`;
    }
}

// 更新图表
function updateCharts(data) {
    const now = new Date(data.timestamp);
    
    // 更新 goroutine 数量图表
    goroutineData.push({
        x: now,
        y: data.num_goroutine || 0
    });
    
    // 更新内存使用图表
    if (data.mem_stats) {
        memoryData.push({
            x: now,
            y: Math.round(data.mem_stats.Alloc / 1024 / 1024)
        });
    }
    
    // 限制数据点数量
    if (goroutineData.length > maxDataPoints) {
        goroutineData.shift();
    }
    if (memoryData.length > maxDataPoints) {
        memoryData.shift();
    }
    
    // 更新图表
    goroutineChart.update('none');
    memoryChart.update('none');
    
    // 更新状态分布图表
    updateStateChart(data.goroutines || []);
}

// 更新状态分布图表
function updateStateChart(goroutines) {
    const states = {};
    goroutines.forEach(g => {
        states[g.state] = (states[g.state] || 0) + 1;
    });
    
    stateChart.data.labels = Object.keys(states);
    stateChart.data.datasets[0].data = Object.values(states);
    stateChart.update('none');
}

// 更新 goroutine 列表
function updateGoroutineList(goroutines) {
    const maxCount = parseInt(maxDisplay.value);
    const displayGoroutines = goroutines.slice(0, maxCount);
    
    let html = `
        <div class="goroutine-header">
            <div>ID</div>
            <div>状态</div>
            <div>函数</div>
            <div>行号</div>
            <div>持续时间</div>
        </div>
    `;
    
    displayGoroutines.forEach(g => {
        const duration = formatDuration(g.duration);
        html += `
            <div class="goroutine-item">
                <div class="goroutine-id">#${g.id}</div>
                <div class="goroutine-state state-${g.state}">${g.state}</div>
                <div class="goroutine-function">${g.function}</div>
                <div>${g.line}</div>
                <div class="goroutine-duration">${duration}</div>
            </div>
        `;
    });
    
    goroutineList.innerHTML = html;
    
    // 自动滚动
    if (autoScroll.checked) {
        goroutineList.scrollTop = goroutineList.scrollHeight;
    }
}

// 格式化持续时间
function formatDuration(ms) {
    if (ms < 1000) {
        return `${ms}ms`;
    } else if (ms < 60000) {
        return `${Math.round(ms / 1000)}s`;
    } else {
        return `${Math.round(ms / 60000)}m`;
    }
}

// 初始化图表
function initializeCharts() {
    // Goroutine 数量趋势图
    const goroutineCtx = document.getElementById('goroutine-chart').getContext('2d');
    goroutineChart = new Chart(goroutineCtx, {
        type: 'line',
        data: {
            datasets: [{
                label: 'Goroutine 数量',
                data: goroutineData,
                borderColor: '#3498db',
                backgroundColor: 'rgba(52, 152, 219, 0.1)',
                fill: true,
                tension: 0.4
            }]
        },
        options: {
            responsive: true,
            maintainAspectRatio: false,
            scales: {
                x: {
                    type: 'time',
                    time: {
                        unit: 'second',
                        displayFormats: {
                            second: 'HH:mm:ss'
                        }
                    },
                    title: {
                        display: true,
                        text: '时间'
                    }
                },
                y: {
                    title: {
                        display: true,
                        text: '数量'
                    },
                    beginAtZero: true
                }
            },
            plugins: {
                legend: {
                    display: false
                }
            },
            animation: {
                duration: 0
            }
        }
    });
    
    // 内存使用趋势图
    const memoryCtx = document.getElementById('memory-chart').getContext('2d');
    memoryChart = new Chart(memoryCtx, {
        type: 'line',
        data: {
            datasets: [{
                label: '内存使用 (MB)',
                data: memoryData,
                borderColor: '#e74c3c',
                backgroundColor: 'rgba(231, 76, 60, 0.1)',
                fill: true,
                tension: 0.4
            }]
        },
        options: {
            responsive: true,
            maintainAspectRatio: false,
            scales: {
                x: {
                    type: 'time',
                    time: {
                        unit: 'second',
                        displayFormats: {
                            second: 'HH:mm:ss'
                        }
                    },
                    title: {
                        display: true,
                        text: '时间'
                    }
                },
                y: {
                    title: {
                        display: true,
                        text: '内存 (MB)'
                    },
                    beginAtZero: true
                }
            },
            plugins: {
                legend: {
                    display: false
                }
            },
            animation: {
                duration: 0
            }
        }
    });
    
    // 状态分布饼图
    const stateCtx = document.getElementById('state-chart').getContext('2d');
    stateChart = new Chart(stateCtx, {
        type: 'doughnut',
        data: {
            labels: [],
            datasets: [{
                data: [],
                backgroundColor: [
                    '#27ae60',  // running
                    '#f39c12',  // runnable
                    '#3498db',  // waiting
                    '#e67e22',  // blocked
                    '#95a5a6'   // dead
                ]
            }]
        },
        options: {
            responsive: true,
            maintainAspectRatio: false,
            plugins: {
                legend: {
                    position: 'bottom'
                }
            },
            animation: {
                duration: 0
            }
        }
    });
}

// 设置事件监听器
function setupEventListeners() {
    // 最大显示数量变化
    maxDisplay.addEventListener('change', function() {
        // 下次数据更新时会自动生效
    });
    
    // 窗口大小变化时重新调整图表
    window.addEventListener('resize', function() {
        goroutineChart.resize();
        memoryChart.resize();
        stateChart.resize();
    });
} 