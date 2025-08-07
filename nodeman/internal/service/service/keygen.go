package service

type Keygen interface {
	GenerateHS256Secret() (string, error)
	GenerateVlessUUID() (string, error)
}
