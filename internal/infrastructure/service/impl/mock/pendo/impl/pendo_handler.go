package impl

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// CreateMetadataAccountCustomValue implement POST /metadata/account/custom/value
// endpoint for the pendo mock.
func (m *mockPendo) CreateMetadataAccountCustomValue(ctx echo.Context) error {
	return echo.NewHTTPError(http.StatusNotImplemented)
}
