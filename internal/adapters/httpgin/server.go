package httpgin

import (
	app "comixsearch/internal/app"
	"comixsearch/internal/app/interfaces"
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Server struct {
	port string
	app  *gin.Engine
}

func NewHTTPServer(port string, a app.SearchApp) *http.Server {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	api := router.Group("api/v1")
	s := &http.Server{Addr: port, Handler: router}
	AppRouter(api, a)

	return s
}

func (s *Server) Listen() error {
	return s.app.Run(s.port)
}

func (s *Server) Handler() http.Handler {
	return s.app
}

func Run(ctx context.Context, N interfaces.Normalizer, P interfaces.Fetcher, S interfaces.Storager, L interfaces.Logger, maxProc int, httpPort string) func() error {
	return func() error {
		httpServer := NewHTTPServer(httpPort, app.NewApp(N, P, S, L, maxProc))
		errCh := make(chan error)

		defer func() {
			shCtx, canel := context.WithTimeout(context.Background(), 30*time.Second)
			defer canel()

			if err := httpServer.Shutdown(shCtx); err != nil {
				log.Printf("can't close http server listing on %s: %s", httpServer.Addr, err.Error())
			}
			close(errCh)
		}()

		go func() {
			if err := httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
				errCh <- err
			}
		}()

		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-errCh:
			return fmt.Errorf("http server can't listen and serve requests: %w", err)
		}
	}
}
