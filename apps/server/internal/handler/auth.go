package handler

import (
	"net/http"

	errs "github.com/Niiaks/campusCart/internal/err"
	"github.com/Niiaks/campusCart/internal/middleware"
	"github.com/Niiaks/campusCart/internal/server"
	"github.com/Niiaks/campusCart/internal/service"
	"github.com/Niiaks/campusCart/pkg/types"
)

const (
	sessionCookieName = "session_id"
	sessionMaxAge     = 7 * 24 * 60 * 60 // 7 days in seconds
)

type AuthHandler struct {
	Handler
	authService *service.AuthService
	isProd      bool
}

func NewAuthHandler(server *server.Server, auth *service.AuthService) *AuthHandler {
	return &AuthHandler{
		Handler:     NewHandler(server),
		authService: auth,
		isProd:      server.Config.Primary.Env != "development",
	}
}

func (h *AuthHandler) setSessionCookie(w http.ResponseWriter, sessionID string) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		Secure:   h.isProd,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   sessionMaxAge,
	})
}

func (h *AuthHandler) Login() http.HandlerFunc {
	return Handle(h.Handler, func(w http.ResponseWriter, r *http.Request, req *types.LoginUser) (*types.LoginResponse, error) {
		resp, err := h.authService.Login(r.Context(), req)
		if err != nil {
			return nil, err
		}

		h.setSessionCookie(w, resp.SessionID)
		return resp, nil
	}, http.StatusOK, func() *types.LoginUser { return &types.LoginUser{} })
}

func (h *AuthHandler) Register() http.HandlerFunc {
	return Handle(h.Handler, func(w http.ResponseWriter, r *http.Request, req *types.RegisterUser) (*types.RegisterResponse, error) {
		resp, err := h.authService.Register(r.Context(), req)
		if err != nil {
			return nil, err
		}

		return resp, nil
	}, http.StatusCreated, func() *types.RegisterUser { return &types.RegisterUser{} })
}

func (h *AuthHandler) VerifyEmail() http.HandlerFunc {
	return Handle(h.Handler, func(w http.ResponseWriter, r *http.Request, req *types.VerifyEmailRequest) (*types.LoginResponse, error) {
		resp, err := h.authService.VerifyEmail(r.Context(), req)
		if err != nil {
			return nil, err
		}

		h.setSessionCookie(w, resp.SessionID)
		return resp, nil
	}, http.StatusOK, func() *types.VerifyEmailRequest { return &types.VerifyEmailRequest{} })
}

// Logout requires the Authenticate middleware (session ID comes from context).
func (h *AuthHandler) Logout() http.HandlerFunc {
	return HandleNoContent(h.Handler, func(w http.ResponseWriter, r *http.Request, req *types.EmptyRequest) error {
		sessionID := middleware.GetSessionID(r.Context())

		if err := h.authService.Logout(r.Context(), sessionID); err != nil {
			return err
		}

		// Clear the session cookie
		http.SetCookie(w, &http.Cookie{
			Name:     sessionCookieName,
			Value:    "",
			Path:     "/",
			HttpOnly: true,
			Secure:   h.isProd,
			SameSite: http.SameSiteLaxMode,
			MaxAge:   -1,
		})
		return nil
	}, http.StatusNoContent, func() *types.EmptyRequest { return &types.EmptyRequest{} })
}

// GetCurrentUser requires the Authenticate middleware (user comes from context).
func (h *AuthHandler) GetCurrentUser() http.HandlerFunc {
	return Handle(h.Handler, func(w http.ResponseWriter, r *http.Request, req *types.EmptyRequest) (*types.UserResponse, error) {
		user := middleware.GetAuthUser(r.Context())
		if user == nil {
			return nil, errs.NewUnauthorizedError("not authenticated", false)
		}

		userResponse, err := h.authService.GetCurrentUser(r.Context(), user.ID)
		if err != nil {
			return nil, err
		}

		return userResponse, nil
	}, http.StatusOK, func() *types.EmptyRequest { return &types.EmptyRequest{} })
}
