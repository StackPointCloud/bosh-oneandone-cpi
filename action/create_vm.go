package action

import (
	"fmt"
	"time"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	boshuuid "github.com/cloudfoundry/bosh-utils/uuid"

	"github.com/bosh-oneandone-cpi/oneandone/client"
	"github.com/bosh-oneandone-cpi/oneandone/vm"
	"github.com/bosh-oneandone-cpi/registry"
)

// CreateVM action handles the create_vm request
type CreateVM struct {
	connector client.Connector
	logger    boshlog.Logger
	registry  registry.Client
	uuidGen   boshuuid.Generator
}

const logTag = "createVM"

// NewCreateVM creates a CreateVM instance
func NewCreateVM(c client.Connector, l boshlog.Logger, r registry.Client, u boshuuid.Generator) CreateVM {
	return CreateVM{connector: c, logger: l, registry: r, uuidGen: u}
}

func (cv CreateVM) Run(agentID string, stemcellCID StemcellCID, cloudProps VMCloudProperties, networks Networks, disks []DiskCID, env Environment) (VMCID, error) {

	// Create the VM
	name := cv.vmName(cloudProps.Name, cloudProps.Director)
	creator := newVMCreator(cv.connector, cv.logger)

	icfg := vm.InstanceConfiguration{
		Name:           name,
		ImageId:        string(stemcellCID),
		DatacenterId:   cloudProps.Datacenter,
		Ram:            cloudProps.Ram,
		DiskSize:       cloudProps.DiskSize,
		Cores:          cloudProps.Cores,
		Network:        networks.AsNetworkConfiguration(),
		SSHKey:         cloudProps.SSHKey,
		InstanceFlavor: cloudProps.InstanceFlavor,
		ServerIp:       cloudProps.PublicIP,
		SSHKeyPair:     cloudProps.SSHPairPath,
		EphemeralDisk:  cloudProps.EphemeralDisk,
		LBId:           cloudProps.LBId,
	}

	instance, err := creator.CreateInstance(icfg)
	if err != nil {
		return "", bosherr.WrapError(err, "Error creating new instance")
	}

	// Start a local forward ssh tunnel?
	if err == nil && networks.AllDynamic() {
		err = instance.EnsureReachable(cv.connector, cv.logger)
	}

	//privateIp, err := instance.PrivateIP(cv.connector, cv.logger)
	//if err != nil {
	//	return "", bosherr.WrapError(err, "Error launching new instance")
	//}

	publicIP, err := instance.PublicIP(cv.connector, cv.logger)
	if err != nil {
		return "", bosherr.WrapError(err, "Error launching new instance")
	}

	agentNetworks := networks.AsRegistryNetworks(publicIP)
	userdata := registry.NewUserDataObject(name, cv.connector.AgentRegistryEndpoint(), nil, agentNetworks)

	//check if an ssh key pair path was provided. if not set the defualt value to /vcap/.ssh
	if cloudProps.Director {
		cloudProps.SSHPairPath = "/root/.ssh"
	} else {
		cloudProps.SSHPairPath = "/home/vcap/.ssh"
	}

	if err := cv.updateRegistry(agentID, publicIP, name, cloudProps.SSHKey, cloudProps.SSHPairPath, agentNetworks, userdata, env); err != nil {
		return "Updating registry failed", err
	}
	return VMCID(instance.ID()), nil
}

func (cv CreateVM) vmName(prefix string, director bool) string {
	if director {
		prefix = "director " + prefix
	}
	suffix, err := cv.uuidGen.Generate()
	if err != nil {
		suffix = time.Now().Format(time.Stamp)
	}
	return fmt.Sprintf("%s-%s", prefix, suffix)
}

func (cv CreateVM) updateRegistry(agentID string, ipAddress string, vmName string, publicKey string, keyPairPath string,
	agentNetworks registry.NetworksSettings, userdata registry.UserData, env Environment) error {
	/*create vcap ssh directory and copy public key to it
	This is something that the agent does when using the registry,
	but since we are replacing it with FS we have to do this manually*/
	commands := []string{
		"usermod -G admin,bosh_sudoers,bosh_sshers vcap",
		"mkdir -m 0700 /home/vcap/.ssh",
		fmt.Sprintf("echo \"%s\" >> /home/vcap/.ssh/authorized_keys", publicKey),
		"chown -R vcap.vcap /home/vcap/.ssh",
		"chmod 0700 /home/vcap/.ssh",
	}
	var err error

	err = cv.registry.RunCommand(ipAddress, commands, keyPairPath)
	if err != nil {
		cv.logger.Error(logTag, "Error excuting vcap user commands %s", err)
	}

	//copy ssh key pair to each new node
	err = cv.registry.UploadRootKeyPair(ipAddress, keyPairPath)
	if err != nil {
		cv.logger.Error(logTag, "Error uploading ssh keys %s", err)
	}

	// Handle registry update failure. Delete VM or retry?
	defer func() {
		if err != nil {
			cv.logger.Error(logTag, "Registry update failed! FIXME: handle failure")
		}
	}()
	agentOptions := cv.connector.AgentOptions()
	agentSettings := registry.NewAgentSettings(agentID, vmName, agentNetworks,
		registry.EnvSettings(env), agentOptions, publicKey, userdata)

	//upload file with AgentSettings using FS and SCP
	err = cv.registry.UploadFile(ipAddress, agentSettings, keyPairPath)
	if err != nil {
		cv.logger.Error(logTag, "Error uploading userdata settings file %s", err)
	}
	return nil

}
