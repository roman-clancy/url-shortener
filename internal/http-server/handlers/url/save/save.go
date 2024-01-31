package save

import (
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
	resp "url_shortener/internal/util/api/response"
	"url_shortener/internal/util/sl"
)

type Request struct {
	Url   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	resp.Response
	Alias string `json:"alias"`
}

type UrlSaver interface {
	SaveUrl(urlToSave string, alias string) (int64, error)
}

func New(log *slog.Logger, urlSaver UrlSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		var req Request
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("cannot parse request body", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to decode request"))
			return
		}
		if err := validator.New().Struct(req); err != nil {
			log.Error("request validation failed", sl.Err(err))
			var validationErrors validator.ValidationErrors
			errors.As(err, &validationErrors)
			render.JSON(w, r, resp.ValidationErrors(validationErrors))
			return
		}
		alias := req.Alias
		if alias == "" {
			alias = "йцу23" //random.NewRandomString(6)
		}
		_, err = urlSaver.SaveUrl(req.Url, req.Alias)
		if err != nil {
			log.Error("error while saving url", sl.Err(err))
			render.JSON(w, r, resp.Error(err.Error()))
			return
		}
		render.JSON(w, r, resp.OK())
		return
	}
}
