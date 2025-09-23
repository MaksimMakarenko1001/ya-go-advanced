package config

type diConfig struct {
	HTTP HTTPServerConfig `json:"http"`
}

func (cfg *diConfig) loadConfig() {
	cfg.HTTP = HTTPServerConfig{
		Address: `:8080`,
	}
}

type HTTPServerConfig struct {
	Address string `json:"address"`
}
