package config

import (
	"encoding/json"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshsys "github.com/cloudfoundry/bosh-utils/system"

	"github.com/bosh-oneandone-cpi/registry"
)

type errorMsg int

const (
	invalidPluginType = 0 + iota
	invalidOCIConfiguration
	invalidAgentConfiguration
	invalidRegistryClientConfiguration
)

var errMsgs = []string{
	"Unsupported cloud plugin type %s",
	"Invalid configuration",
	"Invalid agent options configuration",
	"Invalid registry client configuration",
}

func (e errorMsg) String() string {
	return errMsgs[e]
}

// Config represents the full CPI configuration
//
// It is passed to CPI by its invoker (BOSH cli or Director)
// via  by the --configFile=<configfile> startup option.
type Config struct {
	Cloud Cloud
}

// Cloud element in the Config
type Cloud struct {
	Plugin     string
	Properties CPIProperties
}

// CPIProperties element in Cloud.Config
type CPIProperties struct {
	OAO      OAOProperties
	Agent    registry.AgentOptions
	Registry registry.ClientOptions
}

// NewConfigFromPath unmarshals(builds) a Config
// from CPI configuration json persisted in a file on the
// file system.
func NewConfigFromPath(configFile string, fs boshsys.FileSystem) (Config, error) {
	var config Config

	if configFile == "" {
		return config, bosherr.Errorf("Must provide a config file")
	}

	bytes, err := fs.ReadFile(configFile)
	if err != nil {
		return config, bosherr.WrapErrorf(err, "Reading config file '%s'", configFile)
	}

	if err = json.Unmarshal(bytes, &config); err != nil {
		return config, bosherr.WrapError(err, "Unmarshalling config contents")
	}

	// Fix relative paths in OCIProperties.
	// Ideally this should be done by the template scripts (.erb files)
	// packaged in the release, since the template generates cpi.json and
	// other keys files (.pem and .pub)
	//
	// However, the template scripts don't have any knowledge of
	// file system location. One wuould expect that such information
	// would be available via a macro or other environment
	// variable, but  there isn't one.
	// So we resort fixing the paths here in the code;
	// read the config from the path given to CPI (arg -configFile)
	// use that information to replace the unmarshalled object
	// with a new one containing the fixed paths.
	old := config.Cloud.Properties.OAO
	config.Cloud.Properties.OAO = newSanitizedConfig(configFile, old)

	if err = config.Validate(); err != nil {
		return config, bosherr.WrapError(err, "Validating config")
	}
	return config, nil
}

func NewConfigFromString(configString string) (Config, error) {
	var config Config
	var err error
	if configString == "" {
		return config, bosherr.Errorf("Must provide a config")
	}

	if err = json.Unmarshal([]byte(configString), &config); err != nil {
		return config, bosherr.WrapError(err, "Unmarshalling config contents")
	}

	if err = config.Validate(); err != nil {
		return config, bosherr.WrapError(err, "Validating config")
	}

	return config, nil
}

// Validate performs a deep validation of a Config
func (c Config) Validate() error {
	if c.Cloud.Plugin != "oneandone" {
		return bosherr.Errorf(errorMsg(invalidPluginType).String(), c.Cloud.Plugin)
	}
	if err := c.Cloud.Properties.OAO.Validate(); err != nil {
		return bosherr.WrapError(err, errorMsg(invalidOCIConfiguration).String())
	}
	if err := c.Cloud.Properties.Agent.Validate(); err != nil {
		return bosherr.WrapError(err, errorMsg(invalidAgentConfiguration).String())
	}
	if err := c.Cloud.Properties.Registry.Validate(); err != nil {
		return bosherr.WrapError(err, errorMsg(invalidRegistryClientConfiguration).String())
	}
	return nil
}
