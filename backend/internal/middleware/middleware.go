package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/vidhyalakshmi-ramanathan/SCRIBHYA-Editor/internal/dto"
	"github.com/vidhyalakshmi-ramanathan/SCRIBHYA-Editor/internal/response"
	"github.com/vidhyalakshmi-ramanathan/SCRIBHYA-Editor/jwt"
)

type ctxKey string

const claimsKey ctxKey = "claims"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if !strings.HasPrefix(header, "Bearer") {
			response.Error(w, response.ErrInvalidToken)
			return
		}

		tokenString := strings.TrimSpace(strings.TrimPrefix(header, "Bearer"))
		claims, err := jwt.VerifyToken(tokenString)
		if err != nil {
			http.Error(w, "invalid or expired token or unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), claimsKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetClaims(ctx context.Context) (*dto.ClaimsData, bool) {
	c, ok := ctx.Value(claimsKey).(*dto.ClaimsData)
	return c, ok
}
