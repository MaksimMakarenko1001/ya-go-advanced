package config

type diConfig struct {
	HTTP HTTPServerConfig `json:"http"`
}

func (cfg *diConfig) loadConfig() {
	cfg.HTTP = HTTPServerConfig{
		Port: "8080",
	}
}

type HTTPServerConfig struct {
	Port string `json:"port"`
}
