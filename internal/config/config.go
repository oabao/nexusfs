package config

// TODO: Implement configuration loading using Viper.
// This will load settings from a file, environment variables, and command-line flags.
type Config struct {
	Gateway GatewayConfig
	Storage StorageConfig
}

type GatewayConfig struct {
	ListenAddress string
}

type StorageConfig struct {
	DataDir string
}
