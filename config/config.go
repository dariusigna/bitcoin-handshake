package config

type Server struct {
	TargetNodeAddress string `envconfig:"TARGET_NODE_ADDRESS" default:"127.0.0.1:9333"`
	Network           string `envconfig:"NETWORK" default:"simnet"`
	LogLevel          string `envconfig:"LOG_LEVEL" default:"debug"`
}
