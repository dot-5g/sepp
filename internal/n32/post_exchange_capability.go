package n32

import (
	"net/http"
	"slices"

	"github.com/labstack/echo/v4"
)

type SecurityCapability string

const TLS = SecurityCapability("TLS")
const ALS = SecurityCapability("ALS")

type SecNegotiateReqData struct {
	Sender                     FQDN
	SupportedSecCapabilityList []SecurityCapability
}

type SecNegotiateRspData struct {
	Sender                FQDN
	SelectedSecCapability SecurityCapability
}

func (n32c *N32C) HandlePostExchangeCapability(c echo.Context) error {
	reqData := new(SecNegotiateReqData)

	if err := c.Bind(reqData); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if reqData.Sender == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Sender is required")
	}

	if len(reqData.SupportedSecCapabilityList) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "SupportedSecCapabilityList is required")
	}

	containsTLS := slices.Contains(reqData.SupportedSecCapabilityList, TLS)
	if !containsTLS {
		return echo.NewHTTPError(http.StatusBadRequest, "Bad SecurityCapability - Only TLS is supported")
	}

	rspData := SecNegotiateRspData{
		Sender:                n32c.FQDN,
		SelectedSecCapability: TLS,
	}

	return c.JSON(http.StatusOK, rspData)
}
