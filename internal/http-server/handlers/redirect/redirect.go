package redirect

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/kudras3r/url_shortener/internal/lib/api/response"
	"github.com/kudras3r/url_shortener/internal/lib/logger/sl"
)

type URLGetter interface {
	GetURL(alias string) (string, error)
}

type Response struct {
	response.Response
	Url string
}

func New(log *slog.Logger, urlGetter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "http-server.handlers.redirect.New"

		log := log.With(
			slog.String("op", op),
			slog.String("request-id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Info("alias is empty")

			render.JSON(w, r, response.Error("alias is empty"))

			return
		}

		url, err := urlGetter.GetURL(alias)
		if err != nil {
			log.Error("faidel to get url", sl.Err(err))

			render.JSON(w, r, response.Error("failed to get url"))

			return
		}

		log.Info("got url", slog.String("url", url))

		http.Redirect(w, r, url, http.StatusFound)
		//render.JSON(w, r, Response{
		// 	Response: response.Ok(),
		// 	Url:      url,
		// })
	}
}
