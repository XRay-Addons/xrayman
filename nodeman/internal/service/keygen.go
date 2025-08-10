package service

type Keygen interface {
	GenerateHS256Secret() ([]byte, error)
	GenerateVlessUUID() (string, error)
}
