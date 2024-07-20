package authware

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"strings"
	"tek-bank/internal/db/repository"
	"tek-bank/pkg/cresponse"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

const (
	currentUserLabel = "user"
)

/*
REQUIRED(Any middleware must have this)

For every middleware we need a config.
In config we also need to define a function which allows us to skip the middleware if return true.
By convention it should be named as "Filter" but any other name will work too.
*/
type Config struct {
	DBConnection            *gorm.DB
	AuthorizationHeaderKey  string
	AuthorizationTypeBearer string
	authorization           func(c *fiber.Ctx) error // middleware specfic
}

/*
Middleware specific
Function for generating default config
*/

func setup(config Config) Config {

	// Set default logging function if not passed
	config.authorization = func(c *fiber.Ctx) error {

		reqToken := c.Get(config.AuthorizationHeaderKey)

		// if authorization header is not found then skip
		if len(strings.TrimSpace(reqToken)) == 0 {
			return cresponse.ErrorResponse(c, fiber.StatusUnauthorized, "Authorization header is missing")
		}

		saltToken, err := ExtractToken(reqToken, config.AuthorizationTypeBearer)
		if err != nil {
			return cresponse.ErrorResponse(c, fiber.StatusUnauthorized, "Malformed token")
		}

		isTokenValid, claims, err := IsTokenValid(saltToken, os.Getenv("JWT_SECRET_KEY"))
		if !isTokenValid {
			return cresponse.ErrorResponse(c, fiber.StatusUnauthorized, "Token is not valid")
		}

		if err != nil {
			return cresponse.ErrorResponse(c, fiber.StatusUnauthorized, err.Error())
		}

		var claimsStruct JWTClaimsPayload
		jsonItem, err := json.Marshal(claims)
		if err != nil {
			return cresponse.ErrorResponse(c, fiber.StatusUnauthorized, "Token is not valid")
		}

		err = json.Unmarshal(jsonItem, &claimsStruct)
		if err != nil {
			return cresponse.ErrorResponse(c, fiber.StatusUnauthorized, "Token is not valid")
		}

		if err != nil {
			return cresponse.ErrorResponse(c, fiber.StatusUnauthorized, "Token is not valid")
		}

		isAuthorized := checkPermission(c, config.DBConnection, claimsStruct)

		if isAuthorized {
			err := c.Next()
			if err != nil {
				return err
			}
		} else {
			return cresponse.ErrorResponse(c, fiber.StatusUnauthorized, "Unauthorized")
		}
		return nil
	}

	return config
}

/*
REQUIRED(Any middleware must have this)

Our main middleware function used to initialize our middleware.
By convention, we name it "New" but any other name will work too.
*/
func New(config Config) fiber.Handler {

	// For setting default config
	cfg := setup(config)

	return func(c *fiber.Ctx) error {
		err := cfg.authorization(c)
		if err != nil {
			return err
		}
		return nil
	}
}

type CurrentUser struct {
	Id          string `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
	Username    string `json:"username"`
}

func checkPermission(ctx *fiber.Ctx, db *gorm.DB, claim JWTClaimsPayload) bool {
	userRepository := repository.NewUserRepository(db)

	user, err := userRepository.FindByEmail(claim.Email)
	if err != nil {
		return false
	}

	if user == nil {
		return false
	} else {

		var currentUser CurrentUser
		jsonItem, err := json.Marshal(user)
		if err != nil {
			return false
		}

		err = json.Unmarshal(jsonItem, &currentUser)
		if err != nil {
			return false
		}

		ctx.Locals(currentUserLabel, currentUser)

		return true
	}
}

func GetCurrentUser(ctx context.Context) (CurrentUser, error) {
	var response CurrentUser
	currentUser := ctx.Value(currentUserLabel)

	if currentUser == nil {
		return response, errors.New("User not found")
	}

	response, ok := currentUser.(CurrentUser)
	if !ok {
		return response, errors.New("User not found")
	}

	return response, nil
}
