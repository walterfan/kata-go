# gRPC

```{contents} 目录
:depth: 3
```

## gRPC 概述

gRPC 是 Google 开发的高性能 RPC 框架，基于 HTTP/2 和 Protocol Buffers。

| 特性 | 描述 |
|------|------|
| 序列化 | Protocol Buffers（二进制） |
| 传输 | HTTP/2（多路复用） |
| 调用模式 | Unary、Server Stream、Client Stream、Bidirectional |

## 安装

```bash
# 安装 protoc 编译器
brew install protobuf

# 安装 Go 插件
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

## 定义服务

### Proto 文件

```protobuf
// user.proto
syntax = "proto3";

package user;
option go_package = "./pb";

message User {
    int64 id = 1;
    string name = 2;
    string email = 3;
}

message GetUserRequest {
    int64 id = 1;
}

message ListUsersRequest {
    int32 page = 1;
    int32 page_size = 2;
}

service UserService {
    // Unary RPC
    rpc GetUser(GetUserRequest) returns (User);
    
    // Server streaming
    rpc ListUsers(ListUsersRequest) returns (stream User);
    
    // Client streaming
    rpc CreateUsers(stream User) returns (CreateUsersResponse);
    
    // Bidirectional streaming
    rpc Chat(stream ChatMessage) returns (stream ChatMessage);
}
```

### 生成代码

```bash
protoc --go_out=. --go-grpc_out=. user.proto
```

## 服务端实现

```go
package main

import (
    "context"
    "net"
    
    "google.golang.org/grpc"
    pb "your/package/pb"
)

type userServer struct {
    pb.UnimplementedUserServiceServer
}

func (s *userServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error) {
    // 实现获取用户逻辑
    return &pb.User{
        Id:    req.Id,
        Name:  "Alice",
        Email: "alice@example.com",
    }, nil
}

func (s *userServer) ListUsers(req *pb.ListUsersRequest, stream pb.UserService_ListUsersServer) error {
    users := getUsersFromDB(req.Page, req.PageSize)
    for _, user := range users {
        if err := stream.Send(user); err != nil {
            return err
        }
    }
    return nil
}

func main() {
    lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatal(err)
    }
    
    server := grpc.NewServer()
    pb.RegisterUserServiceServer(server, &userServer{})
    
    if err := server.Serve(lis); err != nil {
        log.Fatal(err)
    }
}
```

## 客户端实现

```go
package main

import (
    "context"
    "log"
    
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
    pb "your/package/pb"
)

func main() {
    conn, err := grpc.Dial("localhost:50051", 
        grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()
    
    client := pb.NewUserServiceClient(conn)
    
    // Unary 调用
    user, err := client.GetUser(context.Background(), &pb.GetUserRequest{Id: 1})
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("User: %v", user)
    
    // Server streaming
    stream, err := client.ListUsers(context.Background(), &pb.ListUsersRequest{Page: 1, PageSize: 10})
    if err != nil {
        log.Fatal(err)
    }
    for {
        user, err := stream.Recv()
        if err == io.EOF {
            break
        }
        if err != nil {
            log.Fatal(err)
        }
        log.Printf("User: %v", user)
    }
}
```

## ⚠️ 常见陷阱

### 陷阱 1：忘记处理 Context

```go
// ❌ 忽略 context 取消
func (s *server) LongOperation(ctx context.Context, req *pb.Request) (*pb.Response, error) {
    // 长时间操作，不检查 context
    for i := 0; i < 100; i++ {
        doWork()
    }
    return &pb.Response{}, nil
}

// ✅ 检查 context
func (s *server) LongOperation(ctx context.Context, req *pb.Request) (*pb.Response, error) {
    for i := 0; i < 100; i++ {
        select {
        case <-ctx.Done():
            return nil, ctx.Err()
        default:
            doWork()
        }
    }
    return &pb.Response{}, nil
}
```

### 陷阱 2：客户端连接泄漏

```go
// ❌ 每次请求创建新连接
func getUser(id int64) (*pb.User, error) {
    conn, _ := grpc.Dial("localhost:50051", ...)
    // 忘记 defer conn.Close()
    client := pb.NewUserServiceClient(conn)
    return client.GetUser(context.Background(), &pb.GetUserRequest{Id: id})
}

// ✅ 复用连接
var conn *grpc.ClientConn

func init() {
    var err error
    conn, err = grpc.Dial("localhost:50051", ...)
    if err != nil {
        log.Fatal(err)
    }
}

func getUser(id int64) (*pb.User, error) {
    client := pb.NewUserServiceClient(conn)
    return client.GetUser(context.Background(), &pb.GetUserRequest{Id: id})
}
```

### 陷阱 3：流未正确关闭

```go
// ❌ Server stream 未处理错误
func (s *server) ListUsers(req *pb.Request, stream pb.Service_ListUsersServer) error {
    for _, user := range users {
        stream.Send(user)  // 忽略错误
    }
    return nil
}

// ✅ 处理发送错误
func (s *server) ListUsers(req *pb.Request, stream pb.Service_ListUsersServer) error {
    for _, user := range users {
        if err := stream.Send(user); err != nil {
            return err
        }
    }
    return nil
}
```

## 拦截器 (Interceptor)

### 服务端拦截器

```go
func loggingInterceptor(
    ctx context.Context,
    req interface{},
    info *grpc.UnaryServerInfo,
    handler grpc.UnaryHandler,
) (interface{}, error) {
    start := time.Now()
    resp, err := handler(ctx, req)
    log.Printf("Method: %s, Duration: %v, Error: %v",
        info.FullMethod, time.Since(start), err)
    return resp, err
}

server := grpc.NewServer(
    grpc.UnaryInterceptor(loggingInterceptor),
)
```

### 客户端拦截器

```go
func clientInterceptor(
    ctx context.Context,
    method string,
    req, reply interface{},
    cc *grpc.ClientConn,
    invoker grpc.UnaryInvoker,
    opts ...grpc.CallOption,
) error {
    // 添加 metadata
    ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "bearer token")
    return invoker(ctx, method, req, reply, cc, opts...)
}

conn, _ := grpc.Dial("localhost:50051",
    grpc.WithUnaryInterceptor(clientInterceptor),
)
```

## 错误处理

```go
import "google.golang.org/grpc/status"
import "google.golang.org/grpc/codes"

// 服务端返回错误
func (s *server) GetUser(ctx context.Context, req *pb.Request) (*pb.User, error) {
    user, err := s.db.GetUser(req.Id)
    if err != nil {
        if errors.Is(err, ErrNotFound) {
            return nil, status.Error(codes.NotFound, "user not found")
        }
        return nil, status.Error(codes.Internal, "internal error")
    }
    return user, nil
}

// 客户端处理错误
resp, err := client.GetUser(ctx, req)
if err != nil {
    st, ok := status.FromError(err)
    if ok {
        switch st.Code() {
        case codes.NotFound:
            // 处理 not found
        case codes.Internal:
            // 处理内部错误
        }
    }
}
```

## 参考资源

- [gRPC Go](https://grpc.io/docs/languages/go/)
- [Protocol Buffers](https://protobuf.dev/)
