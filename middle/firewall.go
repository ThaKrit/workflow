package middle

import (
	auth "go-api/token"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func Guard(secret string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// By cookie
		authCookie, err := ctx.Cookie("token")
		if err != nil {
			log.Println("Token missing in cookie")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "unauthorized",
			})
			return
		}

		// Remove prefix "Bearer " from auth token
		tokenString := strings.TrimPrefix(authCookie, "Bearer ")

		// Verify token
		token, err := auth.VerifyToken(tokenString, secret)
		if err != nil {
			log.Printf("Token verification failed: %v\\n", err)
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		log.Printf("Token verified successfully. Token: %s\n", token.Claims)

		ctx.Next()
	}
}
