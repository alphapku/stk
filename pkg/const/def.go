package def

type AdapterType string

const (
	MockAdapter     AdapterType = "mock"
	CoinbaseAdapter AdapterType = "coinbase"
)

type EnvMode string

const (
	DevMode  EnvMode = "dev"
	ProdMode EnvMode = "prod"
)
