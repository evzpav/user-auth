package http

import (
	"net/http"

	"gitlab.com/evzpav/documents/internal/domain"
	"gitlab.com/evzpav/documents/pkg/log"

	"github.com/gin-gonic/gin"
)

type handler struct {
	documentService domain.DocumentService

	log log.Logger
}

// NewHandler ...
func NewHandler(ls domain.DocumentService, log log.Logger) http.Handler {
	handler := &handler{
		documentService: ls,
		log:             log,
	}

	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(handler.logger(), handler.traceHeader(), handler.recovery())

	router.GET("/documents", handler.getDocuments)
	router.GET("/document/:id", handler.getDocument)
	router.POST("/document", handler.postDocument)
	router.PUT("/document/:id", handler.putDocument)
	router.DELETE("/document/:id", handler.deleteDocument)
	// router.GET("/status", handler.serverStatus)

	return router
}
