package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/shamal-iroshan/notora/internal/config"
	"github.com/shamal-iroshan/notora/internal/pkg/jwt"
	"github.com/shamal-iroshan/notora/internal/repository"
)

// JWTMiddleware validates the JWT stored in cookies, loads the user from DB,
// and injects user_id, user_status, and is_admin into context.
func JWTMiddleware(cfg *config.Config, userRepo *repository.UserRepository) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		// 1. Read cookie
		accessToken, err := ctx.Cookie("access_token")
		if err != nil || accessToken == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "access token missing",
			})
			return
		}

		// 2. Validate token
		claims, err := jwt.Parse([]byte(cfg.JWTSecret), accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid or expired token",
			})
			return
		}

		// 3. Extract user_id
		rawUserID, exists := claims["user_id"]
		if !exists {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "user_id missing in token",
			})
			return
		}

		userIDFloat, ok := rawUserID.(float64)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid user_id format",
			})
			return
		}
		userID := int64(userIDFloat)

		// 4. Load user from DB
		user, err := userRepo.FindByID(userID)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "user not found",
			})
			return
		}

		// 5. Inject context values
		ctx.Set("user_id", userID)
		ctx.Set("user_status", user.Status)
		ctx.Set("is_admin", user.IsAdmin)

		ctx.Next()
	}
}
