package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
)

// User represents a user object
type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Age      int    `json:"age"`
	Created  string `json:"created"`
	Modified string `json:"modified"`
}

// Config represents the application configuration
type Config struct {
	Server struct {
		Port int `yaml:"port"`
	} `yaml:"server"`
	Data struct {
		DataDir string `yaml:"data_dir"`
	} `yaml:"data"`
	Commands struct {
		Whitelist []string `yaml:"whitelist"`
	} `yaml:"commands"`
}

// Server represents the HTTP server
type Server struct {
	config     *Config
	users      map[string]*User
	usersFile  string
	router     *gin.Engine
	httpServer *http.Server
}

// NewServer creates a new server instance
func NewServer(config *Config) *Server {
	// Set default data directory if not specified
	dataDir := config.Data.DataDir
	if dataDir == "" {
		dataDir = "/data"
	}

	// Ensure data directory exists
	os.MkdirAll(dataDir, 0755)

	return &Server{
		config:    config,
		users:     make(map[string]*User),
		usersFile: dataDir + "/users.json",
	}
}

// LoadUsers loads users from the JSON file
func (s *Server) LoadUsers() error {
	file, err := os.Open(s.usersFile)
	if err != nil {
		if os.IsNotExist(err) {
			// File doesn't exist, start with empty users
			return nil
		}
		return err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	if len(data) == 0 {
		return nil
	}

	return json.Unmarshal(data, &s.users)
}

// SaveUsers saves users to the JSON file
func (s *Server) SaveUsers() error {
	data, err := json.MarshalIndent(s.users, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.usersFile, data, 0644)
}

// SetupRoutes sets up the HTTP routes
func (s *Server) SetupRoutes() {
	s.router = gin.Default()

	// Health check endpoint
	s.router.GET("/health", s.healthHandler)

	// Users CRUD endpoints
	users := s.router.Group("/users")
	{
		users.GET("", s.getUsers)
		users.GET("/:id", s.getUser)
		users.POST("", s.createUser)
		users.PUT("/:id", s.updateUser)
		users.DELETE("/:id", s.deleteUser)
	}

	// Commands endpoint
	s.router.POST("/commands", s.executeCommand)
}

// healthHandler handles the health check endpoint
func (s *Server) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// getUsers returns all users
func (s *Server) getUsers(c *gin.Context) {
	users := make([]*User, 0, len(s.users))
	for _, user := range s.users {
		users = append(users, user)
	}
	c.JSON(http.StatusOK, users)
}

// getUser returns a specific user by ID
func (s *Server) getUser(c *gin.Context) {
	id := c.Param("id")
	user, exists := s.users[id]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

// createUser creates a new user
func (s *Server) createUser(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if user.ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	if _, exists := s.users[user.ID]; exists {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}

	now := time.Now().Format(time.RFC3339)
	user.Created = now
	user.Modified = now

	s.users[user.ID] = &user

	if err := s.SaveUsers(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save user"})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// updateUser updates an existing user
func (s *Server) updateUser(c *gin.Context) {
	id := c.Param("id")
	user, exists := s.users[id]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var updateData User
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update fields
	if updateData.Name != "" {
		user.Name = updateData.Name
	}
	if updateData.Email != "" {
		user.Email = updateData.Email
	}
	if updateData.Age > 0 {
		user.Age = updateData.Age
	}

	user.Modified = time.Now().Format(time.RFC3339)

	if err := s.SaveUsers(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save user"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// deleteUser deletes a user
func (s *Server) deleteUser(c *gin.Context) {
	id := c.Param("id")
	if _, exists := s.users[id]; !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	delete(s.users, id)

	if err := s.SaveUsers(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save users"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// executeCommand executes a shell command from the whitelist
func (s *Server) executeCommand(c *gin.Context) {
	var request struct {
		Command string `json:"command"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if request.Command == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Command is required"})
		return
	}

	// Check if command is in whitelist
	allowed := false
	for _, cmd := range s.config.Commands.Whitelist {
		if cmd == request.Command {
			allowed = true
			break
		}
	}

	if !allowed {
		c.JSON(http.StatusForbidden, gin.H{"error": "Command not allowed"})
		return
	}

	// Execute the command
	cmd := exec.Command("sh", "-c", request.Command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"output":  string(output),
			"command": request.Command,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"output":  string(output),
		"command": request.Command,
	})
}

// Start starts the HTTP server
func (s *Server) Start() error {
	// Load existing users
	if err := s.LoadUsers(); err != nil {
		return fmt.Errorf("failed to load users: %w", err)
	}

	// Setup routes
	s.SetupRoutes()

	// Create HTTP server
	s.httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%d", s.config.Server.Port),
		Handler: s.router,
	}

	// Start server in a goroutine
	go func() {
		fmt.Printf("Server starting on port %d\n", s.config.Server.Port)
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Server error: %v\n", err)
		}
	}()

	return nil
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	if s.httpServer != nil {
		return s.httpServer.Shutdown(ctx)
	}
	return nil
}

// LoadConfig loads configuration from YAML file
func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func main() {
	// Load configuration
	config, err := LoadConfig("config.yaml")
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Create and start server
	server := NewServer(config)
	if err := server.Start(); err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
		os.Exit(1)
	}

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("Shutting down server...")

	// Create a context with timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Gracefully shutdown the server
	if err := server.Shutdown(ctx); err != nil {
		fmt.Printf("Server forced to shutdown: %v\n", err)
	}

	fmt.Println("Server exited")
}
