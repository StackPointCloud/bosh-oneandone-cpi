package client

import (
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	cclient "github.com/1and1/oneandone-cloudserver-sdk-go"

	"github.com/bosh-oneandone-cpi/config"
	"github.com/bosh-oneandone-cpi/registry"
)

const (
	logTag      = "OAOConnector"
	apiBasePath = "https://cloudpanel-api.1and1.com/v1"
)

type Connector interface {
	Connect() error
	Client() *cclient.API
	Token() string
	AgentOptions() registry.AgentOptions
	AgentRegistryEndpoint() string
	SSHTunnelConfig() config.SSHTunnel
}

type connectorImpl struct {
	config config.Cloud
	logger boshlog.Logger

	client *cclient.API
}

func NewConnector(c config.Cloud, logger boshlog.Logger) Connector {

	return &connectorImpl{config: c, logger: logger,
		client: nil}
}

func (c *connectorImpl) Connect() error {
	return c.createServiceClients(c.config.Properties.OAO.APIToken, apiBasePath)
}

func (c *connectorImpl) Client() *cclient.API {

	return c.client
}

func (c *connectorImpl) Token() string {

	return c.config.Properties.OAO.APIToken
}

func (c *connectorImpl) AgentOptions() registry.AgentOptions {
	return c.config.Properties.Agent
}

func (c *connectorImpl) AgentRegistryEndpoint() string {
	return c.config.Properties.Registry.EndpointWithCredentials()
}

func (c *connectorImpl) SSHTunnelConfig() config.SSHTunnel {
	return c.config.Properties.OAO.SSHTunnel
}

func (c *connectorImpl) createServiceClients(token string, basePath string) (error) {

	api := cclient.New(token, basePath)
	c.client = api

	_, err := c.client.ListDatacenters()
	if err != nil {
		c.logger.Error(logTag, "Error connecting to the API %s. Reason: %v", apiBasePath, err)
		return err
	}

	return nil
}

//func (c *connectorImpl) AuthorizedKeys() []string {
//	keys := []string{}
//	userKey, err := c.config.Properties.OAO.UserSSHPublicKeyContent()
//	if err != nil {
//		c.logger.Debug(logTag, "Ignored error while getting user key %v", err)
//	} else {
//		keys = append(keys, userKey)
//	}
//
//	cpiKey, err := c.config.Properties.OAO.CpiSSHPublicKeyContent()
//	if err != nil {
//		c.logger.Debug(logTag, "Ignored error while getting cpi key %v", err)
//	} else {
//		keys = append(keys, cpiKey)
//
//	}
//	return keys
//	return []string{}
//}

func (c *connectorImpl) createCoreServiceClient() {

}

func (c *connectorImpl) createIdentityServiceClient() {

}
