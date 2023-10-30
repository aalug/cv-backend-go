package api

import (
	"github.com/aalug/cv-backend-go/docs"
	"github.com/aalug/cv-backend-go/internal/config"
	db "github.com/aalug/cv-backend-go/internal/db/sqlc"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

	routerV1 := router.Group("/api/v1")

	// CORS
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	routerV1.Use(cors.New(corsConfig))

	// Swagger docs
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	docs.SwaggerInfo.BasePath = "/api/v1"

	// --- cv profiles ---
	routerV1.GET("/cv-profiles/:id", server.getCvProfile)

	// --- skills ---
	routerV1.GET("/skills/:id", server.listSkills)

	// --- projects ---
	routerV1.GET("/projects/skill/:id/:skill", server.listProjectsBySkillName)
	routerV1.GET("/projects/:id", server.listProjects)

	server.router = router
}

// Start runs the HTTP server on a given address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func errorResponse(err error) ErrorResponse {
	return ErrorResponse{Error: err.Error()}
}
