package action

import (
	"errors"

	boshlog "github.com/cloudfoundry/bosh-utils/logger"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/bosh-oneandone-cpi/oneandone/client"
	clientfakes "github.com/bosh-oneandone-cpi/oneandone/client/fakes"
	"github.com/bosh-oneandone-cpi/oneandone/disks"
	diskfakes "github.com/bosh-oneandone-cpi/oneandone/disks/fakes"
	"github.com/bosh-oneandone-cpi/oneandone/resource"
	"github.com/bosh-oneandone-cpi/oneandone/vm"
	vmfakes "github.com/bosh-oneandone-cpi/oneandone/vm/fakes"
	"github.com/bosh-oneandone-cpi/registry"
	registryfakes "github.com/bosh-oneandone-cpi/registry/fakes"
	sdk "github.com/1and1/oneandone-cloudserver-sdk-go"
)

var _ = Describe("AttachDisk", func() {
	var (
		err        error
		attacherVM *resource.Instance

		foundInstance *resource.Instance
		foundVolume   *sdk.BlockStorage

		expectedAgentSettings registry.AgentSettings

		registryClient *registryfakes.FakeClient
		connector      *clientfakes.FakeConnector
		logger         boshlog.Logger

		vmFinder *vmfakes.FakeVMFinder

		diskFinder       *diskfakes.FakeDiskFinder
		attacherDetacher *diskfakes.FakeAttacherDetacher

		attachDisk AttachDisk
	)

	BeforeEach(func() {
		vmFinder = &vmfakes.FakeVMFinder{}
		installVMFinderFactory(func(c client.Connector, l boshlog.Logger) vm.Finder {
			return vmFinder
		})

		diskFinder = &diskfakes.FakeDiskFinder{}
		installDiskFinderFactory(func(c client.Connector, l boshlog.Logger) disks.Finder {
			return diskFinder
		})

		attacherDetacher = &diskfakes.FakeAttacherDetacher{}
		installInstanceAttacherDetacherFactory(func(in *resource.Instance, c client.Connector, l boshlog.Logger) (disks.AttacherDetacher, error) {
			attacherVM = in
			return attacherDetacher, nil
		})

		connector = &clientfakes.FakeConnector{}
		logger = boshlog.NewLogger(boshlog.LevelNone)
		registryClient = &registryfakes.FakeClient{}

		attachDisk = NewAttachDisk(connector, logger, registryClient)
	})

	AfterEach(func() {
		resetAllFactories()
	})

	Describe("Run", func() {

		BeforeEach(func() {
			foundInstance = resource.NewInstance("fake-vm-ocid")
			vmFinder.FindInstanceResult = foundInstance

			diskFinder.FindStorageResult = foundVolume

			registryClient.FetchSettings = registry.AgentSettings{}
			expectedAgentSettings = registry.AgentSettings{
				Disks: registry.DisksSettings{
					Persistent: map[string]registry.PersistentSettings{
						"fake-vol-ocid": {
							ID:   "fake-vol-ocid",
							Path: "/dev/fake-path",
						},
					},
				},
			}
		})

		It("finds the vm", func() {
			_, err = attachDisk.Run("fake-vm-ocid", "fake-vol-ocid")
			Expect(err).NotTo(HaveOccurred())
			Expect(vmFinder.FindInstanceCalled).To(BeTrue())
			Expect(vmFinder.FindInstanceID).To(Equal("fake-vm-ocid"))
		})

		It("creates an attacher for the found vm", func() {
			_, err = attachDisk.Run("fake-vm-ocid", "fake-vol-ocid")
			Expect(err).NotTo(HaveOccurred())
			Expect(attacherVM).To(Equal(vmFinder.FindInstanceResult))
		})

		It("finds the disk", func() {
			_, err = attachDisk.Run("fake-vm-ocid", "fake-vol-ocid")
			Expect(err).NotTo(HaveOccurred())
			Expect(diskFinder.FindStorageCalled).To(BeTrue())
			Expect(diskFinder.FindStorageID).To(Equal("fake-vol-ocid"))
		})

		It("attaches the disk", func() {
			_, err = attachDisk.Run("fake-vm-ocid", "fake-vol-ocid")
			Expect(err).NotTo(HaveOccurred())
			Expect(attacherDetacher.AttachVolumeCalled).To(BeTrue())
			Expect(attacherDetacher.AttachedVolume).To(Equal(foundVolume))
			Expect(attacherDetacher.AttachedInstance).To(Equal(foundInstance))
		})

		It("udates the registry", func() {
			_, err = attachDisk.Run("fake-vm-ocid", "fake-vol-ocid")
			Expect(err).NotTo(HaveOccurred())
			Expect(registryClient.UpdateCalled).To(BeTrue())
			Expect(registryClient.UpdateSettings).To(Equal(expectedAgentSettings))

		})
		It("returns an error if vmfinder fails", func() {
			vmFinder.FindInstanceError = errors.New("fake-instance-finder-error")

			_, err = attachDisk.Run("fake-vm-ocid", "fake-vol-ocid")
			Expect(err).To(HaveOccurred())
			Expect(vmFinder.FindInstanceCalled).To(BeTrue())
			Expect(err.Error()).To(ContainSubstring("fake-instance-finder-error"))
			Expect(attacherDetacher.AttachVolumeCalled).To(BeFalse())
			Expect(registryClient.FetchCalled).To(BeFalse())
			Expect(registryClient.UpdateCalled).To(BeFalse())
		})

		It("returns an error if diskfinder fails", func() {
			diskFinder.FindStorageError = errors.New("fake-disk-finder-error")

			_, err = attachDisk.Run("fake-vm-ocid", "fake-vol-ocid")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-disk-finder-error"))
			Expect(diskFinder.FindStorageCalled).To(BeTrue())
			Expect(attacherDetacher.AttachVolumeCalled).To(BeFalse())
			Expect(registryClient.FetchCalled).To(BeFalse())
			Expect(registryClient.UpdateCalled).To(BeFalse())
		})

		It("returns an error if attacher fails", func() {
			attacherDetacher.AttachmentError = errors.New("fake-attachment-error")

			_, err = attachDisk.Run("fake-vm-ocid", "fake-vol-ocid")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-attachment-error"))
			Expect(diskFinder.FindStorageCalled).To(BeTrue())
			Expect(vmFinder.FindInstanceCalled).To(BeTrue())
			Expect(attacherDetacher.AttachVolumeCalled).To(BeTrue())
			Expect(registryClient.FetchCalled).To(BeFalse())
			Expect(registryClient.UpdateCalled).To(BeFalse())
		})

		It("returns an error if registryClient fetch call returns an error", func() {
			registryClient.FetchErr = errors.New("fake-registry-client-error")

			_, err = attachDisk.Run("fake-vm-ocid", "fake-vol-ocid")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-registry-client-error"))
			Expect(diskFinder.FindStorageCalled).To(BeTrue())
			Expect(vmFinder.FindInstanceCalled).To(BeTrue())
			Expect(attacherDetacher.AttachVolumeCalled).To(BeTrue())
			Expect(registryClient.FetchCalled).To(BeTrue())
			Expect(registryClient.UpdateCalled).To(BeFalse())
		})

		It("returns an error if registryClient update call returns an error", func() {
			registryClient.UpdateErr = errors.New("fake-registry-client-error")

			_, err = attachDisk.Run("fake-vm-ocid", "fake-vol-ocid")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-registry-client-error"))
			Expect(diskFinder.FindStorageCalled).To(BeTrue())
			Expect(vmFinder.FindInstanceCalled).To(BeTrue())
			Expect(attacherDetacher.AttachVolumeCalled).To(BeTrue())
			Expect(registryClient.FetchCalled).To(BeTrue())
			Expect(registryClient.UpdateCalled).To(BeTrue())
		})
	})

})
