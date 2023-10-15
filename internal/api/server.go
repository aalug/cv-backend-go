package api

import (
	"github.com/aalug/cv-backend-go/internal/config"
	db "github.com/aalug/cv-backend-go/internal/db/sqlc"
	"github.com/gin-gonic/gin"
)

// Server serves HTTP  requests for the service
type Server struct {
	config config.Config
	store  db.Store
	router *gin.Engine
}

// NewServer creates a new HTTP server and setups routing
func NewServer(cfg config.Config, store db.Store) *Server {
	server := &Server{
		config: cfg,
		store:  store,
	}

	server.setupRouter()

	return server
}

// setupRouter sets up the HTTP routing
func (server *Server) setupRouter() {
	router := gin.Default()

	// --- cv profiles ---
	router.GET("/cv-profiles/:id", server.getCvProfile)

	// --- skills ---
	router.GET("/skills/:id", server.listSkills)

	server.router = router
}

// Start runs the HTTP server on a given address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
