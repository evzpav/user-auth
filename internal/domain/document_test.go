package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/evzpav/user-auth/internal/domain"
)

func TestValidadeDocument(t *testing.T) {
	t.Run("should validate document without errors", func(t *testing.T) {
		doc := &domain.Document{
			Value: "docvalue",
		}

		err := doc.Validate()
		assert.NoError(t, err)
	})

	t.Run("should validate document with id required error", func(t *testing.T) {
		doc := &domain.Document{}

		err := doc.Validate()
		assert.Equal(t, "<DOCUMENT_VALUE_REQUIRED> document value is required", err.Error())
	})

}
