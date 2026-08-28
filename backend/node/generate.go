package main

// place all generators here in correct order:

// 1. ogen - server by openapi
//go:generate go generate ./pkg/api/http/openapi/

// 2. converters
//go:generate go generate ./internal/http/handler/converter/

// 3. mocks
//go:generate go generate ./internal/http/handler/
//go:generate go generate ./internal/http/security/test/
