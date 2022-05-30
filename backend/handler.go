package backend

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func influencer(c echo.Context) error {
	userName := c.Param("user_name")
	useCase := UseCase{}
	resp, err := useCase.Influencer(userName)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]Model{"data": resp})
}

func Route(e *echo.Echo) {
	e.GET("/api/v1/influencer/:user_name", influencer)
}
