package internal

import (
	"encoding/json"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/justinas/alice"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

func LoggingHandler(log zerolog.Logger) alice.Constructor {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	return func(next http.Handler) http.Handler {
		return alice.New(
			// Install the logger handler with default output on the console
			hlog.NewHandler(log),

			hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
				hlog.FromRequest(r).Info().
					Str("method", r.Method).
					Stringer("url", r.URL).
					Int("status", status).
					Int("size", size).
					Dur("duration", duration).
					Msg("")
			}),

			// Install some provided extra handler to set some request's context fields.
			// Thanks to that handler, all our logs will come with some prepopulated fields.
			hlog.RemoteAddrHandler("ip"),
			hlog.UserAgentHandler("user_agent"),
			hlog.RefererHandler("referer"),
			hlog.RequestIDHandler("req_id", "Request-Id"),
		).Then(next)
	}
}

func PanicHandler() alice.Constructor {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				err := recover()
				if err != nil {
					log.Error().Err(err.(error)).Stack().Msgf("recovered from panic")
					log.Error().Msg(string(debug.Stack()))

					jsonBody, _ := json.Marshal(map[string]string{
						"error": "There was an internal server error",
					})

					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusInternalServerError)
					w.Write(jsonBody)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
