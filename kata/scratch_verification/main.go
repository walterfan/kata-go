package main

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Site represents a site object
type Site struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
	Created  string `json:"created"`
	Modified string `json:"modified"`
}

// EncryptionService handles password encryption/decryption
type EncryptionService struct {
	key []byte
}

// NewEncryptionService creates a new encryption service
func NewEncryptionService() (*EncryptionService, error) {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		// .env file not found, continue with environment variables
	}

	// Load AES key from environment variable (from .env file or system env)
	aesKeyString := os.Getenv("AES_KEY")
	if aesKeyString == "" {
		return nil, fmt.Errorf("AES_KEY environment variable is required (set in .env file or system environment)")
	}

	// Hash the key string with SHA256 to get 32 bytes for AES-256
	hash := sha256.Sum256([]byte(aesKeyString))

	return &EncryptionService{
		key: hash[:],
	}, nil
}

// EncryptPassword encrypts a password using AES-GCM
func (es *EncryptionService) EncryptPassword(password string) (string, error) {
	// Create a new AES cipher block
	block, err := aes.NewCipher(es.key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	// Create a new GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	// Generate a random nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Encrypt the password
	ciphertext := gcm.Seal(nonce, nonce, []byte(password), nil)

	// Encode to base64 for storage
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptPassword decrypts a password using AES-GCM
func (es *EncryptionService) DecryptPassword(encryptedPassword string) (string, error) {
	// Decode from base64
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedPassword)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %w", err)
	}

	// Create a new AES cipher block
	block, err := aes.NewCipher(es.key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	// Create a new GCM mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	// Extract nonce and ciphertext
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	// Decrypt the password
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt: %w", err)
	}

	return string(plaintext), nil
}

// encryptCmd represents the encrypt command
var encryptCmd = &cobra.Command{
	Use:   "encrypt [password]",
	Short: "Encrypt a password using AES-GCM",
	Long:  `Encrypt a password using AES-GCM with the key from AES_KEY environment variable.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		password := args[0]

		// Load .env file if it exists
		if err := godotenv.Load(); err != nil {
			// .env file not found, continue with environment variables
		}

		// Create encryption service
		encryption, err := NewEncryptionService()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		// Encrypt the password
		encrypted, err := encryption.EncryptPassword(password)
		if err != nil {
			fmt.Printf("Error encrypting password: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("%s\n", encrypted)
	},
}

// decryptCmd represents the decrypt command
var decryptCmd = &cobra.Command{
	Use:   "decrypt [encrypted_password]",
	Short: "Decrypt a password using AES-GCM",
	Long:  `Decrypt a password using AES-GCM with the key from AES_KEY environment variable.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		encryptedPassword := args[0]

		// Load .env file if it exists
		if err := godotenv.Load(); err != nil {
			// .env file not found, continue with environment variables
		}

		// Create encryption service
		encryption, err := NewEncryptionService()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		// Decrypt the password
		decrypted, err := encryption.DecryptPassword(encryptedPassword)
		if err != nil {
			fmt.Printf("Error decrypting password: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("%s\n", decrypted)
	},
}

// sitesCmd represents the sites command
var sitesCmd = &cobra.Command{
	Use:   "sites",
	Short: "Manage sites",
	Long:  `Manage sites with CRUD operations via command line.`,
}

// sitesListCmd represents the sites list command
var sitesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all sites",
	Long:  `List all sites stored in the database.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Load .env file if it exists
		if err := godotenv.Load(); err != nil {
			// .env file not found, continue with environment variables
		}

		// Load configuration
		config, err := loadConfig()
		if err != nil {
			fmt.Printf("Failed to load config: %v\n", err)
			os.Exit(1)
		}

		// Setup logger
		logger, err := setupLogger(config)
		if err != nil {
			fmt.Printf("Failed to setup logger: %v\n", err)
			os.Exit(1)
		}
		defer logger.Sync()

		// Create server instance
		server, err := NewServer(config, logger)
		if err != nil {
			fmt.Printf("Failed to create server: %v\n", err)
			os.Exit(1)
		}
		defer server.db.Close()

		// Get all sites
		sites, err := server.GetAllSites()
		if err != nil {
			fmt.Printf("Failed to get sites: %v\n", err)
			os.Exit(1)
		}

		if len(sites) == 0 {
			fmt.Println("No sites found.")
			return
		}

		fmt.Printf("Found %d site(s):\n\n", len(sites))
		for _, site := range sites {
			// Decrypt password for display
			decryptedPassword, err := server.encryption.DecryptPassword(site.Password)
			if err != nil {
				decryptedPassword = "[ENCRYPTED]"
			}

			fmt.Printf("ID: %s\n", site.ID)
			fmt.Printf("Name: %s\n", site.Name)
			fmt.Printf("Username: %s\n", site.Username)
			fmt.Printf("Password: %s\n", decryptedPassword)
			fmt.Printf("Created: %s\n", site.Created)
			fmt.Printf("Modified: %s\n", site.Modified)
			fmt.Println("---")
		}
	},
}

// sitesGetCmd represents the sites get command
var sitesGetCmd = &cobra.Command{
	Use:   "get [id]",
	Short: "Get a specific site by ID",
	Long:  `Get a specific site by its ID.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		siteID := args[0]

		// Load .env file if it exists
		if err := godotenv.Load(); err != nil {
			// .env file not found, continue with environment variables
		}

		// Load configuration
		config, err := loadConfig()
		if err != nil {
			fmt.Printf("Failed to load config: %v\n", err)
			os.Exit(1)
		}

		// Setup logger
		logger, err := setupLogger(config)
		if err != nil {
			fmt.Printf("Failed to setup logger: %v\n", err)
			os.Exit(1)
		}
		defer logger.Sync()

		// Create server instance
		server, err := NewServer(config, logger)
		if err != nil {
			fmt.Printf("Failed to create server: %v\n", err)
			os.Exit(1)
		}
		defer server.db.Close()

		// Get site
		site, err := server.GetSite(siteID)
		if err != nil {
			fmt.Printf("Site not found: %v\n", err)
			os.Exit(1)
		}

		// Decrypt password for display
		decryptedPassword, err := server.encryption.DecryptPassword(site.Password)
		if err != nil {
			decryptedPassword = "[ENCRYPTED]"
		}

		fmt.Printf("ID: %s\n", site.ID)
		fmt.Printf("Name: %s\n", site.Name)
		fmt.Printf("Username: %s\n", site.Username)
		fmt.Printf("Password: %s\n", decryptedPassword)
		fmt.Printf("Created: %s\n", site.Created)
		fmt.Printf("Modified: %s\n", site.Modified)
	},
}

// sitesCreateCmd represents the sites create command
var sitesCreateCmd = &cobra.Command{
	Use:   "create [id] [name] [username] [password]",
	Short: "Create a new site",
	Long:  `Create a new site with the specified ID, name, username, and password.`,
	Args:  cobra.ExactArgs(4),
	Run: func(cmd *cobra.Command, args []string) {
		siteID := args[0]
		siteName := args[1]
		username := args[2]
		password := args[3]

		// Load .env file if it exists
		if err := godotenv.Load(); err != nil {
			// .env file not found, continue with environment variables
		}

		// Load configuration
		config, err := loadConfig()
		if err != nil {
			fmt.Printf("Failed to load config: %v\n", err)
			os.Exit(1)
		}

		// Setup logger
		logger, err := setupLogger(config)
		if err != nil {
			fmt.Printf("Failed to setup logger: %v\n", err)
			os.Exit(1)
		}
		defer logger.Sync()

		// Create server instance
		server, err := NewServer(config, logger)
		if err != nil {
			fmt.Printf("Failed to create server: %v\n", err)
			os.Exit(1)
		}
		defer server.db.Close()

		// Check if site already exists
		_, err = server.GetSite(siteID)
		if err == nil {
			fmt.Printf("Site with ID '%s' already exists\n", siteID)
			os.Exit(1)
		}

		// Encrypt password
		encryptedPassword, err := server.encryption.EncryptPassword(password)
		if err != nil {
			fmt.Printf("Failed to encrypt password: %v\n", err)
			os.Exit(1)
		}

		// Create site
		now := time.Now().Format(time.RFC3339)
		site := &Site{
			ID:       siteID,
			Name:     siteName,
			Username: username,
			Password: encryptedPassword,
			Created:  now,
			Modified: now,
		}

		err = server.SaveSite(site)
		if err != nil {
			fmt.Printf("Failed to save site: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Site created successfully:\n")
		fmt.Printf("ID: %s\n", site.ID)
		fmt.Printf("Name: %s\n", site.Name)
		fmt.Printf("Username: %s\n", site.Username)
		fmt.Printf("Password: %s\n", password)
		fmt.Printf("Created: %s\n", site.Created)
	},
}

// sitesUpdateCmd represents the sites update command
var sitesUpdateCmd = &cobra.Command{
	Use:   "update [id] [name] [username] [password]",
	Short: "Update an existing site",
	Long:  `Update an existing site with the specified ID. Use empty strings for fields you don't want to change.`,
	Args:  cobra.ExactArgs(4),
	Run: func(cmd *cobra.Command, args []string) {
		siteID := args[0]
		siteName := args[1]
		username := args[2]
		password := args[3]

		// Load .env file if it exists
		if err := godotenv.Load(); err != nil {
			// .env file not found, continue with environment variables
		}

		// Load configuration
		config, err := loadConfig()
		if err != nil {
			fmt.Printf("Failed to load config: %v\n", err)
			os.Exit(1)
		}

		// Setup logger
		logger, err := setupLogger(config)
		if err != nil {
			fmt.Printf("Failed to setup logger: %v\n", err)
			os.Exit(1)
		}
		defer logger.Sync()

		// Create server instance
		server, err := NewServer(config, logger)
		if err != nil {
			fmt.Printf("Failed to create server: %v\n", err)
			os.Exit(1)
		}
		defer server.db.Close()

		// Get existing site
		site, err := server.GetSite(siteID)
		if err != nil {
			fmt.Printf("Site not found: %v\n", err)
			os.Exit(1)
		}

		// Update fields if provided
		if siteName != "" {
			site.Name = siteName
		}
		if username != "" {
			site.Username = username
		}
		if password != "" {
			encryptedPassword, err := server.encryption.EncryptPassword(password)
			if err != nil {
				fmt.Printf("Failed to encrypt password: %v\n", err)
				os.Exit(1)
			}
			site.Password = encryptedPassword
		}

		site.Modified = time.Now().Format(time.RFC3339)

		err = server.SaveSite(site)
		if err != nil {
			fmt.Printf("Failed to save site: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Site updated successfully:\n")
		fmt.Printf("ID: %s\n", site.ID)
		fmt.Printf("Name: %s\n", site.Name)
		fmt.Printf("Username: %s\n", site.Username)
		fmt.Printf("Password: %s\n", password)
		fmt.Printf("Modified: %s\n", site.Modified)
	},
}

// sitesDeleteCmd represents the sites delete command
var sitesDeleteCmd = &cobra.Command{
	Use:   "delete [id]",
	Short: "Delete a site by ID",
	Long:  `Delete a site by its ID.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		siteID := args[0]

		// Load .env file if it exists
		if err := godotenv.Load(); err != nil {
			// .env file not found, continue with environment variables
		}

		// Load configuration
		config, err := loadConfig()
		if err != nil {
			fmt.Printf("Failed to load config: %v\n", err)
			os.Exit(1)
		}

		// Setup logger
		logger, err := setupLogger(config)
		if err != nil {
			fmt.Printf("Failed to setup logger: %v\n", err)
			os.Exit(1)
		}
		defer logger.Sync()

		// Create server instance
		server, err := NewServer(config, logger)
		if err != nil {
			fmt.Printf("Failed to create server: %v\n", err)
			os.Exit(1)
		}
		defer server.db.Close()

		// Delete site
		err = server.DeleteSite(siteID)
		if err != nil {
			fmt.Printf("Failed to delete site: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Site '%s' deleted successfully\n", siteID)
	},
}

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the HTTP server",
	Long:  `Start the HTTP server with all endpoints for site management and command execution.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Load configuration
		config, err := loadConfig()
		if err != nil {
			fmt.Printf("Failed to load config: %v\n", err)
			os.Exit(1)
		}

		// Setup logger
		logger, err := setupLogger(config)
		if err != nil {
			fmt.Printf("Failed to setup logger: %v\n", err)
			os.Exit(1)
		}
		defer logger.Sync()

		logger.Info("Application starting")

		// Create and start server
		server, err := NewServer(config, logger)
		if err != nil {
			logger.Error("Failed to create server", zap.Error(err))
			os.Exit(1)
		}

		if err := server.Start(); err != nil {
			logger.Error("Failed to start server", zap.Error(err))
			os.Exit(1)
		}

		// Wait for interrupt signal to gracefully shutdown the server
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		logger.Info("Shutting down server...")

		// Create a context with timeout for shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// Gracefully shutdown the server
		if err := server.Shutdown(ctx); err != nil {
			logger.Error("Server forced to shutdown", zap.Error(err))
		}

		logger.Info("Server exited")
	},
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "scratch-verification",
	Short: "A simple HTTP server with site management and command execution",
	Long: `A simple HTTP server built with Go and the Gin framework that provides:
- Site CRUD operations via /sites endpoints
- Command execution via /commands endpoint
- Password encryption using AES-GCM
- Structured logging with zap and lumberjack
- Graceful shutdown capabilities`,
}

// Server represents the HTTP server
type Server struct {
	config     *viper.Viper
	db         *sql.DB
	router     *gin.Engine
	httpServer *http.Server
	logger     *zap.Logger
	encryption *EncryptionService
}

// NewServer creates a new server instance
func NewServer(config *viper.Viper, logger *zap.Logger) (*Server, error) {
	// Set default data directory if not specified
	dataDir := config.GetString("data.data_dir")
	if dataDir == "" {
		dataDir = "./data"
	}

	// Ensure data directory exists
	os.MkdirAll(dataDir, 0755)

	// Initialize database
	dbPath := dataDir + "/sites.db"
	db, err := initDatabase(dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// Initialize encryption service
	encryption, err := NewEncryptionService()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize encryption service: %w", err)
	}

	return &Server{
		config:     config,
		db:         db,
		logger:     logger,
		encryption: encryption,
	}, nil
}

// initDatabase initializes the SQLite database
func initDatabase(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Create sites table if it doesn't exist
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS sites (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		username TEXT NOT NULL,
		password TEXT NOT NULL,
		created TEXT NOT NULL,
		modified TEXT NOT NULL
	);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	return db, nil
}

// LoadSites loads sites from the database (no-op for database)
func (s *Server) LoadSites() error {
	// Database is already initialized, no need to load
	s.logger.Info("Database initialized, ready to serve sites")
	return nil
}

// SaveSite saves a single site to the database
func (s *Server) SaveSite(site *Site) error {
	query := `
	INSERT OR REPLACE INTO sites (id, name, username, password, created, modified)
	VALUES (?, ?, ?, ?, ?, ?)`

	_, err := s.db.Exec(query, site.ID, site.Name, site.Username, site.Password, site.Created, site.Modified)
	if err != nil {
		return fmt.Errorf("failed to save site: %w", err)
	}

	s.logger.Info("Saved site to database", zap.String("site_id", site.ID))
	return nil
}

// DeleteSite deletes a site from the database
func (s *Server) DeleteSite(id string) error {
	query := `DELETE FROM sites WHERE id = ?`
	result, err := s.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete site: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("site not found")
	}

	s.logger.Info("Deleted site from database", zap.String("site_id", id))
	return nil
}

// GetSite retrieves a site from the database
func (s *Server) GetSite(id string) (*Site, error) {
	query := `SELECT id, name, username, password, created, modified FROM sites WHERE id = ?`
	row := s.db.QueryRow(query, id)

	site := &Site{}
	err := row.Scan(&site.ID, &site.Name, &site.Username, &site.Password, &site.Created, &site.Modified)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("site not found")
		}
		return nil, fmt.Errorf("failed to get site: %w", err)
	}

	return site, nil
}

// GetAllSites retrieves all sites from the database
func (s *Server) GetAllSites() ([]*Site, error) {
	query := `SELECT id, name, username, password, created, modified FROM sites ORDER BY created`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query sites: %w", err)
	}
	defer rows.Close()

	var sites []*Site
	for rows.Next() {
		site := &Site{}
		err := rows.Scan(&site.ID, &site.Name, &site.Username, &site.Password, &site.Created, &site.Modified)
		if err != nil {
			return nil, fmt.Errorf("failed to scan site: %w", err)
		}
		sites = append(sites, site)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return sites, nil
}

// SetupRoutes sets up the HTTP routes
func (s *Server) SetupRoutes() {
	// Create router with custom logger
	gin.SetMode(gin.ReleaseMode)
	s.router = gin.New()

	// Add logging middleware
	s.router.Use(gin.Recovery())
	s.router.Use(s.loggingMiddleware())

	// Health check endpoint
	s.router.GET("/health", s.healthHandler)

	// Sites CRUD endpoints
	sites := s.router.Group("/sites")
	{
		sites.GET("", s.getSites)
		sites.GET("/:id", s.getSite)
		sites.POST("", s.createSite)
		sites.PUT("/:id", s.updateSite)
		sites.DELETE("/:id", s.deleteSite)
	}

	// Commands endpoint
	s.router.POST("/commands", s.executeCommand)
}

// loggingMiddleware creates a middleware for HTTP request/response logging
func (s *Server) loggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Log after request is processed
		latency := time.Since(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		bodySize := c.Writer.Size()

		if raw != "" {
			path = path + "?" + raw
		}

		s.logger.Info("HTTP Request",
			zap.String("method", method),
			zap.String("path", path),
			zap.String("client_ip", clientIP),
			zap.Int("status", statusCode),
			zap.Duration("latency", latency),
			zap.Int("body_size", bodySize),
			zap.String("user_agent", c.Request.UserAgent()),
		)
	}
}

// healthHandler handles the health check endpoint
func (s *Server) healthHandler(c *gin.Context) {
	s.logger.Debug("Health check requested")
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// getSites returns all sites
func (s *Server) getSites(c *gin.Context) {
	sites, err := s.GetAllSites()
	if err != nil {
		s.logger.Error("Failed to get sites", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve sites"})
		return
	}

	// Decrypt passwords for response
	for _, site := range sites {
		decryptedPassword, err := s.encryption.DecryptPassword(site.Password)
		if err != nil {
			s.logger.Error("Failed to decrypt password", zap.String("site_id", site.ID), zap.Error(err))
			site.Password = "[ENCRYPTED]"
		} else {
			site.Password = decryptedPassword
		}
	}

	s.logger.Info("Retrieved all sites", zap.Int("count", len(sites)))
	c.JSON(http.StatusOK, sites)
}

// getSite returns a specific site by ID
func (s *Server) getSite(c *gin.Context) {
	id := c.Param("id")
	site, err := s.GetSite(id)
	if err != nil {
		s.logger.Warn("Site not found", zap.String("site_id", id), zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": "Site not found"})
		return
	}

	// Decrypt password for response
	decryptedPassword, err := s.encryption.DecryptPassword(site.Password)
	if err != nil {
		s.logger.Error("Failed to decrypt password", zap.String("site_id", id), zap.Error(err))
		site.Password = "[ENCRYPTED]"
	} else {
		site.Password = decryptedPassword
	}

	s.logger.Info("Retrieved site", zap.String("site_id", id))
	c.JSON(http.StatusOK, site)
}

// createSite creates a new site
func (s *Server) createSite(c *gin.Context) {
	var site Site
	if err := c.ShouldBindJSON(&site); err != nil {
		s.logger.Error("Failed to bind JSON for site creation", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if site.ID == "" {
		s.logger.Warn("Site creation failed: ID is required")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Site ID is required"})
		return
	}

	// Check if site already exists
	_, err := s.GetSite(site.ID)
	if err == nil {
		s.logger.Warn("Site creation failed: site already exists", zap.String("site_id", site.ID))
		c.JSON(http.StatusConflict, gin.H{"error": "Site already exists"})
		return
	}

	// Store the original password for response
	originalPassword := site.Password

	// Encrypt the password before storing
	encryptedPassword, err := s.encryption.EncryptPassword(site.Password)
	if err != nil {
		s.logger.Error("Failed to encrypt password", zap.String("site_id", site.ID), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt password"})
		return
	}
	site.Password = encryptedPassword

	now := time.Now().Format(time.RFC3339)
	site.Created = now
	site.Modified = now

	if err := s.SaveSite(&site); err != nil {
		s.logger.Error("Failed to save site", zap.String("site_id", site.ID), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save site"})
		return
	}

	// Return the site with original password for response
	siteResponse := Site{
		ID:       site.ID,
		Name:     site.Name,
		Username: site.Username,
		Password: originalPassword,
		Created:  site.Created,
		Modified: site.Modified,
	}

	s.logger.Info("Site created successfully", zap.String("site_id", site.ID), zap.String("name", site.Name))
	c.JSON(http.StatusCreated, &siteResponse)
}

// updateSite updates an existing site
func (s *Server) updateSite(c *gin.Context) {
	id := c.Param("id")
	site, err := s.GetSite(id)
	if err != nil {
		s.logger.Warn("Site update failed: site not found", zap.String("site_id", id), zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": "Site not found"})
		return
	}

	var updateData Site
	if err := c.ShouldBindJSON(&updateData); err != nil {
		s.logger.Error("Failed to bind JSON for site update", zap.String("site_id", id), zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update fields
	if updateData.Name != "" {
		site.Name = updateData.Name
	}
	if updateData.Username != "" {
		site.Username = updateData.Username
	}
	if updateData.Password != "" {
		// Encrypt the new password before storing
		encryptedPassword, err := s.encryption.EncryptPassword(updateData.Password)
		if err != nil {
			s.logger.Error("Failed to encrypt password during update", zap.String("site_id", id), zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt password"})
			return
		}
		site.Password = encryptedPassword
	}

	site.Modified = time.Now().Format(time.RFC3339)

	if err := s.SaveSite(site); err != nil {
		s.logger.Error("Failed to save site after update", zap.String("site_id", id), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save site"})
		return
	}

	// Return the site with decrypted password for response
	decryptedPassword, err := s.encryption.DecryptPassword(site.Password)
	if err != nil {
		s.logger.Error("Failed to decrypt password for response", zap.String("site_id", id), zap.Error(err))
		site.Password = "[ENCRYPTED]"
	} else {
		site.Password = decryptedPassword
	}

	s.logger.Info("Site updated successfully", zap.String("site_id", id))
	c.JSON(http.StatusOK, site)
}

// deleteSite deletes a site
func (s *Server) deleteSite(c *gin.Context) {
	id := c.Param("id")
	err := s.DeleteSite(id)
	if err != nil {
		s.logger.Warn("Site deletion failed: site not found", zap.String("site_id", id), zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": "Site not found"})
		return
	}

	s.logger.Info("Site deleted successfully", zap.String("site_id", id))
	c.JSON(http.StatusOK, gin.H{"message": "Site deleted successfully"})
}

// executeCommand executes a shell command from the whitelist
func (s *Server) executeCommand(c *gin.Context) {
	var request struct {
		Command string `json:"command"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		s.logger.Error("Failed to bind JSON for command execution", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if request.Command == "" {
		s.logger.Warn("Command execution failed: command is required")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Command is required"})
		return
	}

	// Check if command is in whitelist
	whitelist := s.config.GetStringSlice("commands.whitelist")
	allowed := false
	for _, cmd := range whitelist {
		if cmd == request.Command {
			allowed = true
			break
		}
	}

	if !allowed {
		s.logger.Warn("Command execution failed: command not allowed", zap.String("command", request.Command))
		c.JSON(http.StatusForbidden, gin.H{"error": "Command not allowed"})
		return
	}

	// Execute the command
	s.logger.Info("Executing command", zap.String("command", request.Command))
	cmd := exec.Command("sh", "-c", request.Command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		s.logger.Error("Command execution failed",
			zap.String("command", request.Command),
			zap.Error(err),
			zap.String("output", string(output)))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"output":  string(output),
			"command": request.Command,
		})
		return
	}

	s.logger.Info("Command executed successfully",
		zap.String("command", request.Command),
		zap.String("output", string(output)))
	c.JSON(http.StatusOK, gin.H{
		"output":  string(output),
		"command": request.Command,
	})
}

// Start starts the HTTP server
func (s *Server) Start() error {
	// Load existing sites
	if err := s.LoadSites(); err != nil {
		return fmt.Errorf("failed to load sites: %w", err)
	}

	// Setup routes
	s.SetupRoutes()

	// Create HTTP server
	port := s.config.GetInt("server.port")
	s.httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: s.router,
	}

	// Start server in a goroutine
	go func() {
		s.logger.Info("Server starting", zap.Int("port", port))
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Error("Server error", zap.Error(err))
		}
	}()

	return nil
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	if s.httpServer != nil {
		s.logger.Info("Server shutting down")
		err := s.httpServer.Shutdown(ctx)
		if err != nil {
			return err
		}
	}

	// Close database connection
	if s.db != nil {
		s.logger.Info("Closing database connection")
		return s.db.Close()
	}

	return nil
}

// setupLogger creates and configures the zap logger with lumberjack
func setupLogger(config *viper.Viper) (*zap.Logger, error) {
	// Configure lumberjack for log rotation
	logFile := config.GetString("logging.file.filename")
	if logFile == "" {
		logFile = "/data/app.log"
	}

	lumberjackLogger := &lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    config.GetInt("logging.file.max_size"), // MB
		MaxAge:     config.GetInt("logging.file.max_age"),  // days
		MaxBackups: config.GetInt("logging.file.max_backups"),
		Compress:   config.GetBool("logging.file.compress"),
	}

	// Create zap configuration
	zapConfig := zap.NewProductionConfig()

	// Set log level
	level := config.GetString("logging.level")
	if level == "" {
		level = "info"
	}

	var zapLevel zapcore.Level
	if err := zapLevel.UnmarshalText([]byte(level)); err != nil {
		zapLevel = zapcore.InfoLevel
	}
	zapConfig.Level = zap.NewAtomicLevelAt(zapLevel)

	// Configure file output
	zapConfig.OutputPaths = []string{"stdout"}
	zapConfig.ErrorOutputPaths = []string{"stderr"}

	// Create logger
	logger, err := zapConfig.Build()
	if err != nil {
		return nil, err
	}

	// Add lumberjack as a core
	core := zapcore.NewTee(
		logger.Core(),
		zapcore.NewCore(
			zapcore.NewJSONEncoder(zapConfig.EncoderConfig),
			zapcore.AddSync(lumberjackLogger),
			zapLevel,
		),
	)

	return zap.New(core), nil
}

// loadConfig loads configuration using viper
func loadConfig() (*viper.Viper, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/app/")
	viper.AddConfigPath("$HOME/.app")

	// Set default values
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("data.data_dir", "/data")
	viper.SetDefault("logging.level", "info")
	viper.SetDefault("logging.file.filename", "/data/app.log")
	viper.SetDefault("logging.file.max_size", 100)
	viper.SetDefault("logging.file.max_age", 30)
	viper.SetDefault("logging.file.max_backups", 10)
	viper.SetDefault("logging.file.compress", true)

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
		// Config file not found, use defaults
	}

	// Enable environment variable support
	viper.AutomaticEnv()

	return viper.GetViper(), nil
}

func main() {
	// Add subcommands
	rootCmd.AddCommand(encryptCmd)
	rootCmd.AddCommand(decryptCmd)
	rootCmd.AddCommand(serverCmd)

	// Add sites subcommands
	sitesCmd.AddCommand(sitesListCmd)
	sitesCmd.AddCommand(sitesGetCmd)
	sitesCmd.AddCommand(sitesCreateCmd)
	sitesCmd.AddCommand(sitesUpdateCmd)
	sitesCmd.AddCommand(sitesDeleteCmd)
	rootCmd.AddCommand(sitesCmd)

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
