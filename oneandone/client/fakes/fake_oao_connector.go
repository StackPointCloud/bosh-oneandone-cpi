package fakes

import (
	"github.com/bosh-oneandone-cpi/config"
	"github.com/bosh-oneandone-cpi/registry"
	cclient "github.com/1and1/oneandone-cloudserver-sdk-go"
)

type FakeConnector struct {
	AgentOptionsResult registry.AgentOptions
}

func (fc *FakeConnector) Connect() error {
	return nil
}

func (fc *FakeConnector) Client() *cclient.API {
	return nil
}

func (fc *FakeConnector) AuthorizedKeys() []string {
	return []string{"ssh-rsa-fake"}
}
func (fc *FakeConnector) AgentOptions() registry.AgentOptions {
	return fc.AgentOptionsResult
}

func (fc *FakeConnector) AgentRegistryEndpoint() string {
	return "fake-agent-registry-endpoint"
}

func (fc *FakeConnector) SSHTunnelConfig() config.SSHTunnel {
	return config.SSHTunnel{}
}

func (fc *FakeConnector) SSHConfig() config.SSHConfig {
	return config.SSHConfig{}
}

func (c *FakeConnector) Token() string {

	return "fake-token"
}
