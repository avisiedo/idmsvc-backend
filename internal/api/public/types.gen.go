// Package public provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.12.4 DO NOT EDIT.
package public

import (
	"time"
)

const (
	X_rh_identityScopes = "x_rh_identity.Scopes"
)

// Defines values for CreateDomainDomainType.
const (
	CreateDomainDomainTypeIpa CreateDomainDomainType = "ipa"
)

// Defines values for DomainResponseDomainType.
const (
	DomainResponseDomainTypeIpa DomainResponseDomainType = "ipa"
)

// CheckHosts Define the input data for the /check_host action.
//
// This action is launched from the ipa server to check the host that is being
// auto-enrolled.
type CheckHosts struct {
	// DomainName The domain name where to enroll the host vm.
	DomainName string `json:"domain_name"`

	// DomainType Indicate the type of domain. Actually only ipa is supported.
	DomainType string `json:"domain_type"`

	// InventoryId The id of the host vm into the insight host inventory.
	InventoryId string `json:"inventory_id"`

	// SubscriptionManagerId The subscription manager id for auto-enroll the host vm.
	SubscriptionManagerId string `json:"subscription_manager_id"`
}

// CreateDomain A domain resource
type CreateDomain struct {
	// AutoEnrollmentEnabled Enable or disable host vm auto-enrollment for this domain
	AutoEnrollmentEnabled bool `json:"auto_enrollment_enabled"`

	// DomainDescription Humand readable description for this domain.
	DomainDescription string `json:"domain_description"`

	// DomainName Domain name
	DomainName string `json:"domain_name"`

	// DomainType Type of this domain. Currently only ipa is supported.
	DomainType CreateDomainDomainType `json:"domain_type"`

	// Ipa Options for ipa domains
	Ipa CreateDomainIpa `json:"ipa"`
}

// CreateDomainDomainType Type of this domain. Currently only ipa is supported.
type CreateDomainDomainType string

// CreateDomainIpa Options for ipa domains
type CreateDomainIpa struct {
	// CaCerts A base64 representation of all the list of chain of certificates, including the server ca.
	CaCerts []CreateDomainIpaCert `json:"ca_certs"`

	// RealmDomains TODO What is the meaning of this field.
	RealmDomains []string `json:"realm_domains"`

	// RealmName The kerberos realm name associated to the IPA domain.
	RealmName string `json:"realm_name"`

	// Servers List of auto-enrollment enabled servers for this domain.
	Servers *[]CreateDomainIpaServer `json:"servers,omitempty"`
}

// CreateDomainIpaCert Represent a certificate item in the cacerts list for the Ipa domain type.
type CreateDomainIpaCert struct {
	Issuer         *string    `json:"issuer,omitempty"`
	Nickname       *string    `json:"nickname,omitempty"`
	NotValidAfter  *time.Time `json:"not_valid_after,omitempty"`
	NotValidBefore *time.Time `json:"not_valid_before,omitempty"`
	Pem            *string    `json:"pem,omitempty"`
	SerialNumber   *string    `json:"serial_number,omitempty"`
	Subject        *string    `json:"subject,omitempty"`
}

// CreateDomainIpaServer Server schema for an entry into the Ipa domain type.
type CreateDomainIpaServer struct {
	CaServer            bool   `json:"ca_server"`
	Fqdn                string `json:"fqdn"`
	HccEnrollmentServer bool   `json:"hcc_enrollment_server"`
	HccUpdateServer     bool   `json:"hcc_update_server"`
	PkinitServer        bool   `json:"pkinit_server"`
	RhsmId              string `json:"rhsm_id"`
}

// DomainResponse A domain resource
type DomainResponse struct {
	// AutoEnrollmentEnabled Enable or disable host vm auto-enrollment for this domain
	AutoEnrollmentEnabled bool `json:"auto_enrollment_enabled"`

	// DomainDescription Human readable description abou the domain.
	DomainDescription string `json:"domain_description"`

	// DomainName Domain name
	DomainName string `json:"domain_name"`

	// DomainType Type of this domain. Currently only ipa is supported.
	DomainType DomainResponseDomainType `json:"domain_type"`

	// DomainUuid Internal id for this domain
	DomainUuid string `json:"domain_uuid"`

	// Ipa Options for ipa domains
	Ipa DomainResponseIpa `json:"ipa"`
}

// DomainResponseDomainType Type of this domain. Currently only ipa is supported.
type DomainResponseDomainType string

// DomainResponseIpa Options for ipa domains
type DomainResponseIpa struct {
	// CaCerts A base64 representation of all the list of chain of certificates, including the server ca.
	CaCerts []DomainResponseIpaCert `json:"ca_certs"`

	// RealmDomains List of realm associated to the IPA domain.
	RealmDomains []string `json:"realm_domains"`

	// RealmName The kerberos realm name associated to this IPA domain.
	RealmName string `json:"realm_name"`

	// Servers List of auto-enrollment enabled servers for this domain.
	Servers []DomainResponseIpaServer `json:"servers"`

	// Token One time token returned when the domain is created to let to register
	// an IPA domain by using hcc-ipa agent.
	Token *string `json:"token,omitempty"`

	// TokenExpiration When expire the one time token.
	TokenExpiration *time.Time `json:"token_expiration,omitempty"`
}

// DomainResponseIpaCert Represent a certificate item in the cacerts list for the Ipa domain type.
type DomainResponseIpaCert struct {
	Issuer         string    `json:"issuer"`
	Nickname       string    `json:"nickname"`
	NotValidAfter  time.Time `json:"not_valid_after"`
	NotValidBefore time.Time `json:"not_valid_before"`
	Pem            string    `json:"pem"`
	SerialNumber   string    `json:"serial_number"`
	Subject        string    `json:"subject"`
}

// DomainResponseIpaServer Server schema for an entry into the Ipa domain type.
type DomainResponseIpaServer struct {
	CaServer            bool   `json:"ca_server"`
	Fqdn                string `json:"fqdn"`
	HccEnrollmentServer bool   `json:"hcc_enrollment_server"`
	HccUpdateServer     bool   `json:"hcc_update_server"`
	PkinitServer        bool   `json:"pkinit_server"`
	RhsmId              string `json:"rhsm_id"`
}

// Error General error schema
type Error struct {
	// Detail A human-readable explanation specific to this occurrence of the problem. This field’s value can be localized.
	Detail string `json:"detail"`

	// Id A unique identifier for this particular occurrence of the problem.
	Id string `json:"id"`

	// Status The HTTP status code applicable to this problem, expressed as a string value. This SHOULD be provided.
	Status *string `json:"status,omitempty"`
}

// ErrorResponseSchema General error response returned by the hmsidm API
type ErrorResponseSchema struct {
	// Errors Error objects provide additional information about problems encountered while performing an operation.
	Errors *[]Error `json:"errors,omitempty"`
}

// HostConf Represent the request payload for the /hostconf/:fqdn endpoint.
type HostConf struct {
	// Fqdn The full qualified domain name of the host vm that is being enrolled.
	Fqdn *string `json:"fqdn,omitempty"`

	// SubscriptionManagerId The id for the subscription manager.
	SubscriptionManagerId *string `json:"subscription_manager_id,omitempty"`
}

// HostConfResponseSchema The response for the action to retrieve the host vm information when
// it is being enrolled. This action is taken from the host vm.
type HostConfResponseSchema struct {
	DomainName *string `json:"domain_name,omitempty"`
	DomainType *string `json:"domain_type,omitempty"`
	Ipa        *struct {
		CaCert     *string   `json:"ca_cert,omitempty"`
		RealmName  *string   `json:"realm_name,omitempty"`
		ServerList *[]string `json:"server_list,omitempty"`
	} `json:"ipa,omitempty"`
}

// ListDomainsData The data listed for the domains.
type ListDomainsData struct {
	AutoEnrollmentEnabled *bool   `json:"auto_enrollment_enabled,omitempty"`
	DomainName            *string `json:"domain_name,omitempty"`
	DomainType            *string `json:"domain_type,omitempty"`
	DomainUuid            *string `json:"domain_uuid,omitempty"`
}

// ListDomainsResponseSchema Represent a paginated result for a list of domains
type ListDomainsResponseSchema struct {
	// Data The content for this page.
	Data []ListDomainsData `json:"data"`

	// Links Represent the navigation links for the data paginated.
	Links PaginationLinks `json:"links"`

	// Meta Metadata for the paginated responses.
	Meta PaginationMeta `json:"meta"`
}

// PaginationLinks Represent the navigation links for the data paginated.
type PaginationLinks struct {
	// First Reference to the first page of the request.
	First *string `json:"first,omitempty"`

	// Last Reference to the last page of the request.
	Last *string `json:"last,omitempty"`

	// Next Reference to the next page of the request.
	Next *string `json:"next,omitempty"`

	// Previous Reference to the previous page of the request.
	Previous *string `json:"previous,omitempty"`
}

// PaginationMeta Metadata for the paginated responses.
type PaginationMeta struct {
	// Count Total number of registers without pagination.
	Count *int32 `json:"count,omitempty"`
}

// CreateDomainIpaResponse A domain resource
type CreateDomainIpaResponse = DomainResponse

// CreateDomainResponse A domain resource
type CreateDomainResponse = DomainResponse

// ErrorResponse General error response returned by the hmsidm API
type ErrorResponse = ErrorResponseSchema

// HostConfResponse The response for the action to retrieve the host vm information when
// it is being enrolled. This action is taken from the host vm.
type HostConfResponse = HostConfResponseSchema

// ListDomainsResponse Represent a paginated result for a list of domains
type ListDomainsResponse = ListDomainsResponseSchema

// ReadDomainResponse A domain resource
type ReadDomainResponse = DomainResponse

// CheckHostParams defines parameters for CheckHost.
type CheckHostParams struct {
	// XRhIdentity Identity header.
	XRhIdentity string `json:"X-Rh-Identity"`

	// XRhInsightsRequestId Request id for distributed tracing.
	XRhInsightsRequestId *string `json:"X-Rh-Insights-Request-Id,omitempty"`
}

// ListDomainsParams defines parameters for ListDomains.
type ListDomainsParams struct {
	// Offset pagination offset
	Offset *int `form:"offset,omitempty" json:"offset,omitempty"`

	// Limit Number of items per page
	Limit *int `form:"limit,omitempty" json:"limit,omitempty"`

	// XRhIdentity Identity header for the request
	XRhIdentity string `json:"X-Rh-Identity"`

	// XRhInsightsRequestId Request id for distributed tracing.
	XRhInsightsRequestId *string `json:"X-Rh-Insights-Request-Id,omitempty"`
}

// CreateDomainParams defines parameters for CreateDomain.
type CreateDomainParams struct {
	// XRhIdentity Identity header for the request
	XRhIdentity string `json:"X-Rh-Identity"`

	// XRhInsightsRequestId Request id for distributed tracing.
	XRhInsightsRequestId *string `json:"X-Rh-Insights-Request-Id,omitempty"`
}

// DeleteDomainParams defines parameters for DeleteDomain.
type DeleteDomainParams struct {
	// XRhIdentity Identity header for the request
	XRhIdentity string `json:"X-Rh-Identity"`

	// XRhInsightsRequestId Request id for distributed tracing.
	XRhInsightsRequestId *string `json:"X-Rh-Insights-Request-Id,omitempty"`
}

// ReadDomainParams defines parameters for ReadDomain.
type ReadDomainParams struct {
	// XRhIdentity Identity header for the request
	XRhIdentity string `json:"X-Rh-Identity"`

	// XRhInsightsRequestId Request id for distributed tracing.
	XRhInsightsRequestId *string `json:"X-Rh-Insights-Request-Id,omitempty"`
}

// RegisterIpaDomainParams defines parameters for RegisterIpaDomain.
type RegisterIpaDomainParams struct {
	// XRhIdentity Identity header
	XRhIdentity string `json:"X-Rh-Identity"`

	// XRhInsightsRequestId Request id
	XRhInsightsRequestId *string `json:"X-Rh-Insights-Request-Id,omitempty"`

	// XRhIDMRegistrationToken One time use token to register ipa information.
	XRhIDMRegistrationToken string `json:"X-Rh-IDM-Registration-Token"`
}

// HostConfParams defines parameters for HostConf.
type HostConfParams struct {
	// XRhIdentity The identity header of the request.
	XRhIdentity string `json:"X-Rh-Identity"`

	// XRhInsightsRequestId Unique request id for distributing tracing.
	XRhInsightsRequestId *string `json:"X-Rh-Insights-Request-Id,omitempty"`
}

// CheckHostJSONRequestBody defines body for CheckHost for application/json ContentType.
type CheckHostJSONRequestBody = CheckHosts

// CreateDomainJSONRequestBody defines body for CreateDomain for application/json ContentType.
type CreateDomainJSONRequestBody = CreateDomain

// RegisterIpaDomainJSONRequestBody defines body for RegisterIpaDomain for application/json ContentType.
type RegisterIpaDomainJSONRequestBody = CreateDomainIpa

// HostConfJSONRequestBody defines body for HostConf for application/json ContentType.
type HostConfJSONRequestBody = HostConf
