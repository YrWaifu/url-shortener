package redirect

import (
	"database/sql"
	"errors"
	"github.com/YrWaifu/url-shortener/internal/lib/api/response"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

type URLGetter interface {
	GetURL(alias string) (string, error)
}

func New(log *slog.Logger, urlGetter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.redirect.New"

		log = log.With("op", op)

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Info("alias is empty")

			render.JSON(w, r, response.Error("empty alias"))

			return
		}

		resURL, err := urlGetter.GetURL(alias)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				log.Error("url not found", "alias", alias, "error", err)

				render.JSON(w, r, response.Error("url not found"))

				return
			}

			log.Error("failed to get url", "alias", alias, "error", err)

			render.JSON(w, r, response.Error("failed to get url"))

			return
		}

		log.Info("got url", "alias", alias, "url", resURL)

		http.Redirect(w, r, resURL, http.StatusFound)
	}
}
