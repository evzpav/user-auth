package http_test

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/evzpav/documents/internal/domain"
	document "gitlab.com/evzpav/documents/internal/domain/document"
	internalHttp "gitlab.com/evzpav/documents/internal/infrastructure/server/http"
)

func TestDocument_Get(t *testing.T) {
	t.Run("should get all documents succesfully", func(t *testing.T) {
		objs := make([]*domain.Document, 0)
		objs = append(objs, &domain.Document{ID: "docId", Value: "docValue"})

		docServiceMock := &document.ServiceMock{
			GetAllFn: func(ctx context.Context, filter *domain.DocumentFilter, sort ...*domain.StorageSort) ([]*domain.Document, error) {
				return objs, nil
			},
		}

		server := httptest.NewServer(internalHttp.NewHandler(docServiceMock, testLog))
		defer server.Close()

		URL, _ := url.Parse(server.URL)

		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/documents", URL), nil)

		res, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)

		bodyBytesResp, err := ioutil.ReadAll(res.Body)
		assert.NoError(t, err)

		expectedBody := `
		{
			"documents": [
				{
					"id": "docId",
					"value": "docValue",
					"is_blacklisted": false,
					"doc_type": ""
				}
			]
		}`

		assert.JSONEq(t, expectedBody, string(bodyBytesResp))
		assert.Equal(t, http.StatusOK, res.StatusCode)

	})

	t.Run("should return internal error", func(t *testing.T) {
		docServiceMock := &document.ServiceMock{
			GetAllFn: func(ctx context.Context, filter *domain.DocumentFilter, sort ...*domain.StorageSort) ([]*domain.Document, error) {
				return nil, errors.New("not.found")
			},
		}
		server := httptest.NewServer(internalHttp.NewHandler(docServiceMock, testLog))

		defer server.Close()

		URL, _ := url.Parse(server.URL)

		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/documents", URL), nil)

		res, err := http.DefaultClient.Do(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	})
}

func Testdoc_Post(t *testing.T) {
	t.Run("should post a new document successfully", func(t *testing.T) {
		docServiceMock := &document.ServiceMock{
			CreateFn: func(ctx context.Context, document *domain.Document) (*domain.Document, error) {
				doc := &domain.Document{
					ID:    "docId",
					Value: "docValue",
				}

				assert.Equal(t, document, doc)
				return doc, nil
			},
		}
		server := httptest.NewServer(internalHttp.NewHandler(docServiceMock, testLog))

		defer server.Close()

		URL, _ := url.Parse(server.URL)

		bodyReader := strings.NewReader(`{ "id": "docId", "value":"docValue" }`)

		req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/document", URL), bodyReader)

		res, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusCreated, res.StatusCode)

		responseBody, err := ioutil.ReadAll(res.Body)
		assert.NoError(t, err)

		expectedBody := `
		{
			"id": "docId",
			"value": "docValue",
			"is_blacklisted": false,
			"doc_type": ""
		}`

		assert.JSONEq(t, expectedBody, string(responseBody))
	})

	t.Run("should return error with invalid body", func(t *testing.T) {
		server := httptest.NewServer(internalHttp.NewHandler(nil, testLog))

		defer server.Close()

		URL, _ := url.Parse(server.URL)

		bodyReader := strings.NewReader(`"id": "docId"`)

		req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/document", URL), bodyReader)
		res, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	})

	t.Run("should return internal error", func(t *testing.T) {
		docServiceMock := &document.ServiceMock{
			CreateFn: func(ctx context.Context, doc *domain.Document) (*domain.Document, error) {
				return nil, errors.New("error")
			},
		}
		server := httptest.NewServer(internalHttp.NewHandler(docServiceMock, testLog))
		defer server.Close()

		URL, _ := url.Parse(server.URL)

		bodyReader := strings.NewReader(`{ "id": "docId" }`)

		req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/document", URL), bodyReader)
		res, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	})
}

func TestDocument_Remove(t *testing.T) {
	t.Run("should remove a document successfully", func(t *testing.T) {
		docServiceMock := &document.ServiceMock{
			DeleteFn: func(ctx context.Context, ID string) error {
				return nil
			},
		}
		server := httptest.NewServer(internalHttp.NewHandler(docServiceMock, testLog))
		defer server.Close()

		URL, _ := url.Parse(server.URL)

		req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/document/1", URL), nil)
		res, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusNoContent, res.StatusCode)
	})

	t.Run("should return internal error", func(t *testing.T) {
		docServiceMock := &document.ServiceMock{
			DeleteFn: func(ctx context.Context, ID string) error {
				return errors.New("error")
			},
		}
		server := httptest.NewServer(internalHttp.NewHandler(docServiceMock, testLog))
		defer server.Close()

		URL, _ := url.Parse(server.URL)

		req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/document/1", URL), nil)
		res, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	})

	t.Run("should return error with no parameter", func(t *testing.T) {
		docServiceMock := &document.ServiceMock{
			DeleteFn: func(ctx context.Context, ID string) error {
				return errors.New("id.required")
			},
		}
		server := httptest.NewServer(internalHttp.NewHandler(docServiceMock, testLog))
		defer server.Close()

		URL, _ := url.Parse(server.URL)

		req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/document/", URL), nil)
		res, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusNotFound, res.StatusCode)
	})
}
