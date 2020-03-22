package config

import (
	"fmt"

	"github.com/caarlos0/env"
	"github.com/golang/glog"
)

type DaemonConfig struct {
	PeriodicUpdate int `env:"PERIODIC_UPDATE" envDefault:"5"` // Interval between every check for the added and deleted pods
	GuidPool       GuidPoolConfig
	SubnetManager  SubnetManagerPluginConfig
}

type GuidPoolConfig struct {
	RangeStart string `env:"RANGE_START" envDefault:"02:00:00:00:00:00:00:00"` // First guid in the pool
	RangeEnd   string `env:"RANGE_END"   envDefault:"02:FF:FF:FF:FF:FF:FF:FF"` // Last guid in the pool
}

type SubnetManagerPluginConfig struct {
	Plugin string `env:"PLUGIN"` // Subnet manager plugin name
	Ufm    UFMConfig
}

type UFMConfig struct {
	Username    string `env:"UFM_USERNAME"`    // Username of ufm
	Password    string `env:"UFM_PASSWORD"`    // Password of ufm
	Address     string `env:"UFM_ADDRESS"`     // IP address or hostname of ufm server
	Port        int    `env:"UFM_PORT"`        // REST API port of ufm
	HttpSchema  string `env:"UFM_HTTP_SCHEMA"` // http or https
	Certificate string `env:"UFM_CERTIFICATE"` // Certificate of ufm
}

var supportedPlugins = []string{"noop", "ufm"}

func (dc *DaemonConfig) ReadConfig() error {
	glog.Info("ReadConfig():")
	err := env.Parse(dc)

	return err
}

func (dc *DaemonConfig) ValidateConfig() error {
	glog.Info("ValidateConfig():")
	if dc.PeriodicUpdate <= 0 {
		return fmt.Errorf("ValidateConfig(): invalid \"PeriodicUpdate\" value %v", dc.PeriodicUpdate)
	}

	if dc.SubnetManager.Plugin == "" {
		return fmt.Errorf("ValidateConfig(): no plugin selected, supported plugins %v", supportedPlugins)
	}

	if !dc.isSupportedPlugin() {
		return fmt.Errorf("ValidateConfig(): not supported plugin %s, supprted plugins %v",
			dc.SubnetManager.Plugin, supportedPlugins)
	}

	return nil
}

// isSupportedPlugin check if the plugin is supported
func (dc *DaemonConfig) isSupportedPlugin() bool {
	for _, plugin := range supportedPlugins {
		if dc.SubnetManager.Plugin == plugin {
			return true
		}
	}
	return false
}
