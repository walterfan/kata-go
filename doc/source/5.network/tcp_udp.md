# TCP/UDP 编程

```{contents} 目录
:depth: 3
```

## TCP 编程

### TCP 服务器

```go
package main

import (
    "bufio"
    "fmt"
    "net"
)

func main() {
    listener, err := net.Listen("tcp", ":8080")
    if err != nil {
        panic(err)
    }
    defer listener.Close()
    
    fmt.Println("Server listening on :8080")
    
    for {
        conn, err := listener.Accept()
        if err != nil {
            continue
        }
        go handleConnection(conn)
    }
}

func handleConnection(conn net.Conn) {
    defer conn.Close()
    
    reader := bufio.NewReader(conn)
    for {
        message, err := reader.ReadString('\n')
        if err != nil {
            return
        }
        
        fmt.Printf("Received: %s", message)
        conn.Write([]byte("Echo: " + message))
    }
}
```

### TCP 客户端

```go
func main() {
    conn, err := net.Dial("tcp", "localhost:8080")
    if err != nil {
        panic(err)
    }
    defer conn.Close()
    
    // 发送数据
    conn.Write([]byte("Hello, Server!\n"))
    
    // 接收响应
    response := make([]byte, 1024)
    n, _ := conn.Read(response)
    fmt.Printf("Response: %s", response[:n])
}
```

### ⚠️ TCP 陷阱

#### 陷阱 1：粘包问题

```go
// ❌ 直接读取，可能读到不完整或多个消息
func handleConn(conn net.Conn) {
    buf := make([]byte, 1024)
    n, _ := conn.Read(buf)
    process(buf[:n])  // 可能是半个消息或多个消息
}

// ✅ 使用长度前缀协议
func handleConn(conn net.Conn) {
    reader := bufio.NewReader(conn)
    for {
        // 读取 4 字节长度
        lenBuf := make([]byte, 4)
        io.ReadFull(reader, lenBuf)
        length := binary.BigEndian.Uint32(lenBuf)
        
        // 读取消息体
        data := make([]byte, length)
        io.ReadFull(reader, data)
        
        process(data)
    }
}
```

#### 陷阱 2：忘记设置超时

```go
// ❌ 无超时，可能永久阻塞
conn, _ := net.Dial("tcp", "server:8080")
conn.Read(buf)  // 可能永久阻塞

// ✅ 设置超时
conn, _ := net.DialTimeout("tcp", "server:8080", 5*time.Second)
conn.SetReadDeadline(time.Now().Add(10 * time.Second))
conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
```

#### 陷阱 3：未处理半关闭

```go
// TCP 连接可以半关闭（一端关闭写，但仍可读）
func handleConn(conn net.Conn) {
    // 读取客户端数据
    data, _ := io.ReadAll(conn)
    
    // 关闭读端，但仍可写
    conn.(*net.TCPConn).CloseRead()
    
    // 发送响应
    conn.Write([]byte("Response"))
    conn.Close()
}
```

## UDP 编程

### UDP 服务器

```go
func main() {
    addr, _ := net.ResolveUDPAddr("udp", ":8080")
    conn, err := net.ListenUDP("udp", addr)
    if err != nil {
        panic(err)
    }
    defer conn.Close()
    
    buf := make([]byte, 1024)
    for {
        n, remoteAddr, err := conn.ReadFromUDP(buf)
        if err != nil {
            continue
        }
        
        fmt.Printf("Received from %v: %s\n", remoteAddr, buf[:n])
        
        // 响应
        conn.WriteToUDP([]byte("Echo: "+string(buf[:n])), remoteAddr)
    }
}
```

### UDP 客户端

```go
func main() {
    addr, _ := net.ResolveUDPAddr("udp", "localhost:8080")
    conn, err := net.DialUDP("udp", nil, addr)
    if err != nil {
        panic(err)
    }
    defer conn.Close()
    
    conn.Write([]byte("Hello, UDP!"))
    
    buf := make([]byte, 1024)
    n, _ := conn.Read(buf)
    fmt.Printf("Response: %s\n", buf[:n])
}
```

### ⚠️ UDP 陷阱

#### 陷阱 1：假设消息一定到达

```go
// ❌ UDP 不保证消息送达
conn.Write(importantData)
// 数据可能丢失！

// ✅ 实现确认机制或使用 TCP
func sendWithRetry(conn *net.UDPConn, data []byte, maxRetries int) error {
    for i := 0; i < maxRetries; i++ {
        conn.Write(data)
        conn.SetReadDeadline(time.Now().Add(time.Second))
        
        ack := make([]byte, 4)
        _, err := conn.Read(ack)
        if err == nil && string(ack) == "ACK" {
            return nil
        }
    }
    return errors.New("send failed after retries")
}
```

#### 陷阱 2：假设消息顺序

```go
// UDP 不保证顺序，需要自己处理
type Message struct {
    SeqNum  uint32
    Payload []byte
}

// 接收端需要重排序
```

## 连接池

```go
type ConnPool struct {
    mu    sync.Mutex
    conns chan net.Conn
    addr  string
}

func NewConnPool(addr string, size int) *ConnPool {
    pool := &ConnPool{
        conns: make(chan net.Conn, size),
        addr:  addr,
    }
    
    // 预创建连接
    for i := 0; i < size; i++ {
        conn, _ := net.Dial("tcp", addr)
        pool.conns <- conn
    }
    
    return pool
}

func (p *ConnPool) Get() net.Conn {
    return <-p.conns
}

func (p *ConnPool) Put(conn net.Conn) {
    p.conns <- conn
}
```

## 参考资源

- [net Package](https://pkg.go.dev/net)
- [Network Programming with Go](https://ipfs.io/ipfs/QmfYeDhGH9bZzihBUDEQbCbTc5k5FZKURMUoUvfmc27BwL/socket/tcp_sockets.html)
