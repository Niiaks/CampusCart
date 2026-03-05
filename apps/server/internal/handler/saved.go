package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	errs "github.com/Niiaks/campusCart/internal/err"
	"github.com/Niiaks/campusCart/internal/middleware"
	"github.com/Niiaks/campusCart/internal/model"
	"github.com/Niiaks/campusCart/internal/server"
	"github.com/Niiaks/campusCart/internal/service"
	"github.com/Niiaks/campusCart/pkg/types"
)

// SavedHandler wires HTTP endpoints to saved service.
type SavedHandler struct {
	Handler
	savedService *service.SavedService
}

func NewSavedHandler(server *server.Server, savedService *service.SavedService) *SavedHandler {
	return &SavedHandler{
		Handler:      NewHandler(server),
		savedService: savedService,
	}
}

func (h *SavedHandler) Save() http.HandlerFunc {
	return Handle(h.Handler, func(w http.ResponseWriter, r *http.Request, req *types.Save) (string, error) {
		userID := middleware.GetUserID(r)
		if userID == "" {
			return "failed to save", errs.NewUnauthorizedError("login to save listings", true).WithAction(&errs.Action{
				Type:    errs.ActionTypeRedirect,
				Message: "Log in to save listings",
				Value:   "/login",
			})
		}
		saved := &model.Saved{
			UserID:    userID,
			ListingID: req.ListingID,
		}

		if err := h.savedService.Save(r.Context(), saved); err != nil {
			return "", err
		}
		return "saved successfully", nil
	}, http.StatusCreated, func() *types.Save { return &types.Save{} })
}

func (h *SavedHandler) GetSaved() http.HandlerFunc {
	return Handle(h.Handler, func(w http.ResponseWriter, r *http.Request, _ *types.EmptyRequest) ([]model.Saved, error) {
		userID := middleware.GetUserID(r)
		if userID == "" {
			return nil, errs.NewUnauthorizedError("login to view saved listings", true).WithAction(&errs.Action{
				Type:    errs.ActionTypeRedirect,
				Message: "Log in to view your saved listings",
				Value:   "/login",
			})
		}

		saved, err := h.savedService.GetSaved(r.Context(), userID)
		if err != nil {
			return nil, err
		}
		return saved, nil
	}, http.StatusOK, func() *types.EmptyRequest { return &types.EmptyRequest{} })
}

func (h *SavedHandler) DeleteSaved() http.HandlerFunc {
	return HandleNoContent(h.Handler, func(w http.ResponseWriter, r *http.Request, _ *types.EmptyRequest) error {
		userID := middleware.GetUserID(r)
		if userID == "" {
			return errs.NewUnauthorizedError("login to remove saved listings", true).WithAction(&errs.Action{
				Type:    errs.ActionTypeRedirect,
				Message: "Log in to manage your saved listings",
				Value:   "/login",
			})
		}

		id := chi.URLParam(r, "id")
		if id == "" {
			missingCode := "MISSING_PARAMETER"
			return errs.NewBadRequestError("saved listing id is required", false, &missingCode,
				[]errs.FieldError{{Field: "id", Error: "is required"}}, nil)
		}

		return h.savedService.Remove(r.Context(), id)
	}, http.StatusNoContent, func() *types.EmptyRequest { return &types.EmptyRequest{} })
}
