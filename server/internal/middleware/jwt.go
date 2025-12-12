package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/shamal-iroshan/notora/internal/config"
	jwtutil "github.com/shamal-iroshan/notora/internal/pkg/jwt"
)

// JWTMiddleware validates the access token stored in the HttpOnly cookie.
// If the token is valid, user_id is extracted from JWT claims and added to the context.
// Otherwise, the request is rejected with HTTP 401.
func JWTMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		// ------------------------------------------------------------
		// 1. Read access token from cookie
		// ------------------------------------------------------------
		accessToken, err := ctx.Cookie("access_token")
		if err != nil || accessToken == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "access token missing",
			})
			return
		}

		// ------------------------------------------------------------
		// 2. Validate JWT and extract claims
		// ------------------------------------------------------------
		claims, err := jwtutil.Parse([]byte(cfg.JWTSecret), accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid or expired token",
			})
			return
		}

		// ------------------------------------------------------------
		// 3. Extract user_id from claims safely
		// ------------------------------------------------------------
		rawUserID, exists := claims["user_id"]
		if !exists {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "token missing user_id",
			})
			return
		}

		// Claims in JWT are always float64 when decoded by default.
		userIDFloat, ok := rawUserID.(float64)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid user_id format",
			})
			return
		}

		userID := int64(userIDFloat)

		// Attach user ID to Gin context so handlers can access it.
		ctx.Set("user_id", userID)

		// Continue to the next handler in the chain.
		ctx.Next()
	}
}
