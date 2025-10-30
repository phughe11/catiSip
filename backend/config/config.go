package config

import (
	"encoding/json"
	"os"
)

// Config represents the application configuration
type Config struct {
	SIP      SIPConfig      `json:"sip"`
	Server   ServerConfig   `json:"server"`
}

// SIPConfig represents SIP/FreeSWITCH configuration
type SIPConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Domain   string `json:"domain"`
}

// ServerConfig represents HTTP server configuration
type ServerConfig struct {
	Port int `json:"port"`
}

// Load reads configuration from file or environment
func Load() (*Config, error) {
	cfg := &Config{
		SIP: SIPConfig{
			Host:     getEnv("SIP_HOST", "localhost"),
			Port:     5060,
			Username: getEnv("SIP_USERNAME", "1000"),
			Password: getEnv("SIP_PASSWORD", "1234"),
			Domain:   getEnv("SIP_DOMAIN", "localhost"),
		},
		Server: ServerConfig{
			Port: 8080,
		},
	}

	// Try to load from config file if exists
	if configFile := os.Getenv("CONFIG_FILE"); configFile != "" {
		data, err := os.ReadFile(configFile)
		if err == nil {
			json.Unmarshal(data, cfg)
		}
	}

	return cfg, nil
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
