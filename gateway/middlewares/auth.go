package middlewares

import (
	"fmt"

	"gateway/configs"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("jwt")
		if err != nil {
			c.Error(err)
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(configs.Env.JwtSecretAccess), nil
		})
		if err != nil || !token.Valid {
			c.Error(err)
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			fmt.Printf("Decoded JWT claims: %#v\n", claims)
			c.Set("user_id", uint(claims["user_id"].(float64)))
			c.Set("email", claims["email"])
			c.Set("name", claims["name"])
			if role, ok := claims["role"].(string); ok {
				c.Set("role", role) // ✅ lấy role từ JWT
			}
		}

		c.Next()
	}
}
