package delivery

import (
	"database/sql"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/magneless/pvz/internal/middleware/auth"
	"github.com/magneless/pvz/internal/usecase"
	"github.com/magneless/pvz/pkg/logger/sl"
	"github.com/magneless/pvz/pkg/response"
)

type ReceptionHandler struct {
	usecase *usecase.ReceptionUsecase
	log     *slog.Logger
}

type CreateReceptionRequest struct {
	PvzID uuid.UUID `json:"pvzId" validate:"required"`
}

func NewReceptionHandler(usecase *usecase.ReceptionUsecase, log *slog.Logger) ReceptionHandler {
	return ReceptionHandler{
		usecase: usecase,
		log:     log,
	}
}

func (h *ReceptionHandler) CreateReception() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		const op = "handler.CreateReception"
		log := h.log.With(slog.String("op", op))
		role, ok := r.Context().Value(auth.RoleKey).(string)
		if !ok {
			log.Error("role not found")
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, response.ErrorResponse{Error: "role not found"})
			return
		}
		if role != "employee" {
			log.Error("wrong role")
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, response.ErrorResponse{Error: "wrong role"})
			return
		}

		var req CreateReceptionRequest

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

		reception, err := h.usecase.CreateReception(r.Context(), req.PvzID)
		if err == sql.ErrNoRows {
			log.Error("failed to create Reception, another reception in_progress", sl.Err(err))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.ErrorResponse{Error: "failed to create Reception, another reception in_progress"})
			return
		}
		if err != nil {
			log.Error("failed to create Reception", sl.Err(err))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.ErrorResponse{Error: "iternal error"})
			return
		}

		log.Info("created Reception", slog.String("ID", reception.ID.String()))
		render.Status(r, http.StatusCreated)
		render.JSON(w, r, response.MessageResponse{Message: reception})
	}
}

func (h *ReceptionHandler) CloseReception() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.CloseReception"
		log := h.log.With(slog.String("op", op))

		role, ok := r.Context().Value(auth.RoleKey).(string)
		if !ok {
			log.Error("role not found")
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, response.ErrorResponse{Error: "role not found"})
			return
		}
		if role != "employee" {
			log.Error("wrong role")
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, response.ErrorResponse{Error: "wrong role"})
			return
		}

		pvzIdStr := chi.URLParam(r, "pvzId")
		if pvzIdStr == "" {
			log.Error("pvzId not provided in URL")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.ErrorResponse{Error: "pvzId not provided"})
			return
		}

		pvzId, err := uuid.Parse(pvzIdStr)
		if err != nil {
			log.Error("invalid UUID format", sl.Err(err))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.ErrorResponse{Error: "invalid pvzId format"})
			return
		}

		err = h.usecase.CloseReception(r.Context(), pvzId)
		if errors.Is(err, sql.ErrNoRows) {
			log.Error("not open reception", sl.Err(err))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.ErrorResponse{Error: "not open reception"})
			return
		}
		if err != nil {
			log.Error("failed to close reception", sl.Err(err))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.ErrorResponse{Error: "iternal error"})
			return
		}

		log.Info("reception closed")
		render.Status(r, http.StatusOK)
		render.JSON(w, r, response.MessageResponse{Message: "reception closed"})

	}
}
