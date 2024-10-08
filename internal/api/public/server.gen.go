// Package public provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.16.3 DO NOT EDIT.
package public

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/oapi-codegen/runtime"
)

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// List domains in the organization
	// (GET /domains)
	ListDomains(ctx echo.Context, params ListDomainsParams) error
	// Register a domain.
	// (POST /domains)
	RegisterDomain(ctx echo.Context, params RegisterDomainParams) error
	// Domain registration token request
	// (POST /domains/token)
	CreateDomainToken(ctx echo.Context, params CreateDomainTokenParams) error
	// Delete domain.
	// (DELETE /domains/{uuid})
	DeleteDomain(ctx echo.Context, uuid DomainIdParam, params DeleteDomainParams) error
	// Read a domain.
	// (GET /domains/{uuid})
	ReadDomain(ctx echo.Context, uuid DomainIdParam, params ReadDomainParams) error
	// Update domain information by user.
	// (PATCH /domains/{uuid})
	UpdateDomainUser(ctx echo.Context, uuid DomainIdParam, params UpdateDomainUserParams) error
	// Update domain information by ipa-hcc agent.
	// (PUT /domains/{uuid})
	UpdateDomainAgent(ctx echo.Context, uuid DomainIdParam, params UpdateDomainAgentParams) error
	// Get host vm information.
	// (POST /host-conf/{inventory_id}/{fqdn})
	HostConf(ctx echo.Context, inventoryId HostId, fqdn Fqdn, params HostConfParams) error
	// Signing keys
	// (GET /signing_keys)
	GetSigningKeys(ctx echo.Context, params GetSigningKeysParams) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// ListDomains converts echo context to params.
func (w *ServerInterfaceWrapper) ListDomains(ctx echo.Context) error {
	var err error

	ctx.Set(X_rh_identityScopes, []string{"Type:User", "Type:ServiceAccount"})

	// Parameter object where we will unmarshal all parameters from the context
	var params ListDomainsParams
	// ------------- Optional query parameter "offset" -------------

	err = runtime.BindQueryParameter("form", true, false, "offset", ctx.QueryParams(), &params.Offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter offset: %s", err))
	}

	// ------------- Optional query parameter "limit" -------------

	err = runtime.BindQueryParameter("form", true, false, "limit", ctx.QueryParams(), &params.Limit)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter limit: %s", err))
	}

	headers := ctx.Request().Header
	// ------------- Optional header parameter "X-Rh-Insights-Request-Id" -------------
	if valueList, found := headers[http.CanonicalHeaderKey("X-Rh-Insights-Request-Id")]; found {
		var XRhInsightsRequestId XRhInsightsRequestIdHeader
		n := len(valueList)
		if n != 1 {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Expected one value for X-Rh-Insights-Request-Id, got %d", n))
		}

		err = runtime.BindStyledParameterWithLocation("simple", false, "X-Rh-Insights-Request-Id", runtime.ParamLocationHeader, valueList[0], &XRhInsightsRequestId)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter X-Rh-Insights-Request-Id: %s", err))
		}

		params.XRhInsightsRequestId = &XRhInsightsRequestId
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.ListDomains(ctx, params)
	return err
}

// RegisterDomain converts echo context to params.
func (w *ServerInterfaceWrapper) RegisterDomain(ctx echo.Context) error {
	var err error

	ctx.Set(X_rh_identityScopes, []string{"Type:System"})

	ctx.Set(X_rh_idm_registration_tokenScopes, []string{})

	// Parameter object where we will unmarshal all parameters from the context
	var params RegisterDomainParams

	headers := ctx.Request().Header
	// ------------- Required header parameter "X-Rh-Idm-Registration-Token" -------------
	if valueList, found := headers[http.CanonicalHeaderKey("X-Rh-Idm-Registration-Token")]; found {
		var XRhIdmRegistrationToken XRhIdmRegistrationTokenHeader
		n := len(valueList)
		if n != 1 {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Expected one value for X-Rh-Idm-Registration-Token, got %d", n))
		}

		err = runtime.BindStyledParameterWithLocation("simple", false, "X-Rh-Idm-Registration-Token", runtime.ParamLocationHeader, valueList[0], &XRhIdmRegistrationToken)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter X-Rh-Idm-Registration-Token: %s", err))
		}

		params.XRhIdmRegistrationToken = XRhIdmRegistrationToken
	} else {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Header parameter X-Rh-Idm-Registration-Token is required, but not found"))
	}
	// ------------- Required header parameter "X-Rh-Idm-Version" -------------
	if valueList, found := headers[http.CanonicalHeaderKey("X-Rh-Idm-Version")]; found {
		var XRhIdmVersion XRhIdmVersionHeader
		n := len(valueList)
		if n != 1 {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Expected one value for X-Rh-Idm-Version, got %d", n))
		}

		err = runtime.BindStyledParameterWithLocation("simple", false, "X-Rh-Idm-Version", runtime.ParamLocationHeader, valueList[0], &XRhIdmVersion)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter X-Rh-Idm-Version: %s", err))
		}

		params.XRhIdmVersion = XRhIdmVersion
	} else {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Header parameter X-Rh-Idm-Version is required, but not found"))
	}
	// ------------- Optional header parameter "X-Rh-Insights-Request-Id" -------------
	if valueList, found := headers[http.CanonicalHeaderKey("X-Rh-Insights-Request-Id")]; found {
		var XRhInsightsRequestId XRhInsightsRequestIdHeader
		n := len(valueList)
		if n != 1 {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Expected one value for X-Rh-Insights-Request-Id, got %d", n))
		}

		err = runtime.BindStyledParameterWithLocation("simple", false, "X-Rh-Insights-Request-Id", runtime.ParamLocationHeader, valueList[0], &XRhInsightsRequestId)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter X-Rh-Insights-Request-Id: %s", err))
		}

		params.XRhInsightsRequestId = &XRhInsightsRequestId
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.RegisterDomain(ctx, params)
	return err
}

// CreateDomainToken converts echo context to params.
func (w *ServerInterfaceWrapper) CreateDomainToken(ctx echo.Context) error {
	var err error

	ctx.Set(X_rh_identityScopes, []string{"Type:User", "Type:ServiceAccount"})

	// Parameter object where we will unmarshal all parameters from the context
	var params CreateDomainTokenParams

	headers := ctx.Request().Header
	// ------------- Optional header parameter "X-Rh-Insights-Request-Id" -------------
	if valueList, found := headers[http.CanonicalHeaderKey("X-Rh-Insights-Request-Id")]; found {
		var XRhInsightsRequestId XRhInsightsRequestIdHeader
		n := len(valueList)
		if n != 1 {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Expected one value for X-Rh-Insights-Request-Id, got %d", n))
		}

		err = runtime.BindStyledParameterWithLocation("simple", false, "X-Rh-Insights-Request-Id", runtime.ParamLocationHeader, valueList[0], &XRhInsightsRequestId)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter X-Rh-Insights-Request-Id: %s", err))
		}

		params.XRhInsightsRequestId = &XRhInsightsRequestId
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.CreateDomainToken(ctx, params)
	return err
}

// DeleteDomain converts echo context to params.
func (w *ServerInterfaceWrapper) DeleteDomain(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "uuid" -------------
	var uuid DomainIdParam

	err = runtime.BindStyledParameterWithLocation("simple", false, "uuid", runtime.ParamLocationPath, ctx.Param("uuid"), &uuid)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter uuid: %s", err))
	}

	ctx.Set(X_rh_identityScopes, []string{"Type:User", "Type:ServiceAccount"})

	// Parameter object where we will unmarshal all parameters from the context
	var params DeleteDomainParams

	headers := ctx.Request().Header
	// ------------- Optional header parameter "X-Rh-Insights-Request-Id" -------------
	if valueList, found := headers[http.CanonicalHeaderKey("X-Rh-Insights-Request-Id")]; found {
		var XRhInsightsRequestId XRhInsightsRequestIdHeader
		n := len(valueList)
		if n != 1 {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Expected one value for X-Rh-Insights-Request-Id, got %d", n))
		}

		err = runtime.BindStyledParameterWithLocation("simple", false, "X-Rh-Insights-Request-Id", runtime.ParamLocationHeader, valueList[0], &XRhInsightsRequestId)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter X-Rh-Insights-Request-Id: %s", err))
		}

		params.XRhInsightsRequestId = &XRhInsightsRequestId
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.DeleteDomain(ctx, uuid, params)
	return err
}

// ReadDomain converts echo context to params.
func (w *ServerInterfaceWrapper) ReadDomain(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "uuid" -------------
	var uuid DomainIdParam

	err = runtime.BindStyledParameterWithLocation("simple", false, "uuid", runtime.ParamLocationPath, ctx.Param("uuid"), &uuid)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter uuid: %s", err))
	}

	ctx.Set(X_rh_identityScopes, []string{"Type:User", "Type:ServiceAccount"})

	// Parameter object where we will unmarshal all parameters from the context
	var params ReadDomainParams

	headers := ctx.Request().Header
	// ------------- Optional header parameter "X-Rh-Insights-Request-Id" -------------
	if valueList, found := headers[http.CanonicalHeaderKey("X-Rh-Insights-Request-Id")]; found {
		var XRhInsightsRequestId XRhInsightsRequestIdHeader
		n := len(valueList)
		if n != 1 {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Expected one value for X-Rh-Insights-Request-Id, got %d", n))
		}

		err = runtime.BindStyledParameterWithLocation("simple", false, "X-Rh-Insights-Request-Id", runtime.ParamLocationHeader, valueList[0], &XRhInsightsRequestId)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter X-Rh-Insights-Request-Id: %s", err))
		}

		params.XRhInsightsRequestId = &XRhInsightsRequestId
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.ReadDomain(ctx, uuid, params)
	return err
}

// UpdateDomainUser converts echo context to params.
func (w *ServerInterfaceWrapper) UpdateDomainUser(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "uuid" -------------
	var uuid DomainIdParam

	err = runtime.BindStyledParameterWithLocation("simple", false, "uuid", runtime.ParamLocationPath, ctx.Param("uuid"), &uuid)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter uuid: %s", err))
	}

	ctx.Set(X_rh_identityScopes, []string{"Type:User", "Type:ServiceAccount"})

	// Parameter object where we will unmarshal all parameters from the context
	var params UpdateDomainUserParams

	headers := ctx.Request().Header
	// ------------- Optional header parameter "X-Rh-Insights-Request-Id" -------------
	if valueList, found := headers[http.CanonicalHeaderKey("X-Rh-Insights-Request-Id")]; found {
		var XRhInsightsRequestId XRhInsightsRequestIdHeader
		n := len(valueList)
		if n != 1 {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Expected one value for X-Rh-Insights-Request-Id, got %d", n))
		}

		err = runtime.BindStyledParameterWithLocation("simple", false, "X-Rh-Insights-Request-Id", runtime.ParamLocationHeader, valueList[0], &XRhInsightsRequestId)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter X-Rh-Insights-Request-Id: %s", err))
		}

		params.XRhInsightsRequestId = &XRhInsightsRequestId
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.UpdateDomainUser(ctx, uuid, params)
	return err
}

// UpdateDomainAgent converts echo context to params.
func (w *ServerInterfaceWrapper) UpdateDomainAgent(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "uuid" -------------
	var uuid DomainIdParam

	err = runtime.BindStyledParameterWithLocation("simple", false, "uuid", runtime.ParamLocationPath, ctx.Param("uuid"), &uuid)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter uuid: %s", err))
	}

	ctx.Set(X_rh_identityScopes, []string{"Type:System:domain"})

	// Parameter object where we will unmarshal all parameters from the context
	var params UpdateDomainAgentParams

	headers := ctx.Request().Header
	// ------------- Required header parameter "X-Rh-Idm-Version" -------------
	if valueList, found := headers[http.CanonicalHeaderKey("X-Rh-Idm-Version")]; found {
		var XRhIdmVersion XRhIdmVersionHeader
		n := len(valueList)
		if n != 1 {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Expected one value for X-Rh-Idm-Version, got %d", n))
		}

		err = runtime.BindStyledParameterWithLocation("simple", false, "X-Rh-Idm-Version", runtime.ParamLocationHeader, valueList[0], &XRhIdmVersion)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter X-Rh-Idm-Version: %s", err))
		}

		params.XRhIdmVersion = XRhIdmVersion
	} else {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Header parameter X-Rh-Idm-Version is required, but not found"))
	}
	// ------------- Optional header parameter "X-Rh-Insights-Request-Id" -------------
	if valueList, found := headers[http.CanonicalHeaderKey("X-Rh-Insights-Request-Id")]; found {
		var XRhInsightsRequestId XRhInsightsRequestIdHeader
		n := len(valueList)
		if n != 1 {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Expected one value for X-Rh-Insights-Request-Id, got %d", n))
		}

		err = runtime.BindStyledParameterWithLocation("simple", false, "X-Rh-Insights-Request-Id", runtime.ParamLocationHeader, valueList[0], &XRhInsightsRequestId)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter X-Rh-Insights-Request-Id: %s", err))
		}

		params.XRhInsightsRequestId = &XRhInsightsRequestId
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.UpdateDomainAgent(ctx, uuid, params)
	return err
}

// HostConf converts echo context to params.
func (w *ServerInterfaceWrapper) HostConf(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "inventory_id" -------------
	var inventoryId HostId

	err = runtime.BindStyledParameterWithLocation("simple", false, "inventory_id", runtime.ParamLocationPath, ctx.Param("inventory_id"), &inventoryId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter inventory_id: %s", err))
	}

	// ------------- Path parameter "fqdn" -------------
	var fqdn Fqdn

	err = runtime.BindStyledParameterWithLocation("simple", false, "fqdn", runtime.ParamLocationPath, ctx.Param("fqdn"), &fqdn)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter fqdn: %s", err))
	}

	ctx.Set(X_rh_identityScopes, []string{"Type:System:domain"})

	// Parameter object where we will unmarshal all parameters from the context
	var params HostConfParams

	headers := ctx.Request().Header
	// ------------- Optional header parameter "X-Rh-Insights-Request-Id" -------------
	if valueList, found := headers[http.CanonicalHeaderKey("X-Rh-Insights-Request-Id")]; found {
		var XRhInsightsRequestId XRhInsightsRequestIdHeader
		n := len(valueList)
		if n != 1 {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Expected one value for X-Rh-Insights-Request-Id, got %d", n))
		}

		err = runtime.BindStyledParameterWithLocation("simple", false, "X-Rh-Insights-Request-Id", runtime.ParamLocationHeader, valueList[0], &XRhInsightsRequestId)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter X-Rh-Insights-Request-Id: %s", err))
		}

		params.XRhInsightsRequestId = &XRhInsightsRequestId
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.HostConf(ctx, inventoryId, fqdn, params)
	return err
}

// GetSigningKeys converts echo context to params.
func (w *ServerInterfaceWrapper) GetSigningKeys(ctx echo.Context) error {
	var err error

	ctx.Set(X_rh_identityScopes, []string{"Type:System", "Type:User", "Type:ServiceAccount"})

	// Parameter object where we will unmarshal all parameters from the context
	var params GetSigningKeysParams

	headers := ctx.Request().Header
	// ------------- Optional header parameter "X-Rh-Insights-Request-Id" -------------
	if valueList, found := headers[http.CanonicalHeaderKey("X-Rh-Insights-Request-Id")]; found {
		var XRhInsightsRequestId XRhInsightsRequestIdHeader
		n := len(valueList)
		if n != 1 {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Expected one value for X-Rh-Insights-Request-Id, got %d", n))
		}

		err = runtime.BindStyledParameterWithLocation("simple", false, "X-Rh-Insights-Request-Id", runtime.ParamLocationHeader, valueList[0], &XRhInsightsRequestId)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter X-Rh-Insights-Request-Id: %s", err))
		}

		params.XRhInsightsRequestId = &XRhInsightsRequestId
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetSigningKeys(ctx, params)
	return err
}

// This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {
	RegisterHandlersWithBaseURL(router, si, "")
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.GET(baseURL+"/domains", wrapper.ListDomains)
	router.POST(baseURL+"/domains", wrapper.RegisterDomain)
	router.POST(baseURL+"/domains/token", wrapper.CreateDomainToken)
	router.DELETE(baseURL+"/domains/:uuid", wrapper.DeleteDomain)
	router.GET(baseURL+"/domains/:uuid", wrapper.ReadDomain)
	router.PATCH(baseURL+"/domains/:uuid", wrapper.UpdateDomainUser)
	router.PUT(baseURL+"/domains/:uuid", wrapper.UpdateDomainAgent)
	router.POST(baseURL+"/host-conf/:inventory_id/:fqdn", wrapper.HostConf)
	router.GET(baseURL+"/signing_keys", wrapper.GetSigningKeys)

}
