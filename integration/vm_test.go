package integration

//import (
//	"fmt"
//
//	"github.com/oneandone/oneandone-cloudserver-sdk-go"
//	. "github.com/onsi/ginkgo"
//	. "github.com/onsi/gomega"
//)
//
//var _ = Describe("VM", func() {
//	It("creates a VM with an invalid configuration and receives an error message with logs", func() {
//		request := fmt.Sprintf(`{
//			  "method": "create_vm",
//			  "arguments": [
//				"agent",
//				"%v",
//				{
//				  "name": "boshtest",
//				  "cores": 1,
//				  "diskSize": 20,
//				  "ram": 4
//				},
//				{
//				  "default": {
//					"type": "dynamic",
//					"cloud_properties": {
//					  "open-ports": [
//							{
//								"port-from":22,
//								"port-to":22,
//								"source":"0.0.0.0"
//
//							},
//							{
//								"port-from":80,
//								"port-to":80,
//								"source":"0.0.0.0"
//
//							},
//							{
//								"port-from":443,
//								"port-to":443,
//								"source":"0.0.0.0"
//
//							},
//							{
//								"port-from":8443,
//								"port-to":8443,
//								"source":"0.0.0.0"
//
//							},
//							{
//								"port-from":8447,
//								"port-to":8447,
//								"source":"0.0.0.0"
//							}
//						]
//					}
//				  }
//				}
//			  ]
//			}`, existingStemcell)
//		resp, err := execCPI(request)
//		Expect(err).ToNot(HaveOccurred())
//		Expect(resp.Error.Message).ToNot(BeEmpty())
//	})
//
//	It("executes the VM lifecycle", func() {
//		var vmCID string
//		By("creating a VM")
//		request := fmt.Sprintf(`{
//			  "method": "create_vm",
//			  "arguments": [
//				"agent",
//				"%v",
//				{
//				  "Name": "boshtest",
//				  "Cores": 1,
//				  "DiskSize": 20,
//				  "Ram": 4
//				},
//				{
//				  "default": {
//					"type": "dynamic",
//					"cloud_properties": {
//					  "open-ports": [
//							{
//								"port-from":22,
//								"port-to":22,
//								"source":"0.0.0.0"
//
//							},
//							{
//								"port-from":80,
//								"port-to":80,
//								"source":"0.0.0.0"
//
//							},
//							{
//								"port-from":443,
//								"port-to":443,
//								"source":"0.0.0.0"
//
//							},
//							{
//								"port-from":8443,
//								"port-to":8443,
//								"source":"0.0.0.0"
//
//							},
//							{
//								"port-from":8447,
//								"port-to":8447,
//								"source":"0.0.0.0"
//							}
//						]
//					}
//				  }
//				},
//				{
//				  "bosh": {
//					  "group_name": "1and1 test",
//					  "groups": ["micro-1and1", "dummy", "dummy", "micro-1and1-dummy", "dummy-dummy"]
//				  }
//				}
//			  ]
//			}`, existingStemcell)
//		vmCID = assertSucceedsWithResult(request).(string)
//
//		By("locating the VM")
//		request = fmt.Sprintf(`{
//			  "method": "has_vm",
//			  "arguments": ["%v"]
//			}`, vmCID)
//		exists := assertSucceedsWithResult(request).(bool)
//		Expect(exists).To(Equal(true))
//
//		expectedName := "boshtest"
//		assertValidVM(vmCID, func(instance *oneandone.Server) {
//			// Labels should be an exact match
//			Expect(instance.Name).To(BeEquivalentTo(expectedName))
//		})
//
//		updatedName := "updatedfrombosh"
//		request = fmt.Sprintf(`{
//			  "method": "set_vm_metadata",
//			  "arguments": [
//				"%v",
//				%v
//			  ]
//			}`, vmCID, updatedName)
//		assertSucceeds(request)
//		assertValidVM(vmCID, func(instance *oneandone.Server) {
//			// Labels should be an exact match
//			Expect(instance.Name).To(BeEquivalentTo(expectedName))
//		})
//
//		//By("rebooting the VM")
//		//request = fmt.Sprintf(`{
//		//	  "method": "reboot_vm",
//		//	  "arguments": ["%v"]
//		//	}`, vmCID)
//		//assertSucceeds(request)
//
//		By("deleting the VM")
//		request = fmt.Sprintf(`{
//			  "method": "delete_vm",
//			  "arguments": ["%v"]
//			}`, vmCID)
//		assertSucceeds(request)
//
//	})
//
//	//It("can create a VM with tags", func() {
//	//	By("creating a VM")
//	//	var vmCID string
//	//	request := fmt.Sprintf(`{
//	//		  "method": "create_vm",
//	//		  "arguments": [
//	//			"agent",
//	//			"%v",
//	//			{
//	//			   "machine_type": "n1-standard-1",
//	//				"zone": "%v",
//	//			   "tags": ["tag1", "tag2"]
//	//			},
//	//			{
//	//			  "default": {
//	//				"type": "dynamic",
//	//				"cloud_properties": {
//	//				  "tags": ["integration-delete"],
//	//				  "network_name": "%v"
//	//				}
//	//			  }
//	//			},
//	//			[],
//	//			{}
//	//		  ]
//	//		}`, existingStemcell, zone, networkName)
//	//	vmCID = assertSucceedsWithResult(request).(string)
//	//	//assertValidVM(vmCID, func(instance *compute.Instance) {
//	//	//	Expect(instance.Tags.Items).To(ConsistOf("integration-delete", "tag1", "tag2"))
//	//	//})
//	//
//	//	By("deleting the VM")
//	//	request = fmt.Sprintf(`{
//	//		  "method": "delete_vm",
//	//		  "arguments": ["%v"]
//	//		}`, vmCID)
//	//	assertSucceeds(request)
//	//})
//	//
//	//It("can create a VM with overlapping VM and network tags and VM properties that override network properties", func() {
//	//	By("creating a VM")
//	//	var vmCID string
//	//	request := fmt.Sprintf(`{
//	//		  "method": "create_vm",
//	//		  "arguments": [
//	//			"agent",
//	//			"%v",
//	//			{
//	//			   "machine_type": "n1-standard-1",
//	//			   "zone": "%v",
//	//			   "tags": ["tag1", "tag2", "integration-delete"],
//	//			   "ephemeral_external_ip": false,
//	//			   "ip_forwarding": false
//	//			},
//	//			{
//	//			  "default": {
//	//				"type": "dynamic",
//	//				"cloud_properties": {
//	//				  "tags": ["integration-delete"],
//	//				  "network_name": "%v",
//	//				  "ephemeral_external_ip": true,
//	//				  "ip_forwarding": true
//	//				}
//	//			  }
//	//			},
//	//			[],
//	//			{}
//	//		  ]
//	//		}`, existingStemcell, zone, networkName)
//	//	vmCID = assertSucceedsWithResult(request).(string)
//	//	//assertValidVM(vmCID, func(instance *compute.Instance) {
//	//	//	Expect(instance.Tags.Items).To(ConsistOf("integration-delete", "tag1", "tag2"))
//	//	//	Expect(instance.CanIpForward).To(Equal(false))
//	//	//	Expect(instance.NetworkInterfaces[0].AccessConfigs).To(BeEmpty())
//	//	//})
//	//
//	//	By("deleting the VM")
//	//	request = fmt.Sprintf(`{
//	//		  "method": "delete_vm",
//	//		  "arguments": ["%v"]
//	//		}`, vmCID)
//	//	assertSucceeds(request)
//	//})
//	//
//	//It("can create a VM with a public IP in a network with public IPs disabled ", func() {
//	//	By("creating a VM")
//	//	var vmCID string
//	//	request := fmt.Sprintf(`{
//	//		  "method": "create_vm",
//	//		  "arguments": [
//	//			"agent",
//	//			"%v",
//	//			{
//	//			   "machine_type": "n1-standard-1",
//	//			   "zone": "%v",
//	//			   "tags": ["tag1", "tag2", "integration-delete"],
//	//			   "ephemeral_external_ip": true,
//	//			   "ip_forwarding": false
//	//			},
//	//			{
//	//			  "default": {
//	//				"type": "dynamic",
//	//				"cloud_properties": {
//	//				  "tags": ["integration-delete"],
//	//				  "network_name": "%v",
//	//				  "ephemeral_external_ip": false,
//	//				  "ip_forwarding": true
//	//				}
//	//			  }
//	//			},
//	//			[],
//	//			{}
//	//		  ]
//	//		}`, existingStemcell, zone, networkName)
//	//	vmCID = assertSucceedsWithResult(request).(string)
//	//	//assertValidVM(vmCID, func(instance *compute.Instance) {
//	//	//	Expect(instance.Tags.Items).To(ConsistOf("integration-delete", "tag1", "tag2"))
//	//	//	Expect(instance.CanIpForward).To(Equal(false))
//	//	//	Expect(instance.NetworkInterfaces[0].AccessConfigs[0].Name).ToNot(BeEmpty())
//	//	//})
//	//
//	//	By("deleting the VM")
//	//	request = fmt.Sprintf(`{
//	//		  "method": "delete_vm",
//	//		  "arguments": ["%v"]
//	//		}`, vmCID)
//	//	assertSucceeds(request)
//	//})
//	//It("executes the VM lifecycle with disk attachment hints", func() {
//	//	By("creating two disks")
//	//	var request, diskCID, diskCID2, vmCID string
//	//	request = fmt.Sprintf(`{
//	//		  "method": "create_disk",
//	//		  "arguments": [32768, {"zone": "%v"}, ""]
//	//		}`, zone)
//	//	diskCID = assertSucceedsWithResult(request).(string)
//	//	diskCID2 = assertSucceedsWithResult(request).(string)
//	//
//	//	By("creating a VM with the disk attachment hints")
//	//	request = fmt.Sprintf(`{
//	//		  "method": "create_vm",
//	//		  "arguments": [
//	//			"agent",
//	//			"%v",
//	//			{
//	//			  "machine_type": "n1-standard-1",
//	//			  "zone": "%v"
//	//			},
//	//			{
//	//			  "default": {
//	//				"type": "dynamic",
//	//				"cloud_properties": {
//	//				  "tags": ["integration-delete"],
//	//				  "network_name": "%v"
//	//				}
//	//			  }
//	//			},
//	//			["%v", "%v"],
//	//			{}
//	//		  ]
//	//		}`, existingStemcell, zone, networkName, diskCID, diskCID2)
//	//	vmCID = assertSucceedsWithResult(request).(string)
//	//
//	//	By("deleting the disks")
//	//	request = fmt.Sprintf(`{
//	//		  "method": "delete_disk",
//	//		  "arguments": ["%v"]
//	//		}`, diskCID)
//	//	assertSucceeds(request)
//	//
//	//	request = fmt.Sprintf(`{
//	//		  "method": "delete_disk",
//	//		  "arguments": ["%v"]
//	//		}`, diskCID2)
//	//	assertSucceeds(request)
//	//
//	//	By("deleting the VM")
//	//	request = fmt.Sprintf(`{
//	//		  "method": "delete_vm",
//	//		  "arguments": ["%v"]
//	//		}`, vmCID)
//	//	assertSucceeds(request)
//	//})
//	//
//	//It("executes the VM lifecycle with custom machine type", func() {
//	//	By("creating a VM")
//	//	var vmCID string
//	//	request := fmt.Sprintf(`{
//	//		  "method": "create_vm",
//	//		  "arguments": [
//	//			"agent",
//	//			"%v",
//	//			{
//	//			  "cpu": 2,
//	//			  "ram": 5120,
//	//			  "zone": "%v"
//	//			},
//	//			{
//	//			  "default": {
//	//				"type": "dynamic",
//	//				"cloud_properties": {
//	//				  "tags": ["integration-delete"],
//	//				  "network_name": "%v"
//	//				}
//	//			  }
//	//			},
//	//			[],
//	//			{}
//	//		  ]
//	//		}`, existingStemcell, zone, networkName)
//	//	vmCID = assertSucceedsWithResult(request).(string)
//	//
//	//	By("deleting the VM")
//	//	request = fmt.Sprintf(`{
//	//		  "method": "delete_vm",
//	//		  "arguments": ["%v"]
//	//		}`, vmCID)
//	//	assertSucceeds(request)
//	//})
//	//
//	//It("executes the VM lifecycle in a specific zone", func() {
//	//	By("creating a VM")
//	//	var vmCID string
//	//
//	//	request := fmt.Sprintf(`{
//	//		  "method": "create_vm",
//	//		  "arguments": [
//	//			"agent",
//	//			"%v",
//	//			{
//	//			  "machine_type": "n1-standard-1",
//	//			  "zone": "%v"
//	//			},
//	//			{
//	//			  "default": {
//	//				"type": "dynamic",
//	//				"cloud_properties": {
//	//				  "tags": ["integration-delete"],
//	//				  "network_name": "%v"
//	//				}
//	//			  }
//	//			},
//	//			[],
//	//			{}
//	//		  ]
//	//		}`, existingStemcell, zone, networkName)
//	//	vmCID = assertSucceedsWithResult(request).(string)
//	//
//	//	By("deleting the VM")
//	//	request = fmt.Sprintf(`{
//	//		  "method": "delete_vm",
//	//		  "arguments": ["%v"]
//	//		}`, vmCID)
//	//})
//	//
//	//It("executes the VM lifecycle with automatic restart disabled", func() {
//	//	By("creating a VM")
//	//	var vmCID string
//	//	request := fmt.Sprintf(`{
//	//		  "method": "create_vm",
//	//		  "arguments": [
//	//			"agent",
//	//			"%v",
//	//			{
//	//			  "machine_type": "n1-standard-1",
//	//			  "zone": "%v",
//	//			  "automatic_restart": false
//	//			},
//	//			{
//	//			  "default": {
//	//				"type": "dynamic",
//	//				"cloud_properties": {
//	//				  "tags": ["integration-delete"],
//	//				  "network_name": "%v"
//	//				}
//	//			  }
//	//			},
//	//			[],
//	//			{}
//	//		  ]
//	//		}`, existingStemcell, zone, networkName)
//	//	vmCID = assertSucceedsWithResult(request).(string)
//	//
//	//	By("deleting the VM")
//	//	request = fmt.Sprintf(`{
//	//		  "method": "delete_vm",
//	//		  "arguments": ["%v"]
//	//		}`, vmCID)
//	//	assertSucceeds(request)
//	//})
//	//
//	//It("execute the VM lifecycle with OnHostMaintenance modified", func() {
//	//	By("creating a VM")
//	//	var vmCID string
//	//
//	//	request := fmt.Sprintf(`{
//	//		  "method": "create_vm",
//	//		  "arguments": [
//	//			"agent",
//	//			"%v",
//	//			{
//	//			  "machine_type": "n1-standard-1",
//	//			  "zone": "%v",
//	//			  "on_host_maintenance": "TERMINATE"
//	//			},
//	//			{
//	//			  "default": {
//	//				"type": "dynamic",
//	//				"cloud_properties": {
//	//				  "tags": ["integration-delete"],
//	//				  "network_name": "%v"
//	//				}
//	//			  }
//	//			},
//	//			[],
//	//			{}
//	//		  ]
//	//		}`, existingStemcell, zone, networkName)
//	//
//	//	By("deleting the VM")
//	//	vmCID = assertSucceedsWithResult(request).(string)
//	//	request = fmt.Sprintf(`{
//	//		  "method": "delete_vm",
//	//		  "arguments": ["%v"]
//	//		}`, vmCID)
//	//	assertSucceeds(request)
//	//})
//	//
//	//It("can execute the VM lifecycle with a preemptible VM", func() {
//	//	By("creating a VM")
//	//	var vmCID string
//	//
//	//	request := fmt.Sprintf(`{
//	//		  "method": "create_vm",
//	//		  "arguments": [
//	//			"agent",
//	//			"%v",
//	//			{
//	//			  "machine_type": "n1-standard-1",
//	//			  "zone": "%v",
//	//			  "preemtible": true
//	//			},
//	//			{
//	//			  "default": {
//	//				"type": "dynamic",
//	//				"cloud_properties": {
//	//				  "tags": ["integration-delete"],
//	//				  "network_name": "%v"
//	//				}
//	//			  }
//	//			},
//	//			[],
//	//			{}
//	//		  ]
//	//		}`, existingStemcell, zone, networkName)
//	//	vmCID = assertSucceedsWithResult(request).(string)
//	//
//	//	By("deleting the VM")
//	//	request = fmt.Sprintf(`{
//	//		  "method": "delete_vm",
//	//		  "arguments": ["%v"]
//	//		}`, vmCID)
//	//	assertSucceeds(request)
//	//})
//	//
//	//var vmCID string
//	//It("executes the VM lifecycle with default service scopes and no service account", func() {
//	//	By("creating a VM")
//	//	request := fmt.Sprintf(`{
//	//		  "method": "create_vm",
//	//		  "arguments": [
//	//			"agent",
//	//			"%v",
//	//			{
//	//			  "machine_type": "n1-standard-1",
//	//			  "zone": "%v",
//	//			  "service_scopes": ["cloud-platform", "devstorage.read_write"]
//	//			},
//	//			{
//	//			  "default": {
//	//				"type": "dynamic",
//	//				"cloud_properties": {
//	//				  "tags": ["integration-delete"],
//	//				  "network_name": "%v"
//	//				}
//	//			  }
//	//			},
//	//			[],
//	//			{}
//	//		  ]
//	//		}`, existingStemcell, zone, networkName)
//	//	vmCID = assertSucceedsWithResult(request).(string)
//	//	//assertValidVM(vmCID, func(instance *compute.Instance) {
//	//	//	// Labels should be an exact match
//	//	//	Expect(instance.ServiceAccounts[0].Scopes).To(Not(BeEmpty()))
//	//	//})
//	//
//	//	By("deleting the VM")
//	//	request = fmt.Sprintf(`{
//	//		  "method": "delete_vm",
//	//		  "arguments": ["%v"]
//	//		}`, vmCID)
//	//	assertSucceeds(request)
//	//})
//	//
//	//It("executes the VM lifecycle with a custom service account and scopes", func() {
//	//	By("creating a VM")
//	//	var vmCID string
//	//	request := fmt.Sprintf(`{
//	//		  "method": "create_vm",
//	//		  "arguments": [
//	//			"agent",
//	//			"%v",
//	//			{
//	//			  "machine_type": "n1-standard-1",
//	//			  "zone": "%v",
//	//			  "service_account": "%v",
//	//			  "service_scopes": ["devstorage.read_write"]
//	//			},
//	//			{
//	//			  "default": {
//	//				"type": "dynamic",
//	//				"cloud_properties": {
//	//				  "tags": ["integration-delete"],
//	//				  "network_name": "%v"
//	//				}
//	//			  }
//	//			},
//	//			[],
//	//			{}
//	//		  ]
//	//		}`, existingStemcell, zone, serviceAccount, networkName)
//	//	vmCID = assertSucceedsWithResult(request).(string)
//	//	//assertValidVM(vmCID, func(instance *compute.Instance) {
//	//	//	// Labels should be an exact match
//	//	//	Expect(instance.ServiceAccounts[0].Scopes).To(Not(BeEmpty()))
//	//	//	Expect(instance.ServiceAccounts[0].Email).To(Equal(serviceAccount))
//	//	//})
//	//
//	//	By("deleting the VM")
//	//	request = fmt.Sprintf(`{
//	//		  "method": "delete_vm",
//	//		  "arguments": ["%v"]
//	//		}`, vmCID)
//	//	assertSucceeds(request)
//	//})
//	//
//	//It("executes the VM lifecycle with a custom service account and no scopes", func() {
//	//	By("creating a VM")
//	//	var vmCID string
//	//	request := fmt.Sprintf(`{
//	//		  "method": "create_vm",
//	//		  "arguments": [
//	//			"agent",
//	//			"%v",
//	//			{
//	//			  "machine_type": "n1-standard-1",
//	//			  "zone": "%v",
//	//			  "service_account": "%v"
//	//			},
//	//			{
//	//			  "default": {
//	//				"type": "dynamic",
//	//				"cloud_properties": {
//	//				  "tags": ["integration-delete"],
//	//				  "network_name": "%v"
//	//				}
//	//			  }
//	//			},
//	//			[],
//	//			{}
//	//		  ]
//	//		}`, existingStemcell, zone, serviceAccount, networkName)
//	//	vmCID = assertSucceedsWithResult(request).(string)
//	//	//assertValidVM(vmCID, func(instance *compute.Instance) {
//	//	//	// Labels should be an exact match
//	//	//	Expect(instance.ServiceAccounts[0].Scopes).To(Not(BeEmpty()))
//	//	//	Expect(instance.ServiceAccounts[0].Email).To(Equal(serviceAccount))
//	//	//})
//	//
//	//	By("deleting the VM")
//	//	request = fmt.Sprintf(`{
//	//		  "method": "delete_vm",
//	//		  "arguments": ["%v"]
//	//		}`, vmCID)
//	//	assertSucceeds(request)
//	//})
//
//})