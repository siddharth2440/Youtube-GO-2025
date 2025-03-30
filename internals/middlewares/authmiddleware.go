package middlewares

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/itsmonday/youtube/internals/utils"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		tokenString, err := ctx.Cookie("youtubecookie")
		fmt.Printf("\n token  %v\n", tokenString)
		fmt.Printf("\n length token  %v\n", len(tokenString))
		if err != nil {
			ctx.AbortWithStatusJSON(
				401,
				gin.H{
					"success": false,
					"error":   "Unauthorized",
				},
			)
			return
		}
		if tokenString == "" {
			fmt.Println("Token is not Present")
			ctx.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
			return
		}

		// decode the token
		jwttoken, err := utils.JWT_Verification(tokenString)
		if err != nil {
			fmt.Printf("Token vverification failed %v", err)
			ctx.JSON(401, gin.H{"error": "Invalid token"})
			return
		}
		fmt.Printf("\n Actual token  %v\n", jwttoken)
		claims, ok := jwttoken.Claims.(jwt.MapClaims)
		if !ok || !jwttoken.Valid {
			fmt.Println("Invalid Token")
			ctx.AbortWithStatusJSON(
				401,
				gin.H{"success": false, "error": "Invalid token claims"})
			return
		}
		fmt.Printf("\nClaims %v\n", claims)
		ctx.Set("authuserid", claims["id"])

		val, _ := ctx.Get("authuserid")
		fmt.Printf("\nContext Value := %v\n", val)
		ctx.Next()
	}
}
