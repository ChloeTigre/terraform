package vsphere

import (
	"testing"
	"fmt"
	"log"
	"os"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

// test we can create a DVS
func TestAccVSphereDVS_create(t *testing.T) {
	handleName := "dvs_testacc"
	resourceName := "dvs_testacceptance"
	dvsConfigFilled := fmt.Sprintf(
		dvsConfig,
		handleName, // resource handle name
		resourceName, // resource name
		os.Getenv("VSPHERE_DATACENTER"), // datacenter
		"", // extension_key
		"An Acceptance Test DVS - temporary", // description
		"Terraform Test", // contact.name
		"Non Existent <terraform@example.com>", // contact.infos
		"true", // auto_preinstall_allowed
		"true", // auto_upgrade_allowed
		"true", // partial_upgrade_allowed,
		"198.51.100.1", // switch_ip_address
		"5")// num_standalone_ports
	log.Printf("create: using config\n%s", dvsConfigFilled)
	resource.Test(t, resource.TestCase{
		PreCheck:	func() { testAccPreCheck(t) },
		Providers:	testAccProviders,
		CheckDestroy:	nil,
		Providers:	testAccProviders,
		Steps:		[]resource.TestStep{
			resource.TestStep{
				Config: dvsConfigFilled,
				Check: resource.ComposeTestCheckFunc(wTestDVSExists(resourceName))
			},
		},

	})
}

func wTestDVSExists(handleName string) {
	return func(s *terraform.State) error {
		// this internal function must test whether
		// `handleName` exists in the visible vSphere.
		// if return nil: success, else failure
		return nil
	}
}

// test we can read a DVS
func TestAccVSphereDVS_read(t *testing.T) {
}

// test we can update a DVS
func TestAccVSphereDVS_update(t *testing.T) {

}
// test we can delete a DVS
func TestAccVSphereDVS_delete(t *testing.T) {
}

// test we can create a DVSPortGroup
func TestAccVSphereDVSPortGroup_create(t *testing.T) {
}

// test we can read a DVSPortGroup
func TestAccVSphereDVSPortGroup_read(t *testing.T) {
}

// test we can update a DVSPortGroup
func TestAccVSphereDVSPortGroup_update(t *testing.T) {

}
// test we can delete a DVSPortGroup
func TestAccVSphereDVSPortGroup_delete(t *testing.T) {
}

// test DVSMapHostDvs
func TestAccVSphereDVSMapHostDVS_create(t *testing.T) {}
func TestAccVSphereDVSMapHostDVS_read(t *testing.T) {}
func TestAccVSphereDVSMapHostDVS_update(t *testing.T) {}
func TestAccVSphereDVSMapHostDVS_delete(t *testing.T) {}

// test DVSMapVmDVSPG
func TestAccVSphereDVSMapVMDVPG_create(t *testing.T) {}
func TestAccVSphereDVSMapVMDVPG_read(t *testing.T) {}
func TestAccVSphereDVSMapVMDVPG_update(t *testing.T) {}
func TestAccVSphereDVSMapVMDVPG_delete(t *testing.T) {}



// definition of the basic DVS config (without host)
/*
	maps
	DVSConfigSpec
		-> DVSContactInfo
		-> DVSPolicy
*/
const dvsConfig = `
resource "vsphere_dvs" "%s" {
	name = "%s"
	datacenter = "%s"
	extension_key = "%s"
	description = "%s"
	contact {
		name = "%s"
		infos = "%s"
	}
	switch_usage_policy {
		auto_preinstall_allowed = "%s"
		auto_upgrade_allowed = "%s"
		partial_upgrade_allowed = "%s"
  	}
	switch_ip_address = "%s"
	num_standalone_ports = "%s"
}
`


// represent a DistributedVirtualSwitchHostMemberConfigSpec
// minimal, default support
const dvsMapHostDvs = `
resource "vsphere_dvs_host_map" "%s" {
	host_name = "%s"
	switch_name = "%s"
}
`

// represent a DVPortgroupConfigSpec
// type: earlyBinding|ephemeral

const dvsPortGroup = `
resource "vsphere_dvs_port_group" "%s" {
	name = "%s"
	pg_type = "%s"
	description = "%s"
	auto_expand = "%s"
	num_ports = "%s"
	port_name_format = "%s"
	policy {
		allow_block_override = "%s"
		allow_live_port_moving = "%s"
		allow_network_rp_override = "%s"
		port_config_reset_disconnect = "%s"
		allow_shaping_override = "%s"
		allow_traffic_filter_override = "%s"
		allow_vendor_config_override = "%s"
	}
}
`

// represent a VirtualEthernetCardDistributedVirtualPortBackingInfo
// (yes VMware loves huge names)
const dvsVMPort = `
resource "vsphere_dvs_vm_port" "%s" {
	vm = "%s"
	nic_index = "%s"
	portgroup = "%s"
	port = "%s"
}
`
