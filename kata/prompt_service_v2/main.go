package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/walterfan/prompt-service/internal/log"
	"github.com/walterfan/prompt-service/pkg/auth"
	"github.com/walterfan/prompt-service/pkg/authz"
	"github.com/walterfan/prompt-service/pkg/config"
	"github.com/walterfan/prompt-service/pkg/database"
	"github.com/walterfan/prompt-service/pkg/handlers"
	"github.com/walterfan/prompt-service/pkg/metrics"
	"go.uber.org/zap"
)

var (
	logger *zap.Logger
)

func init() {
	var err error
	logger, err = log.InitLogger()
	if err != nil {
		panic(err)
	}

	if err := godotenv.Load(); err != nil {
		zap.L().Warn("No .env file found, using environment variables")
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("error reading config file: %v", err))
	}
}

func startServer(port string) {
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal("Failed to load config", zap.Error(err))
	}

	database.InitDB(cfg.DatabasePath)
	auth.InitJwt(cfg.JwtSecret)
	authz.InitAuthz(cfg.AuthzModelPath)

	metrics.Register()

	r := gin.Default()
	r.Use(metrics.MetricsMiddleware())

	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	r.POST("/login", auth.LoginHandler)

	api := r.Group("/api/v1/prompts")
	api.Use(auth.JwtMiddleware())
	api.Use(authz.CasbinMiddleware())
	{
		api.POST("/", handlers.CreatePrompt)
		api.GET("/:id", handlers.GetPrompt)
		api.PUT("/:id", handlers.UpdatePrompt)
		api.DELETE("/:id", handlers.DeletePrompt)
		api.GET("/", handlers.SearchPrompts)
	}

	// Assuming you've already imported handlers and authz

	userApi := r.Group("/api/v1/users")
	userApi.Use(auth.JwtMiddleware())
	userApi.Use(authz.CasbinMiddleware())
	{
		userApi.POST("/", handlers.CreateUser)
		userApi.GET("/:id", handlers.GetUser)
		userApi.PUT("/:id", handlers.UpdateUser)
		userApi.DELETE("/:id", handlers.DeleteUser)
		userApi.GET("/", handlers.SearchUsers)
	}

	logger.Info("Starting server", zap.String("port", port))
	if err := r.Run(":" + port); err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
	}
}
func main() {
	cmd := &cobra.Command{
		Use:   "prompt-service",
		Short: "Prompt Service",
		Run: func(cmd *cobra.Command, args []string) {
			port, _ := cmd.Flags().GetString("port")
			startServer(port)
		},
	}

	cmd.Flags().StringP("port", "p", "8080", "Port to listen on")

	if err := cmd.Execute(); err != nil {
		logger.Fatal("Command execution failed", zap.Error(err))
	}
}
