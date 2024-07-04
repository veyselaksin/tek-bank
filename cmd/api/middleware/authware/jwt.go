package authware

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"strings"
	"time"
)

type JWTClaimsPayload struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	jwt.RegisteredClaims
}

func GenerateJwtToken(payload JWTClaimsPayload, SecretKey string) (string, error) {
	var jwtKey = []byte(SecretKey)
	now := time.Now().UTC()

	claims := JWTClaimsPayload{
		ID:        payload.ID,
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Username:  payload.Username,
		Email:     payload.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{
				Time: now.Add(time.Hour * 24).UTC(),
			},
			IssuedAt:  &jwt.NumericDate{Time: now},
			NotBefore: &jwt.NumericDate{Time: now},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func IsTokenValid(token string, secretKey string) (bool, jwt.MapClaims, error) {

	claims := jwt.MapClaims{}
	decryptedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return false, jwt.MapClaims{}, err
	}

	if decryptedToken.Valid {
		return true, claims, nil
	} else {
		return false, jwt.MapClaims{}, err
	}
}

func ExtractToken(c *fiber.Ctx, authorizationHeaderKey, authorizationTypeBearer string) (string, error) {

	if len(authorizationHeaderKey) == 0 {
		return "", fmt.Errorf("Authorization header is missing")
	}
	splitToken := strings.Split(authorizationHeaderKey, " ")

	if len(splitToken) != 2 {
		return "", fmt.Errorf("Malformed token")
	}

	authorizationType := strings.ToLower(splitToken[0])
	if authorizationType != strings.ToLower(authorizationTypeBearer) {
		return "", fmt.Errorf("Authorization type is not Bearer")
	}

	authToken := strings.TrimSpace(splitToken[1])

	return authToken, nil
}
