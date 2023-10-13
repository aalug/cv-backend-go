package api

import (
	"github.com/aalug/cv-backend-go/internal/config"
	db "github.com/aalug/cv-backend-go/internal/db/sqlc"
	"github.com/gin-gonic/gin"
	"os"
	"testing"
)

func newTestServer(store db.Store) *Server {
	server := NewServer(config.Config{}, store)

	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	os.Exit(m.Run())
}
