package middleware

import (
	"context"
	"net/http"

	errs "github.com/Niiaks/campusCart/internal/err"
	"github.com/Niiaks/campusCart/internal/lib/tokenhash"
	"github.com/Niiaks/campusCart/internal/model"
	"github.com/Niiaks/campusCart/internal/repository"
)

const (
	sessionCookieName = "cc_refresh_token"
	sessionMaxAge     = 7 * 24 * 60 * 60 // 7 days in seconds

	SessionIDKey = "cc_refresh_token"
	AuthUserKey  = "auth_user"
	BrandIDKey   = "brand_id"
	// Role constants
	AdminRole = "admin"
)

type AuthMiddleware struct {
	sessionRepo repository.SessionRepo
	brandRepo   *repository.BrandRepository
	isProd      bool
}

func NewAuthMiddleware(sessionRepo repository.SessionRepo, brandRepo *repository.BrandRepository, isProd bool) *AuthMiddleware {
	return &AuthMiddleware{
		sessionRepo: sessionRepo,
		brandRepo:   brandRepo,
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

		sessionToken := cookie.Value
		hashedToken := tokenhash.Hash(sessionToken)

		user, err := am.sessionRepo.GetUserBySession(r.Context(), hashedToken)
		if err != nil {
			panic(errs.NewUnauthorizedError("invalid or expired session", false))
		}

		// Sliding expiration: refresh session + cookie
		_ = am.sessionRepo.RefreshSession(r.Context(), hashedToken)
		am.setSessionCookie(w, sessionToken)

		// Store session ID and user in context
		ctx := r.Context()
		ctx = context.WithValue(ctx, SessionIDKey, hashedToken)
		ctx = context.WithValue(ctx, AuthUserKey, user)
		ctx = context.WithValue(ctx, UserIDKey, user.ID)
		ctx = context.WithValue(ctx, UserRoleKey, user.Role)

		if am.brandRepo != nil {
			brandID, err := am.brandRepo.GetBrandIDBySeller(r.Context(), user.ID)
			if err != nil {
				panic(errs.NewUnauthorizedError("brand not found for user", false))
			}
			ctx = context.WithValue(ctx, BrandIDKey, brandID)
		}

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

// GetBrandID retrieves the authenticated user's brand ID from context.
func GetBrandID(ctx context.Context) string {
	if id, ok := ctx.Value(BrandIDKey).(string); ok {
		return id
	}
	return ""
}
