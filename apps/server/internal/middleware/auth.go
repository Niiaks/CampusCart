package middleware

import (
	"context"
	"net/http"

	errs "github.com/Niiaks/campusCart/internal/err"
	"github.com/Niiaks/campusCart/internal/model"
	"github.com/Niiaks/campusCart/internal/repository"
)

const (
	sessionCookieName = "session_id"
	sessionMaxAge     = 7 * 24 * 60 * 60 // 7 days in seconds

	SessionIDKey = "session_id"
	AuthUserKey  = "auth_user"
)

type AuthMiddleware struct {
	sessionRepo repository.SessionRepo
	isProd      bool
}

func NewAuthMiddleware(sessionRepo repository.SessionRepo, isProd bool) *AuthMiddleware {
	return &AuthMiddleware{
		sessionRepo: sessionRepo,
		isProd:      isProd,
	}
}

// Authenticate validates the session, refreshes it (sliding expiration), and stores
// the session ID + user in context. Returns 401 if no valid session.
func (am *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(sessionCookieName)
		if err != nil || cookie.Value == "" {
			panic(errs.NewUnauthorizedError("not authenticated", false))
		}

		sessionID := cookie.Value

		user, err := am.sessionRepo.GetUserBySession(r.Context(), sessionID)
		if err != nil {
			panic(errs.NewUnauthorizedError("invalid or expired session", false))
		}

		// Sliding expiration: refresh session + cookie
		_ = am.sessionRepo.RefreshSession(r.Context(), sessionID)
		am.setSessionCookie(w, sessionID)

		// Store session ID and user in context
		ctx := r.Context()
		ctx = context.WithValue(ctx, SessionIDKey, sessionID)
		ctx = context.WithValue(ctx, AuthUserKey, user)
		ctx = context.WithValue(ctx, UserIDKey, user.ID)
		ctx = context.WithValue(ctx, UserRoleKey, user.Role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (am *AuthMiddleware) setSessionCookie(w http.ResponseWriter, sessionID string) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		Secure:   am.isProd,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   sessionMaxAge,
	})
}

func (am *AuthMiddleware) ClearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   am.isProd,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})
}

// GetSessionID retrieves the session ID from the request context.
func GetSessionID(ctx context.Context) string {
	if id, ok := ctx.Value(SessionIDKey).(string); ok {
		return id
	}
	return ""
}

// GetAuthUser retrieves the authenticated user from the request context.
func GetAuthUser(ctx context.Context) *model.User {
	if user, ok := ctx.Value(AuthUserKey).(*model.User); ok {
		return user
	}
	return nil
}
