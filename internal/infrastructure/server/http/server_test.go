package http_test

import (
	"testing"
	"time"

	"gitlab.com/evzpav/documents/pkg/log"

	"github.com/stretchr/testify/assert"
	"gitlab.com/evzpav/documents/internal/infrastructure/server/http"
)

func TestServer_ListenAndServe(t *testing.T) {
	log := log.NewZeroLog("", "", log.Error)
	server := http.New(nil, "localhost", "9995", log)
	server.ListenAndServe()

	stopChan := make(chan bool)

	go func() {
		time.Sleep(1 * time.Second)
		stopChan <- true
	}()

	var result bool
	result = <-stopChan
	server.Shutdown()

	assert.True(t, result)
}
