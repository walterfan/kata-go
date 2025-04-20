package agent

/*
- listen command from boss:
	 http, grpc and unix socket
	 command has sync and async properties, ack and noack properties
- execute command and return result to boss
- ack command and report command execute result to boss later
*/
import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/robfig/cron/v3"
	"google.golang.org/grpc"

)

// ----- gRPC Service Implementation -----

type CommandRequest struct {
	Name      string
	Type      string
	Seq       uint32
	TrackId   string
	Timestamp int64
	From      string
	To        string
	Message   string `json:"message"`
}

type CommandResponse struct {
	CommandRequest
	Code int
	Desc string
}

type CommandHandler interface {
	OnCommand(cmd CommandRequest) (ret CommandResponse, err error)
}

// send command and handle command
type LazyAgent interface {
	SendCommand(cmd CommandRequest) (ret int, err error)
	OnCommand(cmd CommandRequest) (ret CommandResponse, err error)
	RegisterHandler(handler CommandHandler) (ret int, err error)
}

// Define a server struct for your gRPC service.
type myGRPCServer struct {
	// pb.UnimplementedYourServiceServer // Uncomment if using newer proto versions
}

// Example gRPC method (replace with your actual RPC methods).
// func (s *myGRPCServer) SomeMethod(ctx context.Context, req *pb.YourRequest) (*pb.YourResponse, error) {
//     // Implement your business logic here.
//     return &pb.YourResponse{Message: "Response from gRPC service"}, nil
// }

// ----- HTTP Handler -----

// A simple HTTP handler.
func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello from HTTP!")
}

func main() {
	// Create a context that cancels on system interrupt or SIGTERM.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// ----- Setup HTTP Server -----
	httpMux := http.NewServeMux()
	httpMux.HandleFunc("/hello", helloHandler)
	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: httpMux,
	}

	// ----- Setup gRPC Server -----
	grpcServer := grpc.NewServer()
	// Register your gRPC service here. For example:
	// pb.RegisterYourServiceServer(grpcServer, &myGRPCServer{})

	// ----- Setup Scheduled Jobs using cron -----
	c := cron.New()
	// Schedule a job to run every minute.
	_, err := c.AddFunc("@every 1m", func() {
		log.Println("Executing scheduled job at", time.Now())
		// Insert your scheduled task logic here.
	})
	if err != nil {
		log.Fatalf("Error scheduling job: %v", err)
	}
	c.Start()
	defer c.Stop()

	// ----- Run HTTP and gRPC Servers concurrently -----
	// Start HTTP server.
	go func() {
		log.Println("HTTP server listening on :8080")
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	// Start gRPC server.
	go func() {
		lis, err := net.Listen("tcp", ":9090")
		if err != nil {
			log.Fatalf("Failed to listen on :9090: %v", err)
		}
		log.Println("gRPC server listening on :9090")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("gRPC server error: %v", err)
		}
	}()

	// ----- Graceful Shutdown -----
	<-ctx.Done() // Wait for interrupt signal.
	log.Println("Shutdown signal received, shutting down gracefully...")

	// Shutdown HTTP server with a timeout.
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("HTTP shutdown error: %v", err)
	}

	// Gracefully stop the gRPC server.
	grpcServer.GracefulStop()

	log.Println("All servers stopped. Exiting.")
}
