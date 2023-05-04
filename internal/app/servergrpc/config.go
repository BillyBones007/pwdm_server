package servergrpc

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"

	"github.com/caarlos0/env"
)

// Default config file.
const DefaultConfigFile = "config.json"

// ServerConfig - configuration server structure.
type ServerConfig struct {
	PortgRPC   string `env:"GRPC_PORT" json:"grpc_port,omitempty"`
	DSN        string `env:"DSN" json:"dns,omitempty"`
	ConfigFile string `env:"CONFIG_FILE"`
}

// Set config from config file.
func (s *ServerConfig) setFileConfig(file string) error {
	if err := readConfigFile(file, s); err != nil {
		return err
	}
	s.ConfigFile = file
	return nil
}

// Set config from environment variables.
func (s *ServerConfig) setEnvConfig() error {
	if err := env.Parse(s); err != nil {
		return err
	}
	return nil
}

// ReplaceConfig - replace current config with a new config.
func (s *ServerConfig) replaceConfig(newCfg ServerConfig) {
	newValues := reflect.ValueOf(newCfg)
	oldValues := reflect.ValueOf(s).Elem()

	for j := 0; j < newValues.NumField(); j++ {
		if !newValues.Field(j).IsZero() {
			oldValues.Field(j).Set(newValues.Field(j))
		}
	}
}

// InitServerConfig - initializing the server configuration.
// The values have the following priority:
// 1 - values from environment variables are prioritized.
// 2 - values from confing file.
func InitServerConfig() *ServerConfig {
	mainConf := ServerConfig{}
	envConf := ServerConfig{}
	fileConf := ServerConfig{}

	if err := envConf.setEnvConfig(); err != nil {
		log.Fatal(err)
	}

	// Default config from config file;
	if envConf.ConfigFile != "" {
		if err := fileConf.setFileConfig(envConf.ConfigFile); err != nil {
			log.Fatal(err)
		}
	} else {
		if err := fileConf.setFileConfig(DefaultConfigFile); err != nil {
			log.Fatal(err)
		}
	}

	// Replace current config by priority.
	mainConf.replaceConfig(fileConf)
	mainConf.replaceConfig(envConf)

	paramConfigServerInfo(&mainConf)

	return &mainConf
}

// Displays information about the server configuration.
func paramConfigServerInfo(cfg *ServerConfig) {
	fmt.Println("Server configuration:")
	fmt.Printf("Port gRPC: %s\n", cfg.PortgRPC)
	fmt.Printf("DSN: %s\n", cfg.DSN)
	fmt.Printf("Config file: %s\n", cfg.ConfigFile)
}

// readConfigFile - read configuration file.
func readConfigFile(file string, config *ServerConfig) error {
	data, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, config)
	if err != nil {
		return err
	}
	return nil
}
