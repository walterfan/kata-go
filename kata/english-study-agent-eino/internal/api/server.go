package api

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/walterfan/english-agent/internal/agent"
	"github.com/walterfan/english-agent/internal/config"
	"github.com/walterfan/english-agent/internal/logger"
	"github.com/walterfan/english-agent/internal/rss"
	"github.com/walterfan/english-agent/internal/storage"
	"go.uber.org/zap"
)

type Server struct {
	agent      *agent.Agent
	router     *gin.Engine
	rssFetcher *rss.Fetcher
}

type ChatRequest struct {
	Text string `json:"text"`
	Task string `json:"task"`
}

type ChatResponse struct {
	Result string `json:"result"`
}

func NewServer(ctx context.Context) (*Server, error) {
	a, err := agent.NewAgent(ctx)
	if err != nil {
		return nil, err
	}

	r := gin.Default()
	s := &Server{
		agent:      a,
		router:     r,
		rssFetcher: rss.NewFetcher(),
	}

	s.setupRoutes()
	return s, nil
}

func (s *Server) setupRoutes() {
	api := s.router.Group("/api")
	{
		api.POST("/chat", s.handleChat)
		api.POST("/chat/stream", s.handleChatStream) // New: streaming endpoint
		api.POST("/explain", s.handleExplain)
		api.GET("/feeds", s.handleFeeds)
		api.GET("/rss-sources", s.handleRssSources)
		api.GET("/review", s.handleReview)
		api.POST("/fetch-url", s.handleFetchURL) // New: fetch article from URL
		
		// Custom feed management
		api.GET("/custom-feeds", s.handleGetCustomFeeds)
		api.POST("/custom-feeds", s.handleAddCustomFeed)
		api.PUT("/custom-feeds/:id", s.handleUpdateCustomFeed)
		api.DELETE("/custom-feeds/:id", s.handleDeleteCustomFeed)
	}
}

func (s *Server) handleChat(c *gin.Context) {
	var req ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logger.Log.Info("chat request", zap.String("task", req.Task))

	result, err := s.agent.Run(c.Request.Context(), req.Text, req.Task)
	if err != nil {
		logger.Log.Error("agent run failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ChatResponse{Result: result})
}

// handleChatStream handles streaming chat requests using Server-Sent Events
func (s *Server) handleChatStream(c *gin.Context) {
	var req ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logger.Log.Info("chat stream request", zap.String("task", req.Task))

	// Set headers for SSE
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")

	// Get streaming channels
	contentCh, errCh := s.agent.RunStream(c.Request.Context(), req.Text, req.Task)

	// Stream the response - write SSE format directly to avoid HTML escaping
	c.Stream(func(w io.Writer) bool {
		select {
		case content, ok := <-contentCh:
			if !ok {
				// Channel closed, send done event
				fmt.Fprintf(w, "event: done\ndata: \n\n")
				return false
			}
			// Send content chunk without HTML escaping
			// Replace newlines in content to preserve SSE format
			escapedContent := strings.ReplaceAll(content, "\n", "\\n")
			fmt.Fprintf(w, "event: message\ndata: %s\n\n", escapedContent)
			return true
		case err, ok := <-errCh:
			if ok && err != nil {
				logger.Log.Error("stream error", zap.Error(err))
				fmt.Fprintf(w, "event: error\ndata: %s\n\n", err.Error())
			}
			return false
		case <-c.Request.Context().Done():
			return false
		}
	})
}

func (s *Server) handleExplain(c *gin.Context) {
	var req ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	task := "Explain the meaning and list useful vocabulary"
	if req.Task != "" {
		task = req.Task
	}

	result, err := s.agent.Run(c.Request.Context(), req.Text, task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ChatResponse{Result: result})
}

func (s *Server) handleFeeds(c *gin.Context) {
	// Get optional source filter from query parameter
	source := c.Query("source")
	
	articles, err := s.rssFetcher.FetchFromSource(c.Request.Context(), source)
	if err != nil {
		logger.Log.Error("fetch feeds failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"articles": articles})
}

func (s *Server) handleRssSources(c *gin.Context) {
	// Get config-based sources
	configSources := s.rssFetcher.GetSources()
	
	// Get custom feeds from database
	customFeeds, err := storage.GetEnabledCustomFeeds()
	if err != nil {
		logger.Log.Warn("failed to get custom feeds", zap.Error(err))
	}
	
	// Combine both sources
	var allSources []config.FeedConfig
	for _, src := range configSources {
		allSources = append(allSources, src)
	}
	for _, cf := range customFeeds {
		allSources = append(allSources, config.FeedConfig{
			Title:    cf.Title,
			URL:      cf.URL,
			Category: cf.Category,
		})
	}
	
	c.JSON(http.StatusOK, gin.H{"sources": allSources})
}

func (s *Server) handleReview(c *gin.Context) {
	items, err := storage.GetLearningItems()
	if err != nil {
		logger.Log.Error("fetch review items failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

// FetchURLRequest represents a request to fetch article from URL
type FetchURLRequest struct {
	URL string `json:"url"`
}

func (s *Server) handleFetchURL(c *gin.Context) {
	var req FetchURLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.URL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "url is required"})
		return
	}

	logger.Log.Info("fetching article from URL", zap.String("url", req.URL))

	article, err := s.rssFetcher.FetchFromURL(c.Request.Context(), req.URL)
	if err != nil {
		logger.Log.Error("fetch URL failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"title":   article.Title,
		"content": article.Content,
		"url":     article.URL,
	})
}

// Custom feed handlers

type CustomFeedRequest struct {
	Title    string `json:"title"`
	URL      string `json:"url"`
	Category string `json:"category"`
	Enabled  bool   `json:"enabled"`
}

func (s *Server) handleGetCustomFeeds(c *gin.Context) {
	feeds, err := storage.GetCustomFeeds()
	if err != nil {
		logger.Log.Error("get custom feeds failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"feeds": feeds})
}

func (s *Server) handleAddCustomFeed(c *gin.Context) {
	var req CustomFeedRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Title == "" || req.URL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title and url are required"})
		return
	}

	if err := storage.AddCustomFeed(req.Title, req.URL, req.Category); err != nil {
		logger.Log.Error("add custom feed failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "feed added"})
}

func (s *Server) handleUpdateCustomFeed(c *gin.Context) {
	id := c.Param("id")
	var feedID int
	if _, err := fmt.Sscanf(id, "%d", &feedID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req CustomFeedRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := storage.UpdateCustomFeed(feedID, req.Title, req.URL, req.Category, req.Enabled); err != nil {
		logger.Log.Error("update custom feed failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "feed updated"})
}

func (s *Server) handleDeleteCustomFeed(c *gin.Context) {
	id := c.Param("id")
	var feedID int
	if _, err := fmt.Sscanf(id, "%d", &feedID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := storage.DeleteCustomFeed(feedID); err != nil {
		logger.Log.Error("delete custom feed failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "feed deleted"})
}

func (s *Server) Run() error {
	port := config.Get().Server.Port
	return s.router.Run(":" + port)
}
