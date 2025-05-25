# Lazy-Agent

## build
### Install Dependencies

You need to install the required Go packages:

```bash
go mod init lazy-agent
go get github.com/robfig/cron/v3
go get google.golang.org/grpc
go get gopkg.in/yaml.v2
```


### 1. Define Protobuf Service
Create `agent.proto`:

```proto
syntax = "proto3";
package agent;
option go_package = "./;agent";

import "google/api/annotations.proto";

service JobService {
  rpc AddJob(AddJobRequest) returns (JobResponse) {
    option (google.api.http) = {
      post: "/v1/jobs"
      body: "*"
    };
  }

  rpc ListJobs(Empty) returns (ListJobsResponse) {
    option (google.api.http) = {
      get: "/v1/jobs"
    };
  }

  rpc RemoveJob(RemoveJobRequest) returns (JobResponse) {
    option (google.api.http) = {
      delete: "/v1/jobs/{id}"
    };
  }
}

// Message definitions from previous thought process
```

Generate code:
```bash
protoc -I. --go_out=. --go-grpc_out=. \
  --grpc-gateway_out=. --openapiv2_out=. \
  agent.proto
```

### 2. Implement gRPC Service
```go
// server.go
package main

import (
	"context"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type JobStore struct {
	mu      sync.RWMutex
	jobs    map[string]cron.EntryID
	cron    *cron.Cron
	entries map[cron.EntryID]string // Track entry IDs to job IDs
}

func NewJobStore() *JobStore {
	return &JobStore{
		jobs:    make(map[string]cron.EntryID),
		cron:    cron.New(cron.WithSeconds()),
		entries: make(map[cron.EntryID]string),
	}
}

func (s *JobStore) AddJob(id, schedule, cmd string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.jobs[id]; exists {
		return status.Error(codes.AlreadyExists, "job exists")
	}

	entryID, err := s.cron.AddFunc(schedule, func() {
		// Execute command
		// Implement your command execution logic here
	})
	if err != nil {
		return status.Error(codes.InvalidArgument, err.Error())
	}

	s.jobs[id] = entryID
	s.entries[entryID] = id
	s.cron.Start()
	return nil
}

// Implement ListJobs and RemoveJob similarly
```

### 3. Set Up Servers
```go
// main.go
package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

func main() {
	store := NewJobStore()
	grpcServer := grpc.NewServer()
	RegisterJobServiceServer(grpcServer, &jobServer{store: store})

	ctx, cancel := context.WithCancel(context.Background())
	gwMux := runtime.NewServeMux()
	RegisterJobServiceHandlerServer(ctx, gwMux, &jobServer{store: store})

	g := new(errgroup.Group)
	
	// gRPC Server
	g.Go(func() error {
		lis, _ := net.Listen("tcp", ":50051")
		return grpcServer.Serve(lis)
	})
	
	// HTTP Gateway
	g.Go(func() error {
		return http.ListenAndServe(":8080", gwMux)
	})
	
	// Signal handling
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
	cancel()
	
	grpcServer.GracefulStop()
	store.cron.Stop()
	_ = g.Wait()
}
```

### 4. Job Execution
Enhance the cron job function to execute commands safely:

```go
func executeCommand(cmd string) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Split cmd into command and args
	// Example: cmd = "/bin/sh -c 'echo hello'"
	// Implement proper command parsing based on your needs
	c := exec.CommandContext(ctx, "/bin/sh", "-c", cmd)
	output, err := c.CombinedOutput()
	
	// Handle output/errors
	log.Printf("Job output: %s\nError: %v", output, err)
}
```

### 5. Full Implementation Structure
```
.
├── agent.proto
├── generated
│   ├── agent.pb.go
│   ├── agent_grpc.pb.go
│   └── agent.pb.gw.go
├── go.mod
├── go.sum
└── main.go
```

### 6. Testing
Use `curl` for HTTP API:
```bash
# Add job
curl -X POST http://localhost:8080/v1/jobs \
  -d '{"id":"job1", "schedule":"@every 1m", "command":"date > /tmp/output"}'

# List jobs
curl http://localhost:8080/v1/jobs
```

Use `grpcurl` for gRPC:
```bash
grpcurl -plaintext -d '{"id":"job1"}' localhost:50051 agent.JobService.RemoveJob
```

### Key Considerations:
1. **Concurrency**: Use mutexes to protect shared state
2. **Security**: Validate commands and schedules carefully
3. **Persistence**: Add database integration for job storage
4. **Observability**: Add logging and metrics
5. **Graceful Shutdown**: Handle in-progress jobs during shutdown

This implementation provides a foundation that you can extend with additional features like job status tracking, retries, and output storage.