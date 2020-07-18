package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/evzpav/documents/internal/domain"
)

func (h *handler) postDocument(c *gin.Context) {
	var doc domain.Document
	if err := c.ShouldBindJSON(&doc); err != nil {
		h.responseProblem(c, ErrInvalidBodyRequest)
		return
	}

	docCreated, err := h.documentService.Create(c, &doc)
	if err != nil {
		h.responseProblem(c, err)
		return
	}

	c.JSON(http.StatusCreated, docCreated)
}

func (h *handler) getDocuments(c *gin.Context) {
	filter := &domain.DocumentFilter{}
	documents, err := h.documentService.GetAll(c, filter)
	if err != nil {
		h.responseProblem(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"documents": documents,
	})
}

func (h *handler) getDocument(c *gin.Context) {
	id := c.Param("id")
	document, err := h.documentService.GetOne(c, id)
	if err != nil {
		h.responseProblem(c, err)
		return
	}

	c.JSON(http.StatusOK, document)
}

func (h *handler) putDocument(c *gin.Context) {
	id := c.Param("id")

	var document domain.Document
	if err := c.ShouldBindJSON(&document); err != nil {
		h.responseProblem(c, ErrInvalidBodyRequest)
		return
	}
	document.ID = id

	documentUpdated, err := h.documentService.Update(c, &document)
	if err != nil {
		h.responseProblem(c, err)
		return
	}

	c.JSON(http.StatusOK, &documentUpdated)
}

func (h *handler) deleteDocument(c *gin.Context) {
	id := c.Param("id")

	if err := h.documentService.Delete(c, id); err != nil {
		h.responseProblem(c, err)
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}
