package vsphere

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/object"
	"golang.org/x/net/context"
)

// test we can create a DVS
func TestAccVSphereDVS_create(t *testing.T) {

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

// test we can bind hosts to a DVS
func TestAccVSphereDVS_bind_host(t *testing.T) {

}

// test we can unbind hosts from a DVS
func TestAccVSphereDVS_bind_host(t *testing.T) {

}

// test we can bind VMs to a DVSPG
func TestAccVSphereDVSPG_bind_vm(t *testing.T) {

}

// test we can unbind VMs from a DVSPG
func TestAccVSphereDVSPG_unbind_vm(t *testing.T) {

}

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
	key = "%s"
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
