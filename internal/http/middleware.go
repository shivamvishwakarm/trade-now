package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shivamvishwakarm/trade-now/internal/auth"
)

// RequireAuth returns a Gin middleware that validates the access_token HttpOnly cookie.
// On success it stores *auth.Claims under the "claims" context key and calls Next().
// On failure it aborts with 401 and never calls the handler.
func RequireAuth(accessSecret string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		rawToken, err := ctx.Cookie("access_token")
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing access token"})
			return
		}

		claims, err := auth.ParseToken(rawToken, accessSecret)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired access token"})
			return
		}

		ctx.Set(auth.ClaimsContextKey, claims)
		ctx.Next()
	}
}
