/* test our govmomi wrappings to ensure they work properly */
package vsphere

import "testing"

var testParameters map[string]interface{}

// Test DVS creation and destruction
func TestDVSCreationAndDestruction(t *testing.T) {
	// need:
	// datacenter name, host name 1, host name 2, switch path

	t.FailNow()
}

func TestPortgroupCreationAndDestruction(t *testing.T) {
	// need:
	// datacenter name, switch path, portgroup name

	t.FailNow()
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
	testParameters["datacenter"] = "vm-test-1"
	testParameters["vmPath"] = "vm-test-1"
	testParameters["switchPath"] = "vm-test-1"
	testParameters["portgroupName"] = "vm-test-1"
}
