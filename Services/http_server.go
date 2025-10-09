package Services

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jupitermetalabs/geth-facade/Types"
)

// HTTPServer provides HTTP JSON-RPC server using Gin framework
// //future: May add rate limiting, authentication, and metrics
// //debugging: Includes request logging and error handling
type HTTPServer struct {
	h *Handlers
}

func NewHTTPServer(h *Handlers) *HTTPServer { return &HTTPServer{h: h} }

func (s *HTTPServer) Serve(addr string) error {
	// //debugging: Set Gin to release mode for production
	gin.SetMode(gin.ReleaseMode)

	// Create Gin router
	r := gin.New()

	// //debugging: Add middleware for logging and recovery
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Configure CORS
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	r.Use(cors.New(config))

	// Add health check endpoints
	r.GET("/health", s.healthCheck)
	r.GET("/ready", s.readyCheck)

	// JSON-RPC endpoint
	r.POST("/", s.handleJSONRPC)
	r.GET("/", s.handleJSONRPC) // Support GET for some clients

	// Create HTTP server
	srv := &http.Server{
		Addr:              addr,
		Handler:           r,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	return srv.ListenAndServe()
}

func (s *HTTPServer) handleJSONRPC(c *gin.Context) {
	var req Types.Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Types.RespErr(nil, -32700, "Parse error"))
		return
	}

	resp, _ := s.h.Handle(c.Request.Context(), req)
	c.JSON(http.StatusOK, resp)
}

func (s *HTTPServer) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"timestamp": time.Now().Unix(),
	})
}

func (s *HTTPServer) readyCheck(c *gin.Context) {
	// You could add more sophisticated readiness checks here
	c.JSON(http.StatusOK, gin.H{
		"status":    "ready",
		"timestamp": time.Now().Unix(),
	})
}
