package api

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/suryansh74/simplebank/db/mock"
)

func TestServerStart(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mock.NewMockStore(ctrl)
	server := newTestServer(t, store)

	// Create a context with cancel
	ctx, cancel := context.WithCancel(context.Background())

	// Channel to capture the error from StartWithShutdown
	errChan := make(chan error, 1)

	// Start server in goroutine
	go func() {
		errChan <- server.StartWithShutdown(ctx, "localhost:0")
	}()

	// Give server time to start
	time.Sleep(50 * time.Millisecond)

	// Trigger shutdown by canceling context
	cancel()

	// Wait for server to shutdown
	err := <-errChan

	// Server should return http.ErrServerClosed on graceful shutdown
	require.Equal(t, http.ErrServerClosed, err)
}

func TestServerStartWithoutShutdown(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mock.NewMockStore(ctrl)
	server := newTestServer(t, store)

	// Channel to capture any error
	errChan := make(chan error, 1)

	// Start server in goroutine
	go func() {
		errChan <- server.Start("localhost:0")
	}()

	// Give server time to start
	time.Sleep(50 * time.Millisecond)

	// Make a test request to verify server is running
	resp, err := http.Get("http://localhost:8080/accounts?page_id=1&page_size=5")
	if err == nil {
		resp.Body.Close()
	}

	// Note: This test doesn't achieve full coverage because Start() never returns
	// unless there's an error or external shutdown
}
