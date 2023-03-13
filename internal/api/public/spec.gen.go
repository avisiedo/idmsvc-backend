// Package public provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.12.4 DO NOT EDIT.
package public

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"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",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %s", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	var res = make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	var resolvePath = PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		var pathToFile = url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
