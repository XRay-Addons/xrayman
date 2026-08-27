package main

// place all generators here in correct order:

// 1. ogen - server by openapi + schemas by schemas
//go:generate go generate ./pkg/api/http/openapi/
//go:generate go generate ./pkg/api/http/schemas/

// 2. ogen - node client by openapi
//go:generate go generate ./internal/node/ogenclient/

// 3. sqlc - queries
//go:generate go generate ./internal/dbstorage/sqlc/

// 4. converters
//go:generate go generate ./internal/http/handler/converter/
//go:generate go generate ./internal/node/converter/
//go:generate go generate ./internal/pages/converter/

// 5. mocks
//go:generate go generate ./internal/http/handler/
//go:generate go generate ./internal/jobs/syncman/
