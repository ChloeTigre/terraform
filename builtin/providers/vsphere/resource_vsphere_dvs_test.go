package vsphere

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/vmware/govmomi"
	"golang.org/x/net/context"
)

// test we can create a DVS
func TestAccVSphereDVS_create(t *testing.T) {
	handleName := "dvs_testacc"
	resourceName := "dvs_testacceptance"
	dvsConfigFilled := fmt.Sprintf(
		dvsConfig,
		handleName,                      // resource handle name
		resourceName,                    // resource name
		os.Getenv("VSPHERE_DATACENTER"), // datacenter
		"", // extension_key
		"An Acceptance Test DVS - temporary",   // description
		"Terraform Test",                       // contact.name
		"Non Existent <terraform@example.com>", // contact.infos
		"true",         // auto_preinstall_allowed
		"true",         // auto_upgrade_allowed
		"true",         // partial_upgrade_allowed,
		"198.51.100.1", // switch_ip_address
		"5")            // num_standalone_ports
	log.Printf("create: using config\n%s", dvsConfigFilled)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: nil,
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: dvsConfigFilled,
				Check:  wTestDVSExists(resourceName),
			},
		},
	})
}

func wTestDVSExists(handleName string) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		// this internal function must test whether
		// `handleName` exists in the visible vSphere.
		// if return nil: success, else failure
		client, datacenter, err := boilerplateClient(handleName, s)
		if err != nil {
			return err
		}
		// here we will need to set-up a proper resource type for the Provider
		// itemConf := s.RootModule().Resources[handleName].(dvs)
		dvso := dvs{}
		return loadDVS(client, datacenter, handleName, &dvso)
	}
}

func boilerplateClient(n string, s *terraform.State) (client *govmomi.Client, datacenter string, err error) {
	err = nil
	client = nil
	rs, ok := s.RootModule().Resources[n]
	if !ok {
		err = fmt.Errorf("Resource not found: %s", n)
		return
	}

	if rs.Primary.ID == "" {
		err = fmt.Errorf("No ID is set")
		return
	}

	client = testAccProvider.Meta().(*govmomi.Client)
	return

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
func TestAccVSphereDVSMapHostDVS_read(t *testing.T)   {}
func TestAccVSphereDVSMapHostDVS_update(t *testing.T) {}
func TestAccVSphereDVSMapHostDVS_delete(t *testing.T) {}

// test DVSMapVmDVSPG
func TestAccVSphereDVSMapVMDVPG_create(t *testing.T) {}
func TestAccVSphereDVSMapVMDVPG_read(t *testing.T)   {}
func TestAccVSphereDVSMapVMDVPG_update(t *testing.T) {}
func TestAccVSphereDVSMapVMDVPG_delete(t *testing.T) {}

// unit tests for small pieces of code
func TestApiListNetwork(t *testing.T) {
	if os.Getenv("TESTVAR") == "" {
		t.SkipNow()
	}
	dvsitem := dvs{
		datacenter: os.Getenv("VSPHERE_DATACENTER"),
	}
	cli, err := getTestGovmomiClient()
	if err != nil {
		t.Log("Oops, could not get VMOMI client", err)
		t.FailNow()
	}
	dvss, err := dvsitem.getDVS(cli, os.Getenv("TESTVAR"))
	if err != nil {
		t.Log("Oops, could not load data:", err)
		t.FailNow()
	}
	t.Log("getDVS did not crash", dvss)

	if err := loadDVS(cli, dvsitem.datacenter, os.Getenv("TESTVAR"), &dvsitem); err != nil {
		t.Log("[fail] DVSS:", dvss)
		t.Log("loadDVS failed:", err)
		t.FailNow()
	}
	log.Printf("DVS: %+v", dvsitem)
}

func getTestGovmomiClient() (*govmomi.Client, error) {
	u, err := url.Parse("https://" + os.Getenv("VSPHERE_URL") + "/sdk")
	if err != nil {
		return nil, fmt.Errorf("Cannot parse VSPHERE_URL")
	}
	u.User = url.UserPassword(os.Getenv("VSPHERE_USER"), os.Getenv("VSPHERE_PASSWORD"))

	cli, err := govmomi.NewClient(context.TODO(), u, true)
	if err != nil {
		return nil, err
	}
	return cli, nil
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
	host = "%s"
	switch = "%s"
	nic_name = "%s"
}
`

// represent a DVPortgroupConfigSpec
// type: earlyBinding|ephemeral

const dvsPortGroup = `
resource "vsphere_dvs_port_group" "%s" {
	name = "%s"
	switch = "%s"
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
	nic_label = "%s"
	portgroup = "%s"
	port = "%s"
}
`
