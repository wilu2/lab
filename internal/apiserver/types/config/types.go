package config

type Config struct {
	Name     string `json:"name"`
	LogoData string `json:"logo_data"`
	Version  string `json:"version"`
}

type UpdateConfig struct {
	Name     string `json:"name"`
	LogoData string `json:"logo_data"`
}
