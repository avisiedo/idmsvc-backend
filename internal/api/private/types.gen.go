// Package private provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.16.3 DO NOT EDIT.
package private

// Defines values for SuccessProbe.
const (
	Ok SuccessProbe = "Ok"
)

// Error Schema to define the error response
type Error struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
}

// SuccessProbe Todo schema
type SuccessProbe string

// HealthyFailure Schema to define the error response
type HealthyFailure = Error

// HealthySuccess Todo schema
type HealthySuccess = SuccessProbe
