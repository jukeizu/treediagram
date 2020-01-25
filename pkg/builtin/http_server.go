package builtin

import (
	"context"
	"net/http"
	"time"

	"github.com/jukeizu/contract"
	"github.com/rs/zerolog"
)

type HttpServer struct {
	logger     zerolog.Logger
	httpServer *http.Server
}

func NewHttpServer(logger zerolog.Logger, addr string) HttpServer {
	logger = logger.With().Str("component", "intent.endpoint.builtin").Logger()

	httpServer := http.Server{
		Addr: addr,
	}

	return HttpServer{logger, &httpServer}
}

func (h HttpServer) RegisterHandlers(handlers ...Handler) {
	mux := http.NewServeMux()

	for _, handler := range handlers {
		for _, r := range handler.Registrations() {
			mux.Handle("/builtin/"+r.Name, h.makeLoggingHttpHandlerFunc(r.Name, r.Handler))
		}
	}

	h.httpServer.Handler = mux
}

func (h HttpServer) Start() error {
	if h.httpServer == nil {
		return nil
	}

	h.logger.Info().Str("transport", "http").Msg("starting")

	return h.httpServer.ListenAndServe()

}

func (h HttpServer) Stop() error {
	h.logger.Info().Str("transport", "http").Msg("stopping")

	return h.httpServer.Shutdown(context.Background())
}

func (h HttpServer) makeLoggingHttpHandlerFunc(name string, f func(contract.Request) (*contract.Response, error)) http.HandlerFunc {
	contractHandlerFunc := contract.MakeRequestHttpHandlerFunc(f)

	return func(w http.ResponseWriter, r *http.Request) {
		defer func(begin time.Time) {
			h.logger.Info().
				Str("intent", name).
				Str("took", time.Since(begin).String()).
				Msg("called")
		}(time.Now())

		contractHandlerFunc.ServeHTTP(w, r)
	}
}
