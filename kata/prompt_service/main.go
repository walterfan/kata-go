package main

import (
	"github.com/walterfan/prompt-service/pkg/database"
	"github.com/walterfan/prompt-service/pkg/handlers"
	"github.com/walterfan/prompt-service/pkg/metrics"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	
	"go.uber.org/zap"
	"github.com/spf13/cobra"
)

var (
	log *zap.Logger
)

func init() {
	var err error
	log, err = zap.NewProduction()
	if err != nil {
		panic(err)
	}
}

func startServer(port string) {
	database.InitDB()
	metrics.Register()

	r := gin.Default()
	r.Use(metrics.MetricsMiddleware())

	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	api := r.Group("/api/v1/prompts")
	{
		api.POST("/", handlers.CreatePrompt)
		api.GET("/:id", handlers.GetPrompt)
		api.PUT("/:id", handlers.UpdatePrompt)
		api.DELETE("/:id", handlers.DeletePrompt)
		api.GET("/", handlers.SearchPrompts)
	}

	log.Info("Starting server", zap.String("port", port))
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server", zap.Error(err))
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
		log.Fatal("Command execution failed", zap.Error(err))
	}
}