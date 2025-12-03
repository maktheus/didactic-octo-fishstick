package docs

import _ "embed"

// OpenAPISpec holds the embedded OpenAPI document.
//
//go:embed openapi.json
var OpenAPISpec []byte
