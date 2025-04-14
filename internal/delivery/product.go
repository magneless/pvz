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
	"github.com/magneless/pvz/internal/domain/models"
	"github.com/magneless/pvz/internal/middleware/auth"
	"github.com/magneless/pvz/internal/usecase"
	"github.com/magneless/pvz/pkg/logger/sl"
	"github.com/magneless/pvz/pkg/response"
)

type ProductHandler struct {
	usecase *usecase.ProductUsecase
	log     *slog.Logger
}

type CreateProductRequest struct {
	Type  models.Type `json:"type" validate:"required"`
	PVZID uuid.UUID   `json:"pvzId" validate:"required"`
}

func NewProductHandler(usecase *usecase.ProductUsecase, log *slog.Logger) ProductHandler {
	return ProductHandler{
		usecase: usecase,
		log:     log,
	}
}

func (h *ProductHandler) CreateProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.CreateProduct"
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

		var req CreateProductRequest

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

		product, err := h.usecase.CreateProduct(r.Context(), req.Type, req.PVZID)
		if errors.Is(err, usecase.ErrWrongType) {
			log.Error("failed to create product", sl.Err(err))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.ErrorResponse{Error: "wrong type"})
			return
		}
		if err != nil {
			log.Error("failed to create product", sl.Err(err))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.ErrorResponse{Error: "iternal error"})
			return
		}

		log.Info("created product", slog.String("ID", product.ID.String()))
		render.Status(r, http.StatusCreated)
		render.JSON(w, r, response.MessageResponse{Message: product})
	}
}

func (h *ProductHandler) DeleteProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.DeleteProduct"
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

		err = h.usecase.DeleteProduct(r.Context(), pvzId)
		if errors.Is(err, sql.ErrNoRows) {
			log.Error("no active reception or product available for deletion", sl.Err(err))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.ErrorResponse{Error: "no active reception or product available for deletion"})
			return
		}
		if err != nil {
			log.Error("failed to delete product", sl.Err(err))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.ErrorResponse{Error: "iternal error"})
			return
		}

		log.Info("product deleted")
		render.Status(r, http.StatusOK)
		render.JSON(w, r, response.MessageResponse{Message: "product deleted"})
	}
}
