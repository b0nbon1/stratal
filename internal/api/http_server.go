package api

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/b0nbon1/stratal/internal/queue"
	"github.com/b0nbon1/stratal/internal/security"
	db "github.com/b0nbon1/stratal/internal/storage/db/sqlc"
	"github.com/b0nbon1/stratal/pkg/router"
)

type HTTPServer struct {
	isInitialized bool
	mtx           sync.Mutex
	router        *router.Router
	addr          string
	store         *db.SQLStore
	Server        *http.Server
	ctx           context.Context
	cancel        context.CancelFunc
	queue         queue.TaskQueue
	secretManager *security.SecretManager
}

func NewHTTPServer(addr string, store *db.SQLStore, queue queue.TaskQueue, secretManager *security.SecretManager) *HTTPServer {
	ctx, cancel := context.WithCancel(context.Background())

	s := &HTTPServer{
		addr:          addr,
		store:         store,
		ctx:           ctx,
		cancel:        cancel,
		queue:         queue,
		secretManager: secretManager,
	}

	return s
}

func (httpServer *HTTPServer) Start() error {
	httpServer.mtx.Lock()
	defer httpServer.mtx.Unlock()

	if httpServer.isInitialized {
		return fmt.Errorf("server already initialized")
	}

	httpServer.router = httpServer.registerRoutes()

	httpServer.Server = &http.Server{
		Addr:    httpServer.addr,
		Handler: httpServer.router,
	}

	httpServer.isInitialized = true

	return nil
}

func (httpServer *HTTPServer) Stop() error {
	httpServer.mtx.Lock()
	defer httpServer.mtx.Unlock()

	if !httpServer.isInitialized {
		return fmt.Errorf("server not initialized")
	}

	if err := httpServer.Server.Close(); err != nil {
		return fmt.Errorf("failed to stop server: %w", err)
	}

	httpServer.isInitialized = false
	return nil
}
