package delivery

import (
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/magneless/pvz/internal/domain/models"
	"github.com/magneless/pvz/internal/usecase"
	"github.com/magneless/pvz/pkg/logger/sl"
	"github.com/magneless/pvz/pkg/response"
)

type DummyLoginHandler struct {
	usecase *usecase.DummyLoginUsecase
	log     *slog.Logger
}

type CreateTokenRequest struct {
	Role models.Role `json:"role" validate:"required"`
}

func NewDumyyLoginHandler(usecase *usecase.DummyLoginUsecase, log *slog.Logger) DummyLoginHandler {
	return DummyLoginHandler{
		usecase: usecase,
		log:     log,
	}
}

func (h *DummyLoginHandler) CreateToken() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.CreateToken"
		log := h.log.With(slog.String("op", op))

		var req CreateTokenRequest

		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty", sl.Err(err))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.ErrorResponse{Error: "empty request"})
			return
		}
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.ErrorResponse{Error: "invalid JSON format"})
			return
		}

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)

			log.Error("invalid request", sl.Err(err))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.ValidationError(validateErr))
			return
		}

		token, err := h.usecase.CreateToken(models.DummyLogin{Role: req.Role})
		if errors.Is(err, usecase.ErrWrongRole) {
			log.Error("failed to create token", sl.Err(err))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.ErrorResponse{Error: "wrong request"})
			return
		}
		if err != nil {
			log.Error("failed to create token", sl.Err(err))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.ErrorResponse{Error: "iternal error"})
			return
		}

		log.Info("created Token")
		render.Status(r, http.StatusCreated)
		render.JSON(w, r, response.MessageResponse{Message: token})
	}
}
