package pkg

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	_ "github.com/labstack/echo/v4"
	"net/http"
	"strings"
	"tesodev-korpes/ConsumerService/config"
	"tesodev-korpes/CustomerService/authentication"
	"tesodev-korpes/shared/helpers"
)

var secretKey string

func init() {

	appConf := config.GetAppConfig("dev")
	secretKey = appConf.SecretKey
}

func CorrelationIDMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		correlationID := c.Request().Header.Get("X-Correlation-Id")

		if correlationID == "" {
			correlationID = uuid.New().String()

			c.Request().Header.Set("X-Correlation-Id", correlationID)
		}

		c.Response().Header().Set("X-Correlation-Id", correlationID)

		return next(c)
	}
}

func Authenticate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		skipConditions := helpers.GetSkipConditions()

		reqPath := c.Path()
		reqMethod := c.Request().Method
		for _, condition := range skipConditions {
			if reqMethod == condition.Method && strings.HasPrefix(reqPath, condition.Path) {
				return next(c)
			}
		}
		tokenString := c.Request().Header.Get("Authentication")
		if tokenString == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "No Authentication header provided"})
		}

		if strings.TrimSpace(tokenString) == secretKey {
			return next(c)
		}
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		if err := authentication.VerifyJWT(tokenString); err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
		}

		client := NewRestClient()
		err := client.DoGetRequest("http://customer-service:1907/verify", nil, tokenString)

		if err != nil {
			fmt.Println("Error:", err)
		}

		return next(c)
	}
}
