package integration

import (
	//"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"

	//"github.com/bosh-oneandone-cpi/oneandone/client"
	"fmt"
	"github.com/oneandone/oneandone-cloudserver-sdk-go"
)

var oneandoneClient *oneandone.ApiInstance

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Integration Suite")
}

var _ = SynchronizedBeforeSuite(func() []byte {
	// Clean any straggler VMs
	cleanVMs()
	var imageId = "8E9707A22051CDC1F0784F86411FBB45"

	request := fmt.Sprintf(`{
			  "method": "create_stemcell",
         "arguments": ["",
			{
           "name": "bosh-oneandone-kvm-ubuntu-trusty",
           "version": "3215",
           "infrastructure": "1&1",
           "image-id": "%s"
				}]
			}`, imageId)
	stemcell := assertSucceedsWithResult(request).(string)

	ips = make(chan string, len(ipAddrs))

	// Parse IP addresses to be used and put on a chan
	for _, addr := range ipAddrs {
		ips <- addr
	}

	return []byte(stemcell)
}, func(data []byte) {
	// Ensure stemcell was initialized
	existingStemcell = string(data)
	Expect(existingStemcell).ToNot(BeEmpty())

	//// Required env vars
	//Expect(googleProject).ToNot(Equal(""), "GOOGLE_PROJECT must be set")
	//Expect(externalStaticIP).ToNot(Equal(""), "EXTERNAL_STATIC_IP must be set")
	//Expect(serviceAccount).ToNot(Equal(""), "SERVICE_ACCOUNT must be set")

	//// Initialize a oneandone API client
	//var cc client.Connector
	//cc.Connect()
	//client := cc.Client()
	//Expect(client).ToNot(BeNil())
})

var _ = SynchronizedAfterSuite(func() {}, func() {
	//cleanVMs()
	//request := fmt.Sprintf(`{
	//		  "method": "delete_stemcell",
	//		  "arguments": ["%v"]
	//		}`, existingStemcell)
	//
	//response, err := execCPI(request)
	//Expect(err).ToNot(HaveOccurred())
	//Expect(response.Error).To(BeNil())
	//Expect(response.Result).To(BeNil())
})

func cleanVMs() {
	//// Initialize a compute API client
	//var cc client.Connector
	//cc.Connect()
	//client := cc.Client()
	//
	//// Clean up any VMs left behind from failed tests. Instances with the 'integration-delete' tag will be deleted.
	//var pageToken string
	//toDelete := make([]*compute.Instance, 0)
	//GinkgoWriter.Write([]byte("Looking for VMs with 'integration-delete' tag. Matches will be deleted\n"))
	//for {
	//	// Clean up VMs with 'integration-delete' tag
	//	listCall := computeService.Instances.AggregatedList(googleProject)
	//	listCall.PageToken(pageToken)
	//	aggregatedList, err := listCall.Do()
	//	Expect(err).To(BeNil())
	//	for _, list := range aggregatedList.Items {
	//		for _, instance := range list.Instances {
	//			for _, tag := range instance.Tags.Items {
	//				if tag == "integration-delete" {
	//					toDelete = append(toDelete, instance)
	//				}
	//			}
	//		}
	//	}
	//	if aggregatedList.NextPageToken == "" {
	//		break
	//	}
	//	pageToken = aggregatedList.NextPageToken
	//}
	//
	//for _, vm := range toDelete {
	//	GinkgoWriter.Write([]byte(fmt.Sprintf("Deleting VM %v\n", vm.Name)))
	//	_, err := computeService.Instances.Delete(googleProject, util.ResourceSplitter(vm.Zone), vm.Name).Do()
	//	Expect(err).ToNot(HaveOccurred())
	//}
}
