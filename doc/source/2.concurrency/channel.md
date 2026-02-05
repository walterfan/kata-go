# Channel 详解

```{contents} 目录
:depth: 3
```

## Channel 基础

Channel 是 Go 语言中 Goroutine 之间通信的管道，是 CSP (Communicating Sequential Processes) 并发模型的核心。

```{tip}
Go 的并发哲学：**不要通过共享内存来通信，而要通过通信来共享内存**。
```

```{seealso}
延伸阅读：`通过通信来共享内存 <https://www.fanyamin.com/journal/2025-03-26-tong-guo-tong-xin-lai-gong-xiang-nei-cun-er-bu-shi-tong-guo.html>`_ — C++/Java/Go 三种语言实现事件循环的对比。
```

## Channel 类型

### 无缓冲 Channel (Unbuffered)

```go
ch := make(chan int) // 无缓冲 channel

// 发送操作会阻塞，直到有接收者
go func() {
    ch <- 42 // 阻塞，直到有人接收
}()

value := <-ch // 接收
```

### 有缓冲 Channel (Buffered)

```go
ch := make(chan int, 3) // 容量为 3 的 buffered channel

ch <- 1 // 不阻塞
ch <- 2 // 不阻塞
ch <- 3 // 不阻塞
ch <- 4 // 阻塞！缓冲区已满
```

### 单向 Channel

```go
// 只读 channel
func receive(ch <-chan int) {
    value := <-ch
    // ch <- 1 // 编译错误：不能发送
}

// 只写 channel
func send(ch chan<- int) {
    ch <- 1
    // <-ch // 编译错误：不能接收
}
```

## ⚠️ Channel 常见陷阱

### 陷阱 1：向 nil channel 发送/接收

```go
// ❌ 永远阻塞
var ch chan int // nil channel
ch <- 1         // 永远阻塞
<-ch            // 永远阻塞

// ✅ 正确初始化
ch := make(chan int)
```

### 陷阱 2：向已关闭的 channel 发送数据

```go
ch := make(chan int)
close(ch)

// ❌ panic: send on closed channel
ch <- 1

// ✅ 只有发送方应该关闭 channel
```

### 陷阱 3：重复关闭 channel

```go
ch := make(chan int)
close(ch)

// ❌ panic: close of closed channel
close(ch)

// ✅ 使用 sync.Once 确保只关闭一次
var once sync.Once
safeClose := func() {
    once.Do(func() {
        close(ch)
    })
}
```

### 陷阱 4：从已关闭的 channel 读取

```go
ch := make(chan int, 2)
ch <- 1
ch <- 2
close(ch)

// 可以继续读取缓冲区中的数据
fmt.Println(<-ch) // 1
fmt.Println(<-ch) // 2

// 读取已关闭的空 channel 返回零值
fmt.Println(<-ch) // 0 (零值)

// ✅ 正确方式：检查是否关闭
value, ok := <-ch
if !ok {
    fmt.Println("channel closed")
}
```

## Channel 操作总结

| 操作 | nil channel | 正常 channel | 已关闭 channel |
|------|-------------|--------------|----------------|
| 发送 | 永久阻塞 | 阻塞或成功 | **panic** |
| 接收 | 永久阻塞 | 阻塞或成功 | 零值 + false |
| 关闭 | **panic** | 成功 | **panic** |
| len | 0 | 缓冲区元素数 | 0 |
| cap | 0 | 缓冲区容量 | 缓冲区容量 |

## Select 语句

select 用于同时监听多个 channel 操作。

### 基本用法

```go
select {
case v := <-ch1:
    fmt.Println("received from ch1:", v)
case v := <-ch2:
    fmt.Println("received from ch2:", v)
case ch3 <- 42:
    fmt.Println("sent to ch3")
default:
    fmt.Println("no operation ready")
}
```

### 超时处理

```go
select {
case result := <-ch:
    fmt.Println("got result:", result)
case <-time.After(5 * time.Second):
    fmt.Println("timeout!")
}
```

### 非阻塞操作

```go
select {
case v := <-ch:
    fmt.Println("received:", v)
default:
    fmt.Println("channel empty, not blocking")
}
```

### ⚠️ Select 陷阱：随机选择

当多个 case 同时就绪时，select 会**随机**选择一个执行：

```go
ch1 := make(chan int, 1)
ch2 := make(chan int, 1)
ch1 <- 1
ch2 <- 2

// 结果是不确定的！
select {
case v := <-ch1:
    fmt.Println("ch1:", v)
case v := <-ch2:
    fmt.Println("ch2:", v)
}
```

## 常用 Channel 模式

### 模式 1：Done Channel (信号通知)

```go
func worker(done chan struct{}) {
    for {
        select {
        case <-done:
            fmt.Println("worker stopping")
            return
        default:
            // 工作逻辑
            time.Sleep(100 * time.Millisecond)
        }
    }
}

func main() {
    done := make(chan struct{})
    go worker(done)
    
    time.Sleep(time.Second)
    close(done) // 发送停止信号
    time.Sleep(100 * time.Millisecond)
}
```

### 模式 2：Pipeline (管道)

```go
func generator(nums ...int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for _, n := range nums {
            out <- n
        }
    }()
    return out
}

func square(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for n := range in {
            out <- n * n
        }
    }()
    return out
}

func main() {
    // 组装管道
    nums := generator(1, 2, 3, 4, 5)
    squares := square(nums)
    
    for result := range squares {
        fmt.Println(result) // 1, 4, 9, 16, 25
    }
}
```

### 模式 3：Fan-out, Fan-in

```go
// Fan-out: 多个 goroutine 从同一个 channel 读取
func fanOut(in <-chan int, workers int) []<-chan int {
    outs := make([]<-chan int, workers)
    for i := 0; i < workers; i++ {
        outs[i] = worker(in)
    }
    return outs
}

// Fan-in: 将多个 channel 合并为一个
func fanIn(channels ...<-chan int) <-chan int {
    var wg sync.WaitGroup
    out := make(chan int)
    
    output := func(c <-chan int) {
        defer wg.Done()
        for n := range c {
            out <- n
        }
    }
    
    wg.Add(len(channels))
    for _, c := range channels {
        go output(c)
    }
    
    go func() {
        wg.Wait()
        close(out)
    }()
    
    return out
}
```

### 模式 4：Semaphore (信号量)

```go
// 使用 buffered channel 实现信号量
type Semaphore chan struct{}

func NewSemaphore(n int) Semaphore {
    return make(chan struct{}, n)
}

func (s Semaphore) Acquire() {
    s <- struct{}{}
}

func (s Semaphore) Release() {
    <-s
}

// 使用示例
func main() {
    sem := NewSemaphore(3) // 最多 3 个并发
    
    for i := 0; i < 10; i++ {
        sem.Acquire()
        go func(id int) {
            defer sem.Release()
            fmt.Printf("Worker %d running\n", id)
            time.Sleep(time.Second)
        }(i)
    }
}
```

### 模式 5：Or-Done Channel

```go
// 当 done 关闭时，停止从 c 读取
func orDone(done <-chan struct{}, c <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for {
            select {
            case <-done:
                return
            case v, ok := <-c:
                if !ok {
                    return
                }
                select {
                case out <- v:
                case <-done:
                    return
                }
            }
        }
    }()
    return out
}
```

## Channel 性能考虑

### 缓冲区大小选择

```go
// 无缓冲：强同步，适合信号通知
done := make(chan struct{})

// 小缓冲：减少阻塞，适合突发流量
ch := make(chan Task, 10)

// 大缓冲：解耦生产消费速度，注意内存占用
ch := make(chan LargeData, 1000) // 可能占用大量内存
```

### 避免过度使用 Channel

```go
// ❌ 不必要的 channel 使用
func add(a, b int) int {
    result := make(chan int)
    go func() {
        result <- a + b
    }()
    return <-result
}

// ✅ 简单操作直接返回
func add(a, b int) int {
    return a + b
}
```

## 最佳实践

1. **明确 Channel 所有权**：通常由创建者负责关闭
2. **使用单向 Channel**：在函数签名中明确意图
3. **避免在循环中创建 Channel**：考虑复用
4. **合理设置缓冲区大小**：基于实际场景测试
5. **使用 context 代替 done channel**：更标准、功能更丰富

## 参考资源

- [Go Blog: Pipelines and Cancellation](https://go.dev/blog/pipelines)
- [Go Concurrency Patterns](https://www.youtube.com/watch?v=f6kdp27TYZs)
- [Advanced Go Concurrency Patterns](https://www.youtube.com/watch?v=QDDwwePbDtw)
