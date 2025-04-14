package delivery

import (
	"errors"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/magneless/pvz/internal/domain/models"
	"github.com/magneless/pvz/internal/middleware/auth"
	"github.com/magneless/pvz/internal/usecase"
	"github.com/magneless/pvz/pkg/logger/sl"
	"github.com/magneless/pvz/pkg/response"
)

type PVZHandler struct {
	usecase *usecase.PVZUsecase
	log     *slog.Logger
}

type CreatePVZRequest struct {
	ID               uuid.UUID   `json:"id"`
	RegistrationDate *time.Time  `json:"registrationDate"`
	City             models.City `json:"city" validate:"required"`
}

type GetPVZFullInfoRequest struct {
	StartDate *time.Time `json:"startDate"`
	EndDate   *time.Time `json:"endDate"`
	Page      *int       `json:"page"`
	Limit     *int       `json:"limit"`
}

func NewPVZHandler(usecase *usecase.PVZUsecase, log *slog.Logger) PVZHandler {
	return PVZHandler{
		usecase: usecase,
		log:     log,
	}
}

func (h *PVZHandler) CreatePVZ() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		
		const op = "handler.CreatePVZ"
		log := h.log.With(slog.String("op", op))
		
		role, ok := r.Context().Value(auth.RoleKey).(string)
		if !ok {
			log.Error("role not found")
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, response.ErrorResponse{Error: "role not found"})
			return
		}
		if role != "moderator" {
			log.Error("wrong role")
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, response.ErrorResponse{Error: "wrong role"})
			return
		}

		var req CreatePVZRequest

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

		pvz := models.PVZ{
			ID:               req.ID,
			RegistrationDate: req.RegistrationDate,
			City:             req.City,
		}

		err = h.usecase.CreatePVZ(r.Context(), pvz)
		if err != nil {
			log.Error("failed to create PVZ", sl.Err(err))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.ErrorResponse{Error: "iternal error"})
			return
		}

		log.Info("created PVZ", slog.String("ID", req.ID.String()))
		render.Status(r, http.StatusCreated)
		render.JSON(w, r, response.MessageResponse{Message: pvz})
	}
}

func (h *PVZHandler) GetPVZFullInfo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.GetPVZFullInfo"
		log := h.log.With(slog.String("op", op))

		role, ok := r.Context().Value(auth.RoleKey).(string)
		if !ok {
			log.Error("role not found")
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, response.ErrorResponse{Error: "role not found"})
			return
		}
		if role != "moderator" && role != "employee" {
			log.Error("wrong role")
			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, response.ErrorResponse{Error: "wrong role"})
			return
		}


		var req GetPVZFullInfoRequest

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

		pvzFullInfo, err := h.usecase.GetPVZFullInfo(
			r.Context(),
			req.StartDate,
			req.EndDate,
			req.Page,
			req.Limit,
		)
		if err != nil {
			log.Error("failed to get pvz full info", sl.Err(err))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.ErrorResponse{Error: err.Error()})
			return
		}

		log.Info("got full pvz info")
		render.Status(r, http.StatusOK)
		render.JSON(w, r, response.MessageResponse{Message: pvzFullInfo})
	}
}
