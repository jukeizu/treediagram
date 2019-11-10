package builtin

import (
	"context"
	"net/http"
	"time"

	"github.com/jukeizu/contract"
	"github.com/rs/zerolog"
)

type HttpServer struct {
	logger              zerolog.Logger
	httpServer          *http.Server
	helpHandler         HelpHandler
	selectServerHandler SelectServerHandler
	statsHandler        StatsHandler
	intentHandler       IntentHandler
}

func NewHttpServer(logger zerolog.Logger, addr string,
	helpHandler HelpHandler,
	selectServerHandler SelectServerHandler,
	statsHandler StatsHandler,
	intentHandler IntentHandler,
) HttpServer {
	logger = logger.With().Str("component", "intent.endpoint.builtin").Logger()

	httpServer := http.Server{
		Addr: addr,
	}

	return HttpServer{logger, &httpServer, helpHandler, selectServerHandler, statsHandler, intentHandler}
}

func (h HttpServer) Start() error {
	h.logger.Info().Str("transport", "http").Msg("starting")

	mux := http.NewServeMux()
	mux.HandleFunc("/help", h.makeLoggingHttpHandlerFunc("help", h.helpHandler.Help))
	mux.HandleFunc("/selectserver", h.makeLoggingHttpHandlerFunc("selectserver", h.selectServerHandler.SelectServer))
	mux.HandleFunc("/stats", h.makeLoggingHttpHandlerFunc("stats", h.statsHandler.Stats))
	mux.HandleFunc("/addintent", h.makeLoggingHttpHandlerFunc("addintent", h.intentHandler.AddIntent))
	mux.HandleFunc("/disableintent", h.makeLoggingHttpHandlerFunc("disableintent", h.intentHandler.DisableIntent))

	h.httpServer.Handler = mux

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
