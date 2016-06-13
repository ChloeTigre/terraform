/* test our govmomi wrappings to ensure they work properly */
package vsphere

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/vmware/govmomi"
)

var testParameters map[string]interface{}
var client *govmomi.Client

func init() {
	var err error
	client, err = getTestGovmomiClient()
	if err != nil {
		panic("Cannot get govmomi client.")
	}
}
func buildTestDVS(variant string) *dvs {
	dvsO := dvs{}
	dvsO.datacenter = testParameters["datacenter"].(string)
	dvsO.folder = testParameters["switchFolder"].(string)
	dvsO.name = fmt.Sprintf(testParameters["switchName"].(string), variant)
	dvsO.description = testParameters["switchDescription"].(string)
	dvsO.contact.infos = testParameters["contactInfos"].(string)
	dvsO.contact.name = testParameters["contactName"].(string)
	dvsO.numStandalonePorts = testParameters["numStandalonePorts"].(int)
	dvsO.switchUsagePolicy.autoPreinstallAllowed = true
	dvsO.switchUsagePolicy.autoUpgradeAllowed = true
	dvsO.switchUsagePolicy.partialUpgradeAllowed = false
	dvsO.switchIPAddress = testParameters["switchIPAddress"].(string)
	return &dvsO
}

func buildTestDVPG(variant string, dvsInfo *dvs) *dvs_port_group {
	dvpg := dvs_port_group{}
	dvpg.autoExpand = testParameters["dvpgAutoExpand"].(bool)
	dvpg.name = fmt.Sprintf(testParameters["portgroupName"].(string), variant)
	dvpg.numPorts = testParameters["portgroupPorts"].(int)
	dvpg.switchId = dvsInfo.getID()
	dvpg.description = testParameters["portgroupDescription"].(string)
	dvpg.pgType = "earlyBinding"
	return &dvpg
}

func testCreateDVS(dvsObject *dvs, client *govmomi.Client) error {
	return dvsObject.createSwitch(client)
}

func testDeleteDVS(dvsObject *dvs, client *govmomi.Client) error {
	return dvsObject.Destroy(client)
}

func doCreateDVS(dvsO *dvs, t *testing.T) {
	var err error
	if err = testCreateDVS(dvsO, client); err != nil {
		t.Logf("[ERROR] Cannot create switch: %+v\n", err)
		t.Fail()
	}
	t.Log("Created DVS. Now getting props")

	props, err := dvsO.getProperties(client)
	if err != nil {
		t.Logf("Cannot retrieve DVS properties, failing: [%T]%+v\nProperties obj: [%T]%+v\n", err, err, props, props)

		t.Fail()
	} else {
		t.Logf("Properties: %+v", props)
	}
}

func doDeleteDVS(dvsO *dvs, t *testing.T) {
	if err := testDeleteDVS(dvsO, client); err != nil {
		t.Logf("[ERROR] Cannot delete switch: %+v\n", err)
		t.Fail()
	}
}

// Test DVS creation and destruction
func TestDVSCreationAndDestruction(t *testing.T) {
	// need:
	// datacenter name, host name 1, host name 2, switch path
	dvsO := buildTestDVS("test1")
	doCreateDVS(dvsO, t)
	doDeleteDVS(dvsO, t)
}

func testCreateDVPG(dvpg *dvs_port_group, client *govmomi.Client) error {
	return dvpg.createPortgroup(client)
}

func testDeleteDVPG(dvpg *dvs_port_group, client *govmomi.Client) error {
	return dvpg.deletePortgroup(client)
}

func doCreateDVPortgroup(dvpg *dvs_port_group, t *testing.T) {
	t.Logf("Create DVPG: %+v", dvpg)
	if err := testCreateDVPG(dvpg, client); err != nil {
		t.Logf("[ERROR] Cannot create portgroup: %+v\n", err)
		t.Fail()
	}
	t.Log("Created DVPG. Now getting props")

	props, err := dvpg.getProperties(client)
	if err != nil {
		t.Logf("Cannot retrieve DVPS properties, failing: [%T]%+v\nProperties obj: [%T]%+v\n", err, err, props, props)
		t.Fail()
	} else {
		t.Logf("Properties: %+v", props)
	}
}

func doDeleteDVPortgroup(dvpg *dvs_port_group, t *testing.T) {
	t.Logf("Delete DVPG")
	if err := testDeleteDVPG(dvpg, client); err != nil {
		t.Logf("[ERROR] Cannot delete portgroup: %+v\n", err)
		t.Fail()
	}
}

func TestPortgroupCreationAndDestruction(t *testing.T) {
	// need:
	// datacenter name, switch path, portgroup name
	dvsO := buildTestDVS("test2")
	dvpg := buildTestDVPG("test2", dvsO)
	doCreateDVS(dvsO, t)
	t.Logf("DVPG: %+v", dvpg)
	doCreateDVPortgroup(dvpg, t)
	time.Sleep(10 * time.Second)
	doDeleteDVPortgroup(dvpg, t)
	doDeleteDVS(dvsO, t)
}

// Test VM-DVS binding creation and destruction
func TestVMDVSCreationAndDestruction(t *testing.T) {
	// need:
	// datacenter name, switch path, portgroup name, VM path name

	t.FailNow()
}

// Test read DVS
func TestDVSRead(t *testing.T) {
	// need:
	// datacenter name, switch path
	t.FailNow()
}

// Test read Portgroup
func TestPortgroupRead(t *testing.T) {
	// need:
	// datacenter name, switch path, portgroup name
	t.FailNow()
}

// Test read VM-DVS binding
func TestVMDVSRead(t *testing.T) {
	// need:
	// datacenter name, switch path, portgroup name, VM path name
	t.FailNow()
}

func init() {
	datacenter := os.Getenv("VSPHERE_TEST_DC")
	switchFolder := os.Getenv("VSPHERE_TEST_SWDIR")
	vmFolder := os.Getenv("VSPHERE_TEST_VMDIR")
	if datacenter == "" {
		datacenter = "vm-test-1"
	}
	if switchFolder == "" {
		switchFolder = "/DEVTESTS"
	}
	if vmFolder == "" {
		vmFolder = "/DEVTESTS"
	}
	testParameters = make(map[string]interface{})
	testParameters["datacenter"] = datacenter
	testParameters["switchFolder"] = switchFolder
	testParameters["vmFolder"] = vmFolder
	testParameters["vmPath"] = "VMTEST1-%s"
	testParameters["switchName"] = "DVSTEST-%s"
	testParameters["portgroupName"] = "PORTGROUPTEST1-%s"
	testParameters["switchDescription"] = "lorem test ipsum test"
	testParameters["portgroupDescription"] = "doler test sit amet test"
	testParameters["numStandalonePorts"] = 4
	testParameters["portgroupPorts"] = 16
	testParameters["contactInfos"] = "lorem test <test@example.invalid>"
	testParameters["contactName"] = "Lorem Test Ipsum Invalid"
	testParameters["switchIPAddress"] = "192.0.2.1"
	testParameters["dvpgAutoExpand"] = true
}
